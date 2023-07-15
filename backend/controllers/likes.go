package controllers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

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
	ErrLikeNotRegistered     = errors.New("like not registered")
	ErrUnlikeFailed          = errors.New("failed to unlike")
	ErrUpdateLikeCountFailed = errors.New("failed to update like count")
)

func (a *APIEnv) InitialiseLikeHandler(client *redis.Client) {
	a.LikeDBHandler = &database.LikeDB{
		DB: a.DB,
	}
	a.LikesCacheHandler = &LikesCache{
		redisDB:   client,
		DBHandler: a.LikeDBHandler,
	}
}

func (a *APIEnv) PostLike(ctx *gin.Context) {
	postId := helpers.GetPostIdFromContext(ctx)
	userId := helpers.GetUserIdFromContext(ctx)

	// Ensure that postID is an unsigned integer
	postIdNum, err := strconv.ParseUint(postId, 10, 64)

	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrPostNotFound)
		return
	}

	newLike := &models.Like{
		ID:     helpers.GenerateLikeID(userId, postId),
		UserID: userId,
		PostID: uint(postIdNum),
	}

	like, err := a.LikeDBHandler.CreateLike(newLike)

	if err == gorm.ErrDuplicatedKey {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrAlreadyLiked)
		return
	} else if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrLikeNotRegistered)
		return
	}

	newLikeCount, err := a.LikesCacheHandler.SetCacheCount(ctx, postId)
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, err)
		return
	}

	output := models.LikeUpdate{
		Like:      *like,
		LikeCount: newLikeCount,
	}

	a.CreateLikeNotification(ctx)

	helpers.OutputData(ctx, output)
}

func (a *APIEnv) CreateLikeNotification(ctx *gin.Context) {
	postId := helpers.GetPostIdFromContext(ctx)
	userId := helpers.GetUserIdFromContext(ctx)

	notif := models.Notification{
		SenderId:  userId,
		CreatedAt: time.Now(),
	}
	post, err := a.PostDBHandler.GetPostByPostID(postId)
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, err)
		return
	}
	notif.ReceiverId = post.UserID
	username := post.User.Username
	notif.Content = username + " liked your post"
	a.PostNotificationFromEvent(ctx, notif)
}

func (a *APIEnv) DeleteLike(ctx *gin.Context) {
	postId := helpers.GetPostIdFromContext(ctx)
	userId := helpers.GetUserIdFromContext(ctx)

	// Ensure that postID is an unsigned integer
	_, err := strconv.ParseUint(postId, 10, 64)

	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrPostNotFound)
		return
	}

	err = a.LikeDBHandler.DeleteLike(userId, postId)

	if err == gorm.ErrRecordNotFound {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrPostNotFound)
		return
	} else if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrUnlikeFailed)
		return
	}

	newLikeCount, err := a.LikesCacheHandler.SetCacheCount(ctx, postId)
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, err)
		return
	}

	output := models.LikeUpdate{
		LikeCount: newLikeCount,
	}
	helpers.OutputData(ctx, output)
}

type LikesCacheHandler interface {
	GetCacheCount(context.Context, string) (uint64, error)
	SetCacheCount(context.Context, string) (uint64, error)
}

type LikesCache struct {
	redisDB   *redis.Client
	DBHandler database.LikeDBCountGetter
}

func (c *LikesCache) GetCacheCount(ctx context.Context, postID string) (uint64, error) {
	likeVal, err := c.redisDB.Get(ctx, postID).Result()
	if err == redis.Nil {
		return c.SetCacheCount(ctx, postID)
	}
	return strconv.ParseUint(likeVal, 10, 64)
}

func (c *LikesCache) SetCacheCount(ctx context.Context, postID string) (uint64, error) {
	newLikeCount, err := c.DBHandler.GetLikeCount(postID)
	if err != nil {
		return newLikeCount, ErrLikeCountFailed
	}
	err = c.redisDB.Set(ctx, postID, newLikeCount, 0).Err()
	return newLikeCount, err
}
