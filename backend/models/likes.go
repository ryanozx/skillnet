package models

type Like struct {
	ID     string `json:"-"`
	UserID string
	PostID uint
}

type LikeUpdate struct {
	Like      Like
	LikeCount uint64
}
