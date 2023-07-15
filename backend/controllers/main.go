/*
Contains common functions, structs, and constants used by handler functions.

The response body will contain one of three keys:
- data: Returns the data requested by the client (uses outputData)
- error: Returns an error message (uses outputError)
- message: Returns a message accompanied with status code 200 OK (uses outputMessage)
*/
package controllers

import (
	"errors"

	"cloud.google.com/go/storage"
	goredis "github.com/redis/go-redis/v9"
	"github.com/ryanozx/skillnet/database"
	"gorm.io/gorm"
)

// APIEnv is a wrapper for the shared database instance
type APIEnv struct {
	DB                *gorm.DB
	Redis             *goredis.Client
	PostDBHandler     database.PostDBHandler
	UserDBHandler     database.UserDBHandler
	AuthDBHandler     database.AuthDBHandler
	LikeDBHandler     database.LikeAPIHandler
	GoogleCloud       *storage.Client
	LikesCacheHandler LikesCacheHandler
}

// Creates an APIEnv object from a Database object
func CreateAPIEnv(db *gorm.DB, gc *storage.Client, redis *goredis.Client) *APIEnv {
	apiEnv := &APIEnv{
		DB:          db,
		GoogleCloud: gc,
		Redis:       redis,
	}
	return apiEnv
}

// General
var (
	ErrBadBinding         = errors.New("invalid request")
	ErrCookieSaveFail     = errors.New("cookie failed to save")
	ErrNoValidSession     = errors.New("no valid session")
	ErrSessionClearFailed = errors.New("failed to clear session")
	// ErrTest is a test error that can be used to simulate an unexpected error returned by the database helper functions
	ErrTest = errors.New("test error")
)
