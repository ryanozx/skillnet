package models

import (
	"gorm.io/gorm"
)

type CommentsArray struct {
	Comments    []CommentSchema
	NextPageURL string
}

type CommentSchema struct {
	gorm.Model
	PostID uint
	User   UserMinimal
	Text   string
}
