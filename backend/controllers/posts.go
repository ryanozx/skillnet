package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
)

const postNotFoundMessage = "Post not found"

type APIEnv struct {
	DB *gorm.DB
}

func (a *APIEnv) GetPosts(context *gin.Context) {
	posts, err := database.GetPosts(a.DB)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, posts)
}

func (a *APIEnv) CreatePost(context *gin.Context) {
	var newPost models.Post
	if err := bindInput(context, &newPost); err != nil {
		return
	}
	newPost.UserID = uuid.MustParse(context.Params.ByName("userID"))
	post, err := database.CreatePost(a.DB, &newPost)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, post)
}

func bindInput(context *gin.Context, obj any) error {
	if bindErr := context.ShouldBindJSON(obj); bindErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return bindErr
	}
	return nil
}

func (a *APIEnv) GetPostByID(context *gin.Context) {
	postID := context.Param("id")
	post, err := database.GetPostByID(a.DB, postID)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, post)
}

func (a *APIEnv) UpdatePost(context *gin.Context) {
	postID := context.Param("id")
	var inputUpdate models.Post
	if bindErr := bindInput(context, &inputUpdate); bindErr != nil {
		return
	}
	post, err := database.UpdatePost(a.DB, &inputUpdate, postID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		context.JSON(http.StatusNotFound, gin.H{"error": postNotFoundMessage})
		return
	} else if errors.Is(err, database.ErrNotOwner) {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	} else if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, post)
}

func (a *APIEnv) DeletePost(context *gin.Context) {
	postID := context.Param("id")
	userID := context.Param("userID")
	err := database.DeletePost(a.DB, postID, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		context.JSON(http.StatusNotFound, gin.H{"error": postNotFoundMessage})
		return
	} else if errors.Is(err, database.ErrNotOwner) {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	} else if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, nil)
}
