package models

import "gorm.io/gorm"

type ProjectsArray struct {
	Projects    []ProjectMinimal `json:"projects"`
	NextPageURL string
}

type ProjectMinimal struct {
	Name          string
	URL           string `gorm:"-:all"`
	ProjectImgURI string
	ProjectInfo   string
}

type Project struct {
	gorm.Model
	ProjectMinimal `gorm:"embedded"`
	Owner          UserMinimal
	Members        []UserMinimal `gorm:"-:all"`
	Posts          PostViewArray `gorm:"-:all"`
}
