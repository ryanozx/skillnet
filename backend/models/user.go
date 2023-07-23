package models

import (
	"errors"
	"fmt"
	"regexp"

	"gopkg.in/guregu/null.v3"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

// UserMinimal contains the bare essential user details that can
// be displayed next to posts, comments, etc.
type UserMinimal struct {
	Name       null.String
	URL        string `gorm:"-:all"`
	ProfilePic string
}

func (user *UserMinimal) TestFormat() *UserMinimal {
	return user
}

// UserCredentials handles login information
type UserCredentials struct {
	Username string `gorm:"unique; not null"`
	Password string `gorm:"not null" json:"-"`
}

func (uc *UserCredentials) TestFormat() *UserCredentials {
	output := UserCredentials{
		Username: uc.Username,
	}
	return &output
}

// UserCredentials for signup
type SignupUserCredentials struct {
	UserCredentials
	Email string `gorm:"not null"`
}

func (creds *SignupUserCredentials) TestFormat() *SignupUserCredentials {
	output := SignupUserCredentials{
		UserCredentials: *creds.UserCredentials.TestFormat(),
		Email:           creds.Email,
	}
	return &output
}

// UserView represents the information a visitor to a user's profile page can see
type UserView struct {
	UserMinimal `gorm:"embedded"`
	Title       null.String
	AboutMe     null.String
	ShowTitle   bool
	ShowAboutMe bool
}

func (uv *UserView) TestFormat() *UserView {
	output := UserView{
		UserMinimal: *uv.UserMinimal.TestFormat(),
		Title:       uv.Title,
		AboutMe:     uv.AboutMe,
		ShowTitle:   uv.ShowTitle,
		ShowAboutMe: uv.ShowAboutMe,
	}
	return &output
}

// User is the database representation of a user object
type User struct {
	ID              string `gorm:"<-:create" json:"-"` // UserID will never be revealed to the client; it will never change
	UserView        `gorm:"embedded"`
	UserCredentials `gorm:"embedded"`
	Email           string    `json:"-" gorm:"not null"`
	Likes           []Like    `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	Comments        []Comment `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	Projects        []Project `json:"-" gorm:"constraint:OnDelete:CASCADE;foreignKey:OwnerID"`
}

func (user *User) TestFormat() *User {
	output := User{
		UserView:        *user.UserView.TestFormat(),
		UserCredentials: *user.UserCredentials.TestFormat(),
	}
	return &output
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
	user.URL = GenerateProfileURL(user)
	return nil
}

func (user *User) BeforeUpdate(tx *gorm.DB) error {
	if err := user.whiteSpaceCheck(tx); err != nil {
		return err
	}
	return nil
}

func (user *User) GetUserView(viewerID string) *UserView {
	output := UserView{
		UserMinimal: *user.GetUserMinimal(),
		ShowTitle:   user.ShowTitle,
		ShowAboutMe: user.ShowAboutMe,
	}
	isOwnProfile := user.ID == viewerID
	fmt.Printf("userID: %s ViewerID: %s\n", user.ID, viewerID)
	if isOwnProfile || user.ShowTitle {
		output.Title = user.Title
	}
	if isOwnProfile || user.ShowAboutMe {
		output.AboutMe = user.AboutMe
	}
	return &output
}

func (user *User) GetUserMinimal() *UserMinimal {
	user.URL = GenerateProfileURL(user)
	return &user.UserMinimal
}

func GenerateProfileURL(user *User) string {
	url := fmt.Sprintf("%s/profile/%s", ClientAddress, user.Username)
	return url
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
