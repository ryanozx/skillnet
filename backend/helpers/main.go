package helpers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getParamFromContext(ctx ParamGetter, key string) string {
	return ctx.Param(key)
}

type ParamGetter interface {
	Param(string) string
}

// Retrieves postID from context; the postID is inserted into the context
// by the router when parsing ("/posts/:id")
func GetPostIdFromContext(ctx ParamGetter) string {
	const postIDKey = "id"
	postID := getParamFromContext(ctx, postIDKey)
	return postID
}

// Retrieves userID from context; will be non-empty in private routes since
// AuthRequired adds userID as a parameter in the context
func GetUserIdFromContext(ctx ParamGetter) string {
	userID := getParamFromContext(ctx, IdKey)
	return userID
}

// Retrieves username from context; the username is inserted into the context
// by the router whe parsing ("/users/:username")
func GetUsernameFromContext(ctx ParamGetter) string {
	const usernameKey = "username"
	username := getParamFromContext(ctx, usernameKey)
	return username
}

func AddParamsToContext(ctx *gin.Context, key string, val interface{}) {
	ctx.AddParam(key, fmt.Sprintf("%v", val))
}

// bindInput is used by the handler functions to bind the JSON in the
// request to some struct
func BindInput(ctx *gin.Context, obj interface{}) error {
	err := ctx.BindJSON(obj)
	return err
}

// Serialises data as JSON into response body with status code 200 OK
func OutputData(ctx *gin.Context, obj any) {
	ctx.JSON(http.StatusOK, gin.H{
		"data": obj,
	})
}

// Adds error message to response body along with error status code
func OutputError(ctx *gin.Context, statusCode int, errorMsg string) {
	ctx.JSON(statusCode, gin.H{
		"error": errorMsg,
	})
}

// Adds message to response body along with status code 200 OK
func OutputMessage(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": msg,
	})
}
