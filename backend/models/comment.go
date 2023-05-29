package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentsArray struct {
	Comments    []CommentSchema `json:"comments"`
	NextPageURL string          `json:"nextPageURL"`
}

type CommentSchema struct {
	gorm.Model
	PostID uuid.UUID   `json:"postID"`
	User   UserMinimal `json:"user" gorm:"embedded"`
	Text   string      `json:"text"`
}
