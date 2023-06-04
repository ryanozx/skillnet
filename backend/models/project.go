package models

import "github.com/google/uuid"

/*
type ProjectsArray struct {
	Projects    []ProjectMinimal
	NextPageURL string
}
*/

type ProjectMinimal struct {
	Name          string `json:"name"`
	URL           string `json:"url"`
	ProjectImgURI string `json:"picURI"`
}

type Project struct {
	ID    uuid.UUID   `json:"id" gorm:"primary_key"`
	Name  string      `json:"name"`
	Owner UserMinimal `json:"owner"`
	// Members     []UserMinimal `json:"members"`
	ProjectInfo string    `json:"projectInfo"`
	Posts       PostArray `json:"posts"`
}
