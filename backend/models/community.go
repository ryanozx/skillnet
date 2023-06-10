package models

import "gorm.io/gorm"

type Community struct {
	gorm.Model
	Name     string
	Owner    UserMinimal
	About    string
	Projects ProjectsArray `gorm:"-:all"`
}
