package helpers

import (
	"time"

	"github.com/ryanozx/skillnet/models"
)

func GenerateEventNotification(senderID string, receiverID string, notifText string) *models.Notification {
	output := models.Notification{
		SenderId:   senderID,
		CreatedAt:  time.Now(),
		ReceiverId: receiverID,
		Content:    notifText,
	}
	return &output
}

func GenerateLikeNotification(liker *models.User, receiverID string) *models.Notification {
	notifText := liker.Username + " liked your post"
	return GenerateEventNotification(liker.ID, receiverID, notifText)
}

func GenerateCommentNotification(commenter *models.User, receiverID string) *models.Notification {
	notifText := commenter.Username + " commented on your post"
	return GenerateEventNotification(commenter.ID, receiverID, notifText)
}
