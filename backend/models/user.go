package models

import (
	"time"

	"github.com/google/uuid"
)

type UserMinimal struct {
	ID         uuid.UUID `json:"id" gorm:"primary_key"`
	Name       string    `json:"name"`
	ProfileURL string    `json:"url"`
	ProfilePic string    `json:"profilePic"`
}

type UserCredentials struct {
	Username string    `json:"username"`
	Password string    `json:"password"`
	ID       uuid.UUID `json:"id" gorm:"primary_key"`
}

type User struct {
	ID            uuid.UUID     `json:"id" gorm:"primary_key"`
	Username      string        `json:"username"`
	Name          string        `json:"name"`
	ProfilePicURI string        `json:"profilePic"`
	Title         string        `json:"title"`
	Birthday      time.Time     `json:"birthday"`
	Location      string        `json:"location"`
	AboutMe       string        `json:"aboutMe"`
	Projects      ProjectsArray `json:"projects"`
}

type UserWithSettings struct {
	User              User  `json:"user"`
	TitlePrivileges   uint8 `json:"titlePrivilege"`
	BirthdayPrivilege uint8 `json:"birthdayPrivilege"`
	LocationPrivilege uint8 `json:"locationPrivilege"`
	AboutMePrivilege  uint8 `json:"aboutMePrivilege"`
	ProjectsPrivilege uint8 `json:"projectsPrivilege"`
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
