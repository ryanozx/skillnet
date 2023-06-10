package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostViewArray is a struct for supporting post feed pagination
type PostViewArray struct {
	Posts       []PostView
	NextPageURL string
}

// Post is the database representation of a post object
type Post struct {
	gorm.Model
	UserID  uuid.UUID `json:"-" gorm:"<-:create; not null"`
	User    User      `json:"-"`
	Content string    `gorm:"not null"`
}

// PostView represents the information that the client receives
type PostView struct {
	Post          Post
	UserMinimal   `json:"User"`
	CommentsArray `json:"Comments"`
}

// MultimediaContent will be used to represent multimedia resources
// TODO: Add MultimediaContent array as attribute of Post
type MultimediaContent struct {
	ContentType string
	URI         string
}

// Creates a PostView object
func (post *Post) PostView() *PostView {
	postView := PostView{
		Post:        *post,
		UserMinimal: *post.User.UserMinimal(),
	}
	return &postView
}
