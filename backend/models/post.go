package models

import (
	"gorm.io/gorm"
)

// PostViewArray is a struct for supporting post feed pagination
type PostViewArray struct {
	Posts       []PostView `json:"Posts"`
	NextPageURL string
}

func (pvArray *PostViewArray) TestFormat() *PostViewArray {
	if len(pvArray.Posts) == 0 {
		return &PostViewArray{
			NextPageURL: pvArray.NextPageURL,
		}
	}
	output := PostViewArray{
		Posts:       []PostView{},
		NextPageURL: pvArray.NextPageURL,
	}
	for _, postView := range pvArray.Posts {
		output.Posts = append(output.Posts, *postView.TestFormat())
	}
	return &output
}

// Post is the database representation of a post object
type Post struct {
	gorm.Model
	UserID   string    `json:"-" gorm:"<-:create; not null"`
	User     User      `json:"-"`
	Content  string    `gorm:"not null"`
	Likes    []Like    `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	Comments []Comment `json:"-" gorm:"constraint:OnDelete:CASCADE"`
}

func (post *Post) TestFormat() *Post {
	output := Post{
		Model:    post.Model,
		Content:  post.Content,
		Likes:    post.Likes,
		Comments: post.Comments,
	}
	return &output
}

// PostView represents the information that the client receives
type PostView struct {
	Post         Post
	UserMinimal  `json:"User"`
	IsEditable   bool
	Liked        bool
	LikeCount    uint64
	CommentCount uint64
}

func (pv *PostView) TestFormat() *PostView {
	output := PostView{
		Post:         *pv.Post.TestFormat(),
		UserMinimal:  *pv.UserMinimal.TestFormat(),
		IsEditable:   pv.IsEditable,
		Liked:        pv.Liked,
		LikeCount:    pv.LikeCount,
		CommentCount: pv.CommentCount,
	}
	return &output
}

// MultimediaContent will be used to represent multimedia resources
// TODO: Add MultimediaContent array as attribute of Post
type MultimediaContent struct {
	ContentType string
	URI         string
}

func (mc *MultimediaContent) TestFormat() *MultimediaContent {
	return mc
}

type PostViewParams struct {
	UserID       string
	LikeCount    uint64
	CommentCount uint64
}

// Creates a PostView object
func (post *Post) PostView(params *PostViewParams) *PostView {
	postView := PostView{
		Post:         *post,
		UserMinimal:  *post.User.UserMinimal(),
		IsEditable:   params.UserID == post.UserID,
		Liked:        len(post.Likes) > 0 && post.Likes[0].UserID == params.UserID,
		LikeCount:    params.LikeCount,
		CommentCount: params.CommentCount,
	}
	return &postView
}

func (pv *PostView) GetPost() *Post {
	return &pv.Post
}

func (post *Post) GetUserID() string {
	return post.UserID
}
