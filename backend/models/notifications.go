package models

import (
	"time"
)

type Notification struct {
	ReceiverId string `json:"receiver_id"`
	Content    string `json:"content"`
	SenderId   string `json:"sender_id"`
	CreatedAt  time.Time
}
