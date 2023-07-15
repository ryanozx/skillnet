package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
)

func (a *APIEnv) GetNotifications(context *gin.Context) {
	senderId := helpers.GetUserIdFromContext(context)
	receiverKey := "notifications:" + senderId
	// We will set up the request as Server Sent Events.
	context.Header("Content-Type", "text/event-stream")
	context.Header("Cache-Control", "no-cache")
	context.Header("Connection", "keep-alive")

	ctx := context.Request.Context()

	// Get all pending notifications from Redis
	results, err := a.Redis.ZRange(ctx, receiverKey, 0, -1).Result()
	if err != nil {
		// log.Printf("Error getting notifications: %v\n", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting notifications"})
		return
	}

	// Remove pending notifications from Redis
	a.Redis.Del(ctx, receiverKey)
	flusher, ok := context.Writer.(http.Flusher)
	if !ok {
		panic("expected gin.ResponseWriter to be an http.Flusher")
	}

	// Send all notifications to the client
	for _, result := range results {
		context.Writer.Write([]byte(fmt.Sprintf("%s\n\n", result)))
		flusher.Flush()
	}

	// Use redis pub/sub system
	pubsub := a.Redis.Subscribe(ctx, receiverKey)
	_, errPubSub := pubsub.Receive(ctx)
	if errPubSub != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error connecting to redis"})
		return
	}
	ch := pubsub.Channel()

	// Handle client closing connection
	closeNotifier, ok := context.Writer.(http.CloseNotifier)
	if !ok {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	for {
		select {
		case <-closeNotifier.CloseNotify():
			pubsub.Close()
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			// Each time we receive a message from the Redis channel,
			// we will send an SSE containing the message payload
			// context.SSEvent("message", msg.Payload)
			context.Writer.Write([]byte(fmt.Sprintf("%s\n\n", msg.Payload)))
			flusher.Flush()
		default:
			// This is necessary to prevent the loop from running infinitely
			// and consuming a lot of CPU power.
			time.Sleep(time.Second * 1)
		}
	}
}

func (a *APIEnv) PostNotification(context *gin.Context) {
	senderId := helpers.GetUserIdFromContext(context)

	// var notif models.Notification
	var notif models.Notification
	if err := context.ShouldBindJSON(&notif); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Setting the senderId and createdAt
	notif.SenderId = senderId
	notif.CreatedAt = time.Now()

	// Marshalling the notification to JSON
	notifJson, err := json.Marshal(notif)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Adding the notification to the receiver's sorted set in Redis
	score := float64(time.Now().Unix())
	receiverKey := fmt.Sprintf("notifications:%s", notif.ReceiverId)
	err = a.Redis.ZAdd(context.Request.Context(), receiverKey, goredis.Z{Score: score, Member: notifJson}).Err()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Publish the notification to the receiver's channel
	formattedMessage := "data: " + string(notifJson) + "\n\n"
	err = a.Redis.Publish(context.Request.Context(), receiverKey, formattedMessage).Err()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error publishing notification"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (a *APIEnv) PostNotificationFromEvent(context *gin.Context, notif models.Notification) {
	// Marshalling the notification to JSON
	notifJson, err := json.Marshal(notif)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Adding the notification to the receiver's sorted set in Redis
	score := float64(time.Now().Unix())
	receiverKey := fmt.Sprintf("notifications:%s", notif.ReceiverId)
	err = a.Redis.ZAdd(context.Request.Context(), receiverKey, goredis.Z{Score: score, Member: notifJson}).Err()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Publish the notification to the receiver's channel
	formattedMessage := "data: " + string(notifJson) + "\n\n"
	err = a.Redis.Publish(context.Request.Context(), receiverKey, formattedMessage).Err()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error publishing notification"})
		return
	}
	log.Println("Notification sent")
	log.Println(formattedMessage)
	log.Println(receiverKey)
	// context.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (a *APIEnv) PatchNotification(context *gin.Context) {

}

// CREATE TABLE Notifications (
//     id SERIAL PRIMARY KEY,
//     sender_id INT NOT NULL REFERENCES Users(id),
//     receiver_id INT NOT NULL REFERENCES Users(id),
//     type VARCHAR(30) NOT NULL,
//     PostId INT,
//     CommentId INT,
//     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
//     read BOOLEAN NOT NULL DEFAULT FALSE
// );

// func (a *APIEnv) DeleteNotification(context *gin.Context) {
// 	// Extracting user id and notification id from the request
// 	userId := helpers.GetUserIdFromContext(context)
// 	notificationId := context.Param("id")

// 	// Generate the redis key
// 	receiverKey := "notifications:" + userId

// 	// We have to retrieve all notifications and find the one to delete,
// 	// as sorted set in redis doesn't support direct value deletion
// 	results, err := a.Redis.ZRange(context.Request.Context(), receiverKey, 0, -1).Result()
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting notifications"})
// 		return
// 	}

// 	// Delete the notification with the specific id
// 	for _, result := range results {
// 		var notif models.Notification
// 		err = json.Unmarshal([]byte(result), &notif)
// 		if err != nil {
// 			continue
// 		}
// 		if notif.Id == notificationId {
// 			err = a.Redis.ZRem(context.Request.Context(), receiverKey, result).Err()
// 			if err != nil {
// 				context.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting notification"})
// 				return
// 			}
// 			break
// 		}
// 	}

// 	context.JSON(http.StatusOK, gin.H{"status": "success"})
// }
