package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/models"
)

const postNotFoundMessage = "Post not found"

func GetPosts(context *gin.Context) {
	var posts models.PostArray
	database.Database.Find(&posts)
	context.IndentedJSON(http.StatusOK, gin.H{"data": posts})
}

func PostPost(context *gin.Context) {
	var newPostInput models.CreatePostInput
	if bindErr := context.ShouldBindJSON(&newPostInput); bindErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}

	newPostSchema := models.ConvertInputToPostSchema(newPostInput)
	database.Database.Create(&newPostSchema)

	context.JSON(http.StatusOK, gin.H{"data": newPostSchema})
}

func GetPostByID(context *gin.Context) {
	var postSchema models.PostSchema

	if err := database.Database.Where("id = ?", context.Param("id")).First(&postSchema).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": postNotFoundMessage})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": models.ConvertPostSchemaToPost(postSchema)})
}

func UpdatePost(context *gin.Context) {
	var postSchema models.PostSchema

	if err := database.Database.Where("id = ?", context.Param("id")).First(&postSchema).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": postNotFoundMessage})
		return
	}

	var inputUpdate models.UpdatePostInput
	if err := context.ShouldBindJSON(&inputUpdate); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.Database.Model(&postSchema).Updates(inputUpdate)

	context.JSON(http.StatusOK, gin.H{"data": postSchema})
}

func DeletePost(context *gin.Context) {
	var postSchema models.PostSchema
	if err := database.Database.Where("id = ?", context.Param("id")).First(&postSchema).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": postNotFoundMessage})
		return
	}

	database.Database.Delete(&postSchema)

	context.JSON(http.StatusOK, gin.H{"data": true})
}
