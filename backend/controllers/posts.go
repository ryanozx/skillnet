package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/models"
)

const postNotFoundMessage = "Post not found"

func GetPosts(context *gin.Context) {
	var posts []models.PostSchema
	if err := database.Database.Find(&posts).Error; err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{"data": posts})
}

func PostPost(context *gin.Context) {
	var newPostInput models.CreatePostInput
	if bindErr := context.ShouldBindJSON(&newPostInput); bindErr != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}

	newPostSchema := models.ConvertInputToPostSchema(newPostInput)
	database.Database.Create(&newPostSchema)

	context.IndentedJSON(http.StatusOK, gin.H{"data": newPostSchema})
}

func GetPostByID(context *gin.Context) {
	var postSchema models.PostSchema

	if err := database.Database.Where("id = ?", context.Param("id")).First(&postSchema).Error; err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": postNotFoundMessage})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{"data": models.ConvertPostSchemaToPost(postSchema)})
}

func UpdatePost(context *gin.Context) {
	var postSchema models.PostSchema

	if err := database.Database.Where("id = ?", context.Param("id")).First(&postSchema).Error; err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": postNotFoundMessage})
		return
	}

	var inputUpdate models.UpdatePostInput
	if err := context.ShouldBindJSON(&inputUpdate); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.Database.Model(&postSchema).Updates(inputUpdate)

	context.IndentedJSON(http.StatusOK, gin.H{"data": postSchema})
}

func DeletePost(context *gin.Context) {
	var postSchema models.PostSchema
	if err := database.Database.Where("id = ?", context.Param("id")).First(&postSchema).Error; err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": postNotFoundMessage})
		return
	}

	database.Database.Delete(&postSchema)

	context.IndentedJSON(http.StatusOK, gin.H{"data": true})
}
