package models

import "gorm.io/gorm"

type ProjectsArray struct {
	Projects    []ProjectMinimal `gorm:"-:all"`
	NextPageURL string           `gorm:"-:all"`
}

type ProjectMinimal struct {
	Name          string
	URL           string
	ProjectImgURI string
}

type Project struct {
	gorm.Model
	Name        string
	Owner       UserMinimal
	Members     []UserMinimal
	ProjectInfo string
	Posts       PostViewArray `gorm:"-:all"`
}
