package models

import (
	"gorm.io/gorm"
)

type CommentsArray struct {
	Comments    []CommentView
	NextPageURL string
}

type Comment struct {
	gorm.Model
	PostID uint
	UserID string
	Text   string
}

type CommentView struct {
	Comment     `gorm:"embedded"`
	UserMinimal `json:"User"`
}
