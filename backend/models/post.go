package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostArray struct {
	Posts       []Post
	NextPageURL string
}

type PostSchema struct {
	gorm.Model
	UserID    uuid.UUID
	ProjectID uuid.UUID
	Content   string
}

type CreatePostInput struct {
	Content string
}

type UpdatePostInput struct {
	Content string
}

type Post struct {
	ID        uint
	UserID    uuid.UUID
	User      UserMinimal
	CreatedAt time.Time
	UpdatedAt time.Time
	ProjectID uuid.UUID
	Content   string
	CommentsArray
}

// MultimediaContent will be used to represent multimedia resources
// TODO: Add MultimediaContent array as attribute of Post
type MultimediaContent struct {
	ContentType string
	URI         string
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
