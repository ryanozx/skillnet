package models

import "github.com/google/uuid"

type Community struct {
	ID       uuid.UUID     `json:"id" gorm:"primary_key"`
	Name     string        `json:"name"`
	Owner    UserMinimal   `json:"owner"`
	About    string        `json:"about"`
	Projects ProjectsArray `json:"projects"`
}
