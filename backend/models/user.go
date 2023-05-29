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
	UserID   uuid.UUID
}

type User struct {
	ID            uuid.UUID
	Username      string
	Name          string
	ProfilePicURI string
	Title         string
	Birthday      time.Time
	Location      string
	AboutMe       string
	Projects      ProjectsArray
}

type UserWithSettings struct {
	User              User
	TitlePrivileges   uint8
	BirthdayPrivilege uint8
	LocationPrivilege uint8
	AboutMePrivilege  uint8
	ProjectsPrivilege uint8
}

type UserWithValidationError struct {
	User               UserWithSettings `json:"userWithSettings"`
	UsernameError      error            `json:"usernameError"`
	NameError          error            `json:"nameError"`
	ProfilePicURIError error            `json:"profilePicError"`
	BirthdayError      error            `json:"birthdayError"`
	LocationError      error            `json:"locationError"`
	AboutMeError       error            `json:"aboutMeError"`
}

func generateTestUser(userID uuid.UUID) UserMinimal {
	testUser := UserMinimal{
		ID:         userID,
		Name:       "User1",
		ProfileURL: "www.skillnet.com/user/user1",
		ProfilePic: "user1PicURL",
	}
	return testUser
}
