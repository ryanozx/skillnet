package models

import (
	"time"

	"github.com/google/uuid"
)

type UserMinimal struct {
	ID         uuid.UUID
	Name       string
	ProfileURL string
	ProfilePic string
}

type UserCredentials struct {
	Username string
	Password string
}

type User struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;"`
	Username      string
	Password      string
	Name          string
	ProfilePicURI string
	Title         string
	Birthday      time.Time
	Location      string
	AboutMe       string
	// Projects      ProjectsArray
}

type UserWithSettings struct {
	User              User `gorm:"embedded"`
	TitlePrivileges   uint8
	BirthdayPrivilege uint8
	LocationPrivilege uint8
	AboutMePrivilege  uint8
	ProjectsPrivilege uint8
}

func generateTestUser(userID uuid.UUID) *UserMinimal {
	testUser := UserMinimal{
		ID:         userID,
		Name:       "User1",
		ProfileURL: "www.skillnet.com/user/user1",
		ProfilePic: "user1PicURL",
	}
	return &testUser
}

func (userCreds *UserCredentials) ConvertToUser() *User {
	newUser := User{
		ID:       uuid.New(),
		Username: userCreds.Username,
		Password: userCreds.Password,
	}
	return &newUser
}
