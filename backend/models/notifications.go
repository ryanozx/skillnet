package models

import (
	"time"
)

type Notification struct {
	SenderId   string    `json:"sender_id"`
	ReceiverId string    `json:"receiver_id"`
	Type       string    `json:"type"`
	PostId     *string   `json:"post_id,omitempty"`
	CommentId  *string   `json:"comment_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
