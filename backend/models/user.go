package models

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"gopkg.in/guregu/null.v3"
	"gorm.io/gorm"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

// UserMinimal contains the bare essential user details that can
// be displayed next to posts, comments, etc.
type UserMinimal struct {
	Name       string
	URL        string `gorm:"-:all"`
	ProfilePic string
}

// UserCredentials handles login information
type UserCredentials struct {
	Username string `gorm:"unique; not null"`
	Password string `gorm:"not null" json:"-"`
}

// UserCredentials for signup
type SignupUserCredentials struct {
	UserCredentials
	Email string `gorm:"not null"`
}

// UserView represents the information a visitor to a user's profile page can see
type UserView struct {
	UserMinimal `gorm:"embedded"`
	Title       null.String
	Birthday    time.Time
	Location    null.String
	AboutMe     null.String
	Projects    []ProjectMinimal `gorm:"-:all"`
}

// User is the database representation of a user object
type User struct {
	ID              string `gorm:"<-:create" json:"-"` // UserID will never be revealed to the client; it will never change
	UserView        `gorm:"embedded"`
	UserCredentials `gorm:"embedded"`
	Email           string `gorm:"not null"`
	Likes           []Like `json:"null" gorm:"constraint:OnDelete:CASCADE"`
}

func (userCreds *SignupUserCredentials) NewUser() *User {
	user := User{
		UserCredentials: userCreds.UserCredentials,
		Email:           userCreds.Email,
	}
	return &user
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	user.ID = uuid.NewString()
	if err := user.whiteSpaceCheck(tx); err != nil {
		return err
	}
	return nil
}

func (user *User) AfterFind(tx *gorm.DB) error {
	user.GenerateProfileURL()
	return nil
}

func (user *User) BeforeUpdate(tx *gorm.DB) error {
	if err := user.whiteSpaceCheck(tx); err != nil {
		return err
	}
	return nil
}

func (user *User) AfterUpdate(tx *gorm.DB) error {

	// This catches the cases where there is no change to the database entry, either
	// because there is no valid update, or interestingly enough if only the username
	// or password is being updated and the changed field(s) are empty. Getting the log
	// to print actualUser may be useful in debugging this error.
	emptyUser := &User{}
	actualUser := tx.Statement.Dest.(*User)

	// Ideally a different method should be used for deep equality since cmp.Equal
	// will panic if there is an unexported field - this can be mitigated with
	// cmp.AllowUnexported to ignore unexported fields but this is not recursive,
	// hence each subtype with unexported fields must be explicitly passed in.
	if cmp.Equal(emptyUser, actualUser, cmp.AllowUnexported(User{})) {
		return errors.New("bad request")
	}

	// The URL field is empty since it is not stored in the database, hence we must
	// generate it in the updated User object to the latest profile URL, especially
	// if the username was updated.
	user.GenerateProfileURL()
	return nil
}

func (user *User) UserMinimal() *UserMinimal {
	return user.GetUserView().GetUserMinimal()
}

func (user *User) GetUserView() *UserView {
	user.GenerateProfileURL()
	return &user.UserView
}

func (profile *UserView) GetUserMinimal() *UserMinimal {
	return &profile.UserMinimal
}

func (user *User) GenerateProfileURL() {
	var profileUrlPrefix = fmt.Sprintf("%s/profile/", ClientAddress)
	url := profileUrlPrefix + user.Username
	user.UserView.URL = url
}

func (user *User) whiteSpaceCheck(tx *gorm.DB) error {
	if tx.Statement.Changed("Username") {
		username := tx.Statement.Dest.(*User).Username
		isWhiteSpacePresent := regexp.MustCompile(`\s`).MatchString(username)
		if isWhiteSpacePresent {
			return errors.New("username cannot contain whitespace")
		}
	}
	return nil
}
