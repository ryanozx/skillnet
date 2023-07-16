package helpers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	CutoffKey = "cutoff"
)

func getParamFromContext(ctx ParamGetter, key string) string {
	return ctx.Param(key)
}

type ParamGetter interface {
	Param(string) string
}

// Retrieves username from context; the username is inserted into the context
// by the router when parsing ("/users/:username")
func GetUsernameFromContext(ctx ParamGetter) string {
	const usernameKey = "username"
	username := getParamFromContext(ctx, usernameKey)
	return username
}

func AddParamsToContext(ctx AddParamer, key string, val interface{}) {
	ctx.AddParam(key, fmt.Sprintf("%v", val))
}

type AddParamer interface {
	AddParam(key string, value string)
}

// bindInput is used by the handler functions to bind the JSON in the
// request to some struct
func BindInput(ctx BindJSONer, obj interface{}) error {
	err := ctx.BindJSON(obj)
	return err
}

type BindJSONer interface {
	BindJSON(obj any) error
}

// Serialises data as JSON into response body with status code 200 OK
func OutputData(ctx JSONer, obj any) {
	ctx.JSON(http.StatusOK, gin.H{
		"data": obj,
	})
}

type JSONer interface {
	JSON(code int, obj any)
}

// Adds error message to response body along with error status code
func OutputError(ctx JSONer, statusCode int, err error) {
	log.Println(err.Error())
	ctx.JSON(statusCode, gin.H{
		"error": err.Error(),
	})
}

// Adds message to response body along with status code 200 OK
func OutputMessage(ctx JSONer, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": msg,
	})
}

func CheckUserIsOwner(obj UserIDGetter, userID string) error {
	if obj.GetUserID() != userID {
		return ErrNotOwner
	}
	return nil
}

type UserIDGetter interface {
	GetUserID() string
}

var ErrNotOwner = errors.New("unauthorised action")

func getUnsignedValFromContext(ctx ParamGetter, key string) (uint, error) {
	valStr := getParamFromContext(ctx, key)
	val, err := strconv.ParseUint(valStr, 10, 64)
	return uint(val), err
}

type DefaultQueryer interface {
	DefaultQuery(string, string) string
}

func getUnsignedValFromQuery(ctx DefaultQueryer, key string) (uint, error) {
	valStr := ctx.DefaultQuery(key, "")
	val, err := strconv.ParseUint(valStr, 10, 64)
	return uint(val), err
}

func validateUnsignedOrEmptyQuery(ctx DefaultQueryer, key string) (*NullableUint, error) {
	valStr := ctx.DefaultQuery(key, "")
	return ParseNullableUint(valStr)
}

func GetCutoffFromQuery(ctx DefaultQueryer) (*NullableUint, error) {
	return validateUnsignedOrEmptyQuery(ctx, CutoffKey)
}
