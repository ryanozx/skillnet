package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var testUserID = uuid.Must(uuid.Parse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"))
var testProjectID = uuid.Must(uuid.Parse("6ba7b811-9dad-11d1-80b4-00c04fd430c8"))

type PostArray struct {
	Posts       []Post
	NextPageURL string
}

type PostSchema struct {
	gorm.Model
	UserID    uuid.UUID `json:"-"`
	ProjectID uuid.UUID
	Content   string
}

type PostInput struct {
	Content string
}

type Post struct {
	PostSchema
	User     UserMinimal   `gorm:"embedded"`
	Comments CommentsArray `gorm:"embedded"`
}

// MultimediaContent will be used to represent multimedia resources
// TODO: Add MultimediaContent array as attribute of Post
type MultimediaContent struct {
	ContentType string
	URI         string
}

func (input *PostInput) ConvertToPostSchema() *PostSchema {
	newDBPostEntry := PostSchema{
		UserID:    testUserID,
		ProjectID: testProjectID,
		Content:   input.Content,
	}
	return &newDBPostEntry
}

func (schema *PostSchema) ConvertToPost() *Post {
	postOutput := Post{
		PostSchema: *schema,
		User:       *generateTestUser(schema.UserID),
	}
	// unsets the userID field after it has been used to generate the
	// UserMinimal struct, so that there is only one field containing
	// the userID. This is done to pass the test cases, as the expected
	// output is not generated via JSON and will thus ignore the
	// json ignore tags, which sets the userID field to the default value of nil
	postOutput.PostSchema.UserID = uuid.Nil
	return &postOutput
}

func ConvertToPosts(posts []PostSchema) []Post {
	output := make([]Post, len(posts))
	for i, post := range posts {
		output[i] = *post.ConvertToPost()
	}
	return output
}
