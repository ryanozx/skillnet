package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
)

// Messages
const (
	PostUnlikedMsg = "Post unliked"
)

// Errors
var (
	ErrAlreadyLiked          = errors.New("already liked")
	ErrLikeCountFailed       = errors.New("unable to retrieve like count")
	ErrLikeNotFound          = errors.New("like not found")
	ErrLikeNotRegistered     = errors.New("like not registered")
	ErrUnlikeFailed          = errors.New("failed to unlike")
	ErrUpdateLikeCountFailed = errors.New("failed to update like count")
)

func (a *APIEnv) InitialiseLikeHandler(client *redis.Client) {
	a.LikeDBHandler = &database.LikeDB{
		DB: a.DB,
	}
	a.LikesCacheHandler = &Cache{
		redisDB:   client,
		DBHandler: a.LikeDBHandler,
	}
}

func (a *APIEnv) PostLike(ctx *gin.Context) {
	userID := helpers.GetUserIDFromContext(ctx)

	// Ensure that postID is an unsigned integer
	postID, err := helpers.GetPostIDFromContext(ctx)

	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrPostNotFound)
		return
	}

	newLike := &models.Like{
		ID:     helpers.GenerateLikeID(userID, postID),
		UserID: userID,
		PostID: postID,
	}

	like, err := a.LikeDBHandler.CreateLike(newLike)

	if err == gorm.ErrDuplicatedKey {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrAlreadyLiked)
		return
	} else if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrLikeNotRegistered)
		return
	}

	newLikeCount, err := a.LikesCacheHandler.SetCacheVal(ctx, postID)
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrUpdateLikeCountFailed)
		return
	}

	notif := helpers.GenerateLikeNotification(&like.User, like.Post.UserID)

	// Even if there is an error in creating the notification server-side,
	// this should not throw an error client-side
	a.NotificationPoster.PostNotificationFromEvent(ctx, notif)

	output := models.LikeUpdate{
		Like:      *like,
		LikeCount: newLikeCount,
	}

	helpers.OutputData(ctx, output)
}

// func (a *APIEnv) CreateLikeNotification(ctx *gin.Context, userID string, postID uint) error {
// 	notif := models.Notification{
// 		SenderId:  userID,
// 		CreatedAt: time.Now(),
// 	}
// 	post, err := a.PostDBHandler.GetPostByID(postID, "")
// 	if err != nil {
// 		return err
// 	}
// 	notif.ReceiverId = post.UserID
// 	username := post.User.Username
// 	notif.Content = username + " liked your post"
// 	return a.PostNotificationFromEvent(ctx, notif)
// }

func (a *APIEnv) DeleteLike(ctx *gin.Context) {
	userID := helpers.GetUserIDFromContext(ctx)
	postID, err := helpers.GetPostIDFromContext(ctx)

	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrPostNotFound)
		return
	}

	err = a.LikeDBHandler.DeleteLike(userID, postID)

	if err == gorm.ErrRecordNotFound {
		helpers.OutputError(ctx, http.StatusNotFound, ErrLikeNotFound)
		return
	} else if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrUnlikeFailed)
		return
	}

	newLikeCount, err := a.LikesCacheHandler.SetCacheVal(ctx, postID)
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrUpdateLikeCountFailed)
		return
	}

	output := models.LikeUpdate{
		LikeCount: newLikeCount,
	}
	helpers.OutputData(ctx, output)
}
