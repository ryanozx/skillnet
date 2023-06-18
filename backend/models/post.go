package models

import (
	"gorm.io/gorm"
)

// PostViewArray is a struct for supporting post feed pagination
type PostViewArray struct {
	Posts       []PostView `json:"Posts"`
	NextPageURL string
}

// Post is the database representation of a post object
type Post struct {
	gorm.Model
	UserID  string `json:"-" gorm:"<-:create; not null"`
	User    User   `json:"-"`
	Content string `gorm:"not null"`
}

// PostView represents the information that the client receives
type PostView struct {
	Post          Post
	UserMinimal   `json:"User"`
	CommentsArray `json:"Comments"`
	IsEditable    bool
}

// MultimediaContent will be used to represent multimedia resources
// TODO: Add MultimediaContent array as attribute of Post
type MultimediaContent struct {
	ContentType string
	URI         string
}

// Creates a PostView object
func (post *Post) PostView(userID string) *PostView {
	postView := PostView{
		Post:        *post,
		UserMinimal: *post.User.UserMinimal(),
		IsEditable:  userID == post.UserID,
	}
	return &postView
}

func (pv *PostView) GetPost() *Post {
	return &pv.Post
}
