package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

// UserMinimal contains the bare essential user details that can
// be displayed next to posts, comments, etc.
type UserMinimal struct {
	Name       string
	URL        string `gorm:"-:all"`
	ProfilePic string
}

// UserCredentials handles login/signup information
type UserCredentials struct {
	Username string `gorm:"unique; not null"`
	Password string `gorm:"not null" json:"-"`
}

// UserView represents the information a visitor to a user's profile page can see
type UserView struct {
	UserMinimal   `gorm:"embedded"`
	Title         string
	Birthday      time.Time
	Location      string
	AboutMe       string
	ProjectsArray `json:"projects" gorm:"-:all"`
}

// User is the database representation of a user object
type User struct {
	ID              string   `gorm:"type:uuid;<-:create;default:uuid_generate_v4()" json:"-"` // UserID will never be revealed to the client; it will never change
	Profile         UserView `gorm:"embedded"`
	UserCredentials `gorm:"embedded"`
	Email           string
}

func (userCreds *UserCredentials) NewUser() *User {
	user := User{
		UserCredentials: *userCreds,
	}
	return &user
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	user.ID = uuid.NewString()
	return nil
}

func (user *User) AfterFind(tx *gorm.DB) error {
	user.GenerateProfileURL()
	return nil
}

func (user *User) AfterUpdate(tx *gorm.DB) error {
	user.GenerateProfileURL()
	return nil
}

func (user *User) UserMinimal() *UserMinimal {
	return &user.Profile.UserMinimal
}

func (profile *UserView) GenerateUserMinimal() *UserMinimal {
	return &profile.UserMinimal
}

func (user *User) GenerateProfileURL() {
	const profileUrlPrefix = "localhost:8080/users/"
	url := profileUrlPrefix + user.Username
	user.Profile.URL = url
}
