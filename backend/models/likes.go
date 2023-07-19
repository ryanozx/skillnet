package models

import "time"

type Like struct {
	ID        string    `json:"-"`
	CreatedAt time.Time `gorm:"<-:create" json:"-"`
	UserID    string    `json:"-"`
	User      User      `json:"-"`
	Post      Post      `json:"-"`
	PostID    uint
}

func (like *Like) TestFormat() *Like {
	output := Like{
		PostID: like.PostID,
	}
	return &output
}

type LikeUpdate struct {
	Like      Like
	LikeCount uint64
}

func (update *LikeUpdate) TestFormat() *LikeUpdate {
	output := LikeUpdate{
		Like:      *update.Like.TestFormat(),
		LikeCount: update.LikeCount,
	}
	return &output
}
