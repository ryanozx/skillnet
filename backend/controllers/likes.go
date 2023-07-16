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
	ErrAlreadyLiked      = errors.New("already liked")
	ErrLikeNotRegistered = errors.New("like not registered")
	ErrUnlikeFailed      = errors.New("failed to unlike")
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
		helpers.OutputError(ctx, http.StatusInternalServerError, err)
		return
	}

	output := models.LikeUpdate{
		Like:      *like,
		LikeCount: newLikeCount,
	}
	helpers.OutputData(ctx, output)
}

func (a *APIEnv) DeleteLike(ctx *gin.Context) {
	userID := helpers.GetUserIDFromContext(ctx)
	postID, err := helpers.GetPostIDFromContext(ctx)

	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrPostNotFound)
		return
	}

	err = a.LikeDBHandler.DeleteLike(userID, postID)

	if err == gorm.ErrRecordNotFound {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrPostNotFound)
		return
	} else if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrUnlikeFailed)
		return
	}

	newLikeCount, err := a.LikesCacheHandler.SetCacheVal(ctx, postID)
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, err)
		return
	}

	output := models.LikeUpdate{
		LikeCount: newLikeCount,
	}
	helpers.OutputData(ctx, output)
}
