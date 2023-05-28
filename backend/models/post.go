package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostArray struct {
	Posts       []Post `json:"posts"`
	NextPageURL string `json:"nextPageURL"`
}

type PostSchema struct {
	gorm.Model
	UserID    uuid.UUID `json:"userID" gorm:"text;not null;default:null"`
	ProjectID uuid.UUID `json:"projectID" gorm:"text;not null; default:null"`
	Content   string    `json:"content" gorm:"text;default:null"`
}

type CreatePostInput struct {
	ProjectID uuid.UUID `json:"projectID"`
	Content   string    `json:"content" binding:"required"`
}

type UpdatePostInput struct {
	Content string `json:"content" binding:"required"`
}

type Post struct {
	ID        uint          `json:"id" gorm:"primary_key"`
	User      UserMinimal   `json:"user"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
	ProjectID uuid.UUID     `json:"projectID"`
	Content   string        `json:"content"`
	Comments  CommentsArray `json:"comments"`
}

// MultimediaContent will be used to represent multimedia resources
// TODO: Add MultimediaContent array as attribute of Post
type MultimediaContent struct {
	ContentType string `json:"contentType"`
	URI         string `json:"uri"`
}

func ConvertInputToPostSchema(input CreatePostInput) PostSchema {
	newDBPostEntry := PostSchema{
		UserID:    uuid.New(),
		ProjectID: uuid.New(),
		Content:   input.Content,
	}
	return newDBPostEntry
}

func ConvertPostSchemaToPost(schema PostSchema) Post {
	postOutput := Post{
		ID:        schema.ID,
		User:      generateTestUser(schema.UserID),
		CreatedAt: schema.CreatedAt,
		UpdatedAt: schema.UpdatedAt,
		ProjectID: schema.ProjectID,
		Content:   schema.Content,
	}
	return postOutput
}
