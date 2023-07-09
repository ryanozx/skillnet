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
	DB            *gorm.DB
	Redis         *goredis.Client
	PostDBHandler database.PostDBHandler
	UserDBHandler database.UserDBHandler
	AuthDBHandler database.AuthDBHandler
	LikeDBHandler database.LikeAPIHandler
	GoogleCloud   *storage.Client
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

const (
	GetLoginOKMsg                = "OK"
	LoginSuccessfulMsg           = "Logged in"
	PostDeletedMsg               = "Post successfully deleted"
	SuccessfulAccountCreationMsg = "Account successfully created and logged in"
	SuccessfulAccountDeleteMsg   = "User successfully deleted"
	SuccessfulLogoutMsg          = "Logged out successfully"
)

var (
	ErrAlreadyLoggedIn          = errors.New("already logged in")
	ErrBadBinding               = errors.New("invalid request")
	ErrCannotCreatePost         = errors.New("cannot create post")
	ErrCannotDeletePost         = errors.New("cannot delete post")
	ErrCannotUpdatePost         = errors.New("cannot update post")
	ErrCannotUpdateUser         = errors.New("cannot update user")
	ErrCreateAccountNoCookie    = errors.New("account successfully created but cookie not set")
	ErrCookieSaveFail           = errors.New("cookie failed to save")
	ErrIncorrectUserCredentials = errors.New("incorrect username or password")
	ErrMissingUserCredentials   = errors.New("missing username or password")
	ErrMissingSignupCredentials = errors.New("missing username, password, or email")
	ErrNoValidSession           = errors.New("no valid session")
	ErrPasswordEncryptFailed    = errors.New("password encryption failed")
	ErrPostNotFound             = errors.New("post not found")
	ErrSessionClearFailed       = errors.New("failed to clear session")
	ErrUserNotFound             = errors.New("user not found")
	ErrUsernameAlreadyExists    = errors.New("username already exists")
)
