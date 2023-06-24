package controllers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/helpers"
)

func (a *APIEnv) PostUserPicture(context *gin.Context) {
	userID := helpers.GetUserIdFromContext(context)
	// username := helpers.GetUsernameFromContext(context)
	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer openedFile.Close()

	bucket := a.GoogleCloud.Bucket("skillnet-profile-pictures")
	ctx := context.Request.Context()
	fileName := userID + "-pfp.jpeg"
	writer := bucket.Object(fileName).NewWriter(ctx)

	_, err = io.Copy(writer, openedFile)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := writer.Close(); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	attrs := writer.Attrs()
	url := attrs.MediaLink

	context.JSON(http.StatusOK, gin.H{"url": url})
}
