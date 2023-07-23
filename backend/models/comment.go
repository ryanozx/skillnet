package models

import (
	"gorm.io/gorm"
)

type CommentViewsArray struct {
	Comments    []CommentView
	NextPageURL string
}

func (cvArray *CommentViewsArray) TestFormat() *CommentViewsArray {
	if len(cvArray.Comments) == 0 {
		return &CommentViewsArray{
			NextPageURL: cvArray.NextPageURL,
		}
	}
	output := CommentViewsArray{
		Comments:    []CommentView{},
		NextPageURL: cvArray.NextPageURL,
	}
	for _, cv := range cvArray.Comments {
		output.Comments = append(output.Comments, *cv.TestFormat())
	}
	return &output
}

type Comment struct {
	gorm.Model
	PostID uint   `gorm:"<-:create; not null"`
	Post   Post   `json:"-"`
	UserID string `json:"-" gorm:"<-:create; not null"`
	User   User   `json:"-"`
	Text   string
}

func (comment *Comment) TestFormat() *Comment {
	output := Comment{
		Model:  comment.Model,
		PostID: comment.PostID,
		Text:   comment.Text,
	}
	return &output
}

type CommentView struct {
	Comment     Comment
	UserMinimal `json:"User"`
	IsEditable  bool
}

func (c *Comment) GetUserID() string {
	return c.UserID
}

func (comment *Comment) CommentView(userID string) *CommentView {
	commentView := CommentView{
		Comment:     *comment,
		UserMinimal: *comment.User.GetUserMinimal(),
		IsEditable:  userID == comment.UserID,
	}
	return &commentView
}

func (cv *CommentView) TestFormat() *CommentView {
	output := CommentView{
		Comment:     *cv.Comment.TestFormat(),
		UserMinimal: *cv.UserMinimal.TestFormat(),
		IsEditable:  cv.IsEditable,
	}
	return &output
}

type CommentUpdate struct {
	Comment      CommentView
	CommentCount uint64
}

func (update *CommentUpdate) TestFormat() *CommentUpdate {
	output := CommentUpdate{
		Comment:      *update.Comment.TestFormat(),
		CommentCount: update.CommentCount,
	}
	return &output
}
