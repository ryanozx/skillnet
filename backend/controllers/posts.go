/*
Contains controllers for Post API.
*/
package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
)

// Messages
const (
	PostDeletedMsg = "Post successfully deleted"
)

// Errors
var (
	ErrCannotCreatePost = errors.New("cannot create post")
	ErrCannotDeletePost = errors.New("cannot delete post")
	ErrCannotUpdatePost = errors.New("cannot update post")
	ErrPostNotFound     = errors.New("post not found")
)

func (a *APIEnv) InitialisePostHandler() {
	a.PostDBHandler = &database.PostDB{
		DB: a.DB,
	}
}

func (a *APIEnv) CreatePost(ctx *gin.Context) {
	var newPost models.Post

	// If unable to bind JSON in request to the Post object, return status
	// code 400 Bad Request
	if err := helpers.BindInput(ctx, &newPost); err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrBadBinding)
		return
	}

	// Add userID into the corresponding field in the newPost object so that
	// the client does not have to pass in any userID, and overwrites any userID
	// that a malicious client might have passed in.
	userID := helpers.GetUserIDFromContext(ctx)
	newPost.UserID = userID

	post, err := a.PostDBHandler.CreatePost(&newPost)

	// If post cannot be created, return status code 500 Internal Service Error
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrCannotCreatePost)
		return
	}
	helpers.OutputData(ctx, post.PostView(&models.PostViewParams{UserID: userID}))
}

func (a *APIEnv) DeletePost(ctx *gin.Context) {
	userID := helpers.GetUserIDFromContext(ctx)

	// Ensure that postID is an unsigned integer
	postID, err := helpers.GetPostIDFromContext(ctx)
	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrPostNotFound)
		return
	}

	err = a.PostDBHandler.DeletePost(postID, userID)
	// If post cannot be found in the database return status code 404 Status Not Found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		helpers.OutputError(ctx, http.StatusNotFound, ErrPostNotFound)
		return
	}
	// If user is not the owner of the post, return status code 403 Forbidden
	if errors.Is(err, helpers.ErrNotOwner) {
		helpers.OutputError(ctx, http.StatusForbidden, helpers.ErrNotOwner)
		return
	}
	// If post cannot be deleted for any other reason, return status code 500 Internal Server Error
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrCannotDeletePost)
		return
	}
	helpers.OutputMessage(ctx, PostDeletedMsg)
}

func (a *APIEnv) GetPosts(ctx *gin.Context) {
	userID := helpers.GetUserIDFromContext(ctx)

	// Ensure that cutoff is an unsigned integer or empty
	cutoff, err := helpers.GetCutoffFromQuery(ctx)
	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrBadBinding)
	}

	// Ensure that community ID is an unsigned integer or empty
	communityID, err := helpers.GetCommunityIDFromQuery(ctx)
	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrBadBinding)
	}

	// Ensure that project ID is an unsigned integer or empty
	projectID, err := helpers.GetProjectIDFromQuery(ctx)
	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrBadBinding)
	}

	posts, err := a.PostDBHandler.GetPosts(cutoff, communityID, projectID, userID)
	// If unable to retrieve posts, return status code 404 Not Found
	if err != nil {
		helpers.OutputError(ctx, http.StatusNotFound, ErrPostNotFound)
		return
	}
	var smallestID uint = 0
	var postViews []models.PostView
	// Fill in user details for each post using userID of post creator
	for _, post := range posts {
		smallestID = post.ID
		likeCount, err := a.LikesCacheHandler.GetCacheVal(ctx, post.ID)
		if err != nil {
			continue
		}
		commentCount, err := a.CommentsCacheHandler.GetCacheVal(ctx, post.ID)
		if err != nil {
			continue
		}
		postView := post.PostView(&models.PostViewParams{
			UserID:       userID,
			LikeCount:    likeCount,
			CommentCount: commentCount,
		})
		postViews = append(postViews, *postView)
	}

	additionalURLParams := map[string]interface{}{}
	if !communityID.IsNull() {
		val, _ := communityID.GetValue()
		additionalURLParams[helpers.CommunityIDQueryKey] = val
	} else if !projectID.IsNull() {
		val, _ := projectID.GetValue()
		additionalURLParams[helpers.ProjectIDQueryKey] = val
	}

	nextPageURL := helpers.GeneratePostNextPageURL(models.BackendAddress, smallestID, additionalURLParams)

	postViewArray := models.PostViewArray{
		Posts:       postViews,
		NextPageURL: nextPageURL,
	}
	helpers.OutputData(ctx, postViewArray)
}

func (a *APIEnv) GetPostByID(ctx *gin.Context) {
	userID := helpers.GetUserIDFromContext(ctx)

	// Ensure that postID is an unsigned integer
	postID, err := helpers.GetPostIDFromContext(ctx)
	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrPostNotFound)
		return
	}

	post, err := a.PostDBHandler.GetPostByID(postID, userID)
	// If unable to retrieve post, return status code 404 Not Found
	if err != nil {
		helpers.OutputError(ctx, http.StatusNotFound, ErrPostNotFound)
		return
	}

	likeCount, err := a.LikesCacheHandler.GetCacheVal(ctx, postID)
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, err)
		return
	}

	commentCount, err := a.CommentsCacheHandler.GetCacheVal(ctx, postID)
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, err)
		return
	}

	postView := post.PostView(&models.PostViewParams{
		UserID:       userID,
		LikeCount:    likeCount,
		CommentCount: commentCount,
	})
	helpers.OutputData(ctx, postView)
}

func (a *APIEnv) UpdatePost(ctx *gin.Context) {
	userID := helpers.GetUserIDFromContext(ctx)

	// Ensure that postID is an unsigned integer
	postID, err := helpers.GetPostIDFromContext(ctx)
	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrPostNotFound)
		return
	}

	var inputUpdate models.Post

	// If unable to bind JSON in request to the Post object, return status
	// code 400 Bad Request
	if err := helpers.BindInput(ctx, &inputUpdate); err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrBadBinding)
		return
	}

	post, err := a.PostDBHandler.UpdatePost(&inputUpdate, postID, userID)

	// If post cannot be found in the database, return status code 404 Status Not Found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		helpers.OutputError(ctx, http.StatusNotFound, ErrPostNotFound)
		return
	}
	// If user is not the owner of the post, return status code 403 Forbidden
	if errors.Is(err, helpers.ErrNotOwner) {
		helpers.OutputError(ctx, http.StatusForbidden, helpers.ErrNotOwner)
		return
	}
	// If post cannot be updated for any other reason, return status code 500 Internal Server Error
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrCannotUpdatePost)
		return
	}

	likeCount, err := a.LikesCacheHandler.GetCacheVal(ctx, post.ID)
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, err)
		return
	}

	commentCount, err := a.CommentsCacheHandler.GetCacheVal(ctx, post.ID)
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, err)
		return
	}

	helpers.OutputData(ctx, post.PostView(&models.PostViewParams{
		UserID:       userID,
		LikeCount:    likeCount,
		CommentCount: commentCount,
	}))
}
