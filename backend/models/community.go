package models

import (
	"gopkg.in/guregu/null.v3"
	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name     string
	Owner    UserMinimal
	About    null.String
	Projects ProjectsArray `gorm:"-:all"`
}
