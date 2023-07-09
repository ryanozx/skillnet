package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
)

const (
	PostUnlikedMsg = "Post unliked"
)

var (
	ErrAlreadyLiked      = errors.New("already liked")
	ErrLikeNotRegistered = errors.New("like not registered")
	ErrUnlikeFailed      = errors.New("failed to unlike")
)

func (a *APIEnv) InitialiseLikeHandler() {
	a.LikeDBHandler = &database.LikeDB{
		DB: a.DB,
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
	helpers.OutputData(ctx, like)
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
	helpers.OutputMessage(ctx, PostUnlikedMsg)
}
