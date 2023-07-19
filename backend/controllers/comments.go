/*
Contains controllers for Post API.
*/
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
	CommentDeletedMsg = "Comment successfully deleted"
)

// Errors
var (
	ErrCannotCreateComment = errors.New("cannot create comment")
	ErrCannotDeleteComment = errors.New("cannot delete comment")
	ErrCannotUpdateComment = errors.New("cannot update comment")
	ErrCommentNotFound     = errors.New("comment not found")
)

func (a *APIEnv) InitialiseCommentHandler(client *redis.Client) {
	a.CommentDBHandler = &database.CommentDB{
		DB: a.DB,
	}
	a.CommentsCacheHandler = &Cache{
		redisDB:   client,
		DBHandler: a.CommentDBHandler,
	}
}

func (a *APIEnv) CreateComment(ctx *gin.Context) {
	var newComment models.Comment

	// If unable to bind JSON in request to the Comment object, return status
	// code 400 Bad Request
	if err := helpers.BindInput(ctx, &newComment); err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrBadBinding)
		return
	}

	// Add userID into the corresponding field in the newComment object so that
	// the client does not have to pass in any userID, and overwrites any userID
	// that a malicious client might have passed in.
	userID := helpers.GetUserIDFromContext(ctx)
	newComment.UserID = userID

	// Ensure that commentID is an unsigned integer
	postID, err := helpers.GetPostIDFromQuery(ctx)

	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrPostNotFound)
		return
	}
	newComment.PostID = postID

	comment, err := a.CommentDBHandler.CreateComment(&newComment)

	// If comment cannot be created, return status code 500 Internal Service Error
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrCannotCreateComment)
		return
	}

	newCommentCount, err := a.CommentsCacheHandler.SetCacheVal(ctx, uint(newComment.PostID))
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, err)
		return
	}

	output := models.CommentUpdate{
		Comment:      *comment.CommentView(userID),
		CommentCount: newCommentCount,
	}

	helpers.OutputData(ctx, output)
}

func (a *APIEnv) DeleteComment(ctx *gin.Context) {
	userID := helpers.GetUserIDFromContext(ctx)
	// Ensure that commentID is an unsigned integer
	commentID, err := helpers.GetCommentIDFromContext(ctx)

	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrCommentNotFound)
		return
	}

	postID, err := a.CommentDBHandler.DeleteComment(commentID, userID)
	// If comment cannot be found in the database, return status code 404 Status Not Found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		helpers.OutputError(ctx, http.StatusNotFound, ErrCommentNotFound)
		return
	}
	// If user is not the owner of the comment, return status code 403 Forbidden
	if errors.Is(err, helpers.ErrNotOwner) {
		helpers.OutputError(ctx, http.StatusForbidden, helpers.ErrNotOwner)
		return
	}
	// If comment cannot be deleted for any other reason, return status code 400 Bad Request
	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrCannotDeleteComment)
		return
	}

	newCommentCount, err := a.CommentsCacheHandler.SetCacheVal(ctx, postID)
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, err)
		return
	}

	output := models.CommentUpdate{
		CommentCount: newCommentCount,
	}

	helpers.OutputData(ctx, output)
}

func (a *APIEnv) GetComments(ctx *gin.Context) {
	userID := helpers.GetUserIDFromContext(ctx)
	// Ensure that cutoff is an unsigned integer or empty
	cutoff, err := helpers.GetCutoffFromQuery(ctx)
	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrBadBinding)
		return
	}

	// Ensure that postID is an unsigned integer
	postID, err := helpers.GetPostIDFromQuery(ctx)
	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrPostNotFound)
		return
	}

	comments, err := a.CommentDBHandler.GetComments(postID, cutoff)
	// If unable to retrieve comments, return status code 404 Not Found
	if err != nil {
		helpers.OutputError(ctx, http.StatusNotFound, ErrCommentNotFound)
		return
	}

	var smallestID uint = 0

	var commentViews []models.CommentView
	// Fill in user details for each comment using userID of comment creator
	for _, comment := range comments {
		smallestID = comment.ID
		commentView := comment.CommentView(userID)
		commentViews = append(commentViews, *commentView)
	}

	commentViewArray := models.CommentViewsArray{
		Comments:    commentViews,
		NextPageURL: helpers.GenerateCommentNextPageURL(models.BackendAddress, postID, smallestID),
	}
	helpers.OutputData(ctx, commentViewArray)
}

func (a *APIEnv) UpdateComment(ctx *gin.Context) {
	userID := helpers.GetUserIDFromContext(ctx)

	// Ensure that commentID is an unsigned integer
	commentID, err := helpers.GetCommentIDFromContext(ctx)

	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrCommentNotFound)
		return
	}

	var inputUpdate models.Comment

	// If unable to bind JSON in request to the Post object, return status
	// code 400 Bad Request
	if err := helpers.BindInput(ctx, &inputUpdate); err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrBadBinding)
		return
	}

	comment, err := a.CommentDBHandler.UpdateComment(&inputUpdate, commentID, userID)

	// If comment cannot be found in the database, return status code 404 Status Not Found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		helpers.OutputError(ctx, http.StatusNotFound, ErrCommentNotFound)
		return
	}
	// If user is not the owner of the post, return status code 403 Forbidden
	if errors.Is(err, helpers.ErrNotOwner) {
		helpers.OutputError(ctx, http.StatusForbidden, helpers.ErrNotOwner)
		return
	}
	// If post cannot be updated for any other reason, return status code 500 Internal Server Error
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrCannotUpdateComment)
		return
	}
	helpers.OutputData(ctx, comment.CommentView(userID))
}
