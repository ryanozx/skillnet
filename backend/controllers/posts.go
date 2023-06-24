/*
Contains controllers for Post API.
*/
package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
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
	userID := helpers.GetUserIdFromContext(ctx)
	newPost.UserID = userID

	post, err := a.PostDBHandler.CreatePost(&newPost)

	// If post cannot be created, return status code 500 Internal Service Error
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrCannotCreatePost)
		return
	}
	helpers.OutputData(ctx, post)
}

func (a *APIEnv) DeletePost(ctx *gin.Context) {
	postID := helpers.GetPostIdFromContext(ctx)
	userID := helpers.GetUserIdFromContext(ctx)

	err := a.PostDBHandler.DeletePost(postID, userID)
	// If post cannot be found in the database return status code 404 Status Not Found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		helpers.OutputError(ctx, http.StatusNotFound, ErrPostNotFound)
		return
	}
	// If user is not the owner of the post, return status code 401 Unauthorized
	if errors.Is(err, database.ErrNotOwner) {
		helpers.OutputError(ctx, http.StatusUnauthorized, database.ErrNotOwner)
		return
	}
	// If post cannot be deleted for any other reason, return status code 403 Bad Request
	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrCannotDeletePost)
		return
	}
	helpers.OutputMessage(ctx, PostDeletedMsg)
}

func (a *APIEnv) GetPosts(ctx *gin.Context) {
	cutoff := ctx.DefaultQuery("cutoff", "")
	userID := helpers.GetUserIdFromContext(ctx)
	posts, newCutoff, err := a.PostDBHandler.GetPosts(cutoff, userID)
	log.Println(err)
	// If unable to retrieve posts, return status code 404 Not Found
	if err != nil {
		helpers.OutputError(ctx, http.StatusNotFound, ErrPostNotFound)
		return
	}
	postViewArray := models.PostViewArray{
		Posts:       posts,
		NextPageURL: fmt.Sprintf("%s/auth/posts?cutoff=%d", models.BackendAddress, newCutoff),
	}
	helpers.OutputData(ctx, postViewArray)
}

func (a *APIEnv) GetPostByID(ctx *gin.Context) {
	postID := helpers.GetPostIdFromContext(ctx)
	userID := helpers.GetPostIdFromContext(ctx)
	post, err := a.PostDBHandler.GetPostByID(postID, userID)
	// If unable to retrieve post, return status code 404 Not Found
	if err != nil {
		helpers.OutputError(ctx, http.StatusNotFound, ErrPostNotFound)
		return
	}
	helpers.OutputData(ctx, post)
}

func (a *APIEnv) UpdatePost(ctx *gin.Context) {
	postID := helpers.GetPostIdFromContext(ctx)
	userID := helpers.GetUserIdFromContext(ctx)
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
	// If user is not the owner of the post, return status code 401 Unauthorised
	if errors.Is(err, database.ErrNotOwner) {
		helpers.OutputError(ctx, http.StatusUnauthorized, database.ErrNotOwner)
		return
	}
	// If post cannot be updated for any other reason, return status code 403 Bad Request
	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrCannotUpdatePost)
		return
	}
	helpers.OutputData(ctx, post)
}
