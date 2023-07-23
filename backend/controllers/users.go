package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Errors
const (
	SuccessfulAccountCreationMsg = "Account successfully created and logged in"
	SuccessfulAccountDeleteMsg   = "User successfully deleted"
)

// Messages
var (
	ErrBadEmail                 = errors.New("invalid email")
	ErrBadPassword              = errors.New("password must have at least 8 characters, one uppercase letter, one lowercase letter, one number, and one special character")
	ErrCannotDeleteUser         = errors.New("cannot delete user")
	ErrCannotUpdateUser         = errors.New("cannot update user")
	ErrCreateAccountNoCookie    = errors.New("account successfully created but cookie not set, please login later")
	ErrMissingSignupCredentials = errors.New("missing username, password, or email")
	ErrPasswordEncryptFailed    = errors.New("password encryption failed")
	ErrUserNotFound             = errors.New("user not found")
	ErrUsernameAlreadyExists    = errors.New("username already exists")
)

func (a *APIEnv) InitialiseUserHandler() {
	a.UserDBHandler = &database.UserDB{
		DB: a.DB,
	}
}

// Creates a new user (sign up)
func (a *APIEnv) CreateUser(ctx *gin.Context) {
	userCredentials := helpers.ExtractSignupUserCredentials(ctx)

	// If request is badly formatted, return status code 400 Bad Request
	if helpers.IsSignupUserCredsEmpty(userCredentials) {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrMissingSignupCredentials)
		return
	}

	// If password does not meet requirements, return status code 400 Bad Request
	if !helpers.ValidatePassword(userCredentials.Password) {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrBadPassword)
		return
	}

	// If email is not validated, return status code 400 Bad Request
	if !helpers.ValidateEmail(userCredentials.Email) {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrBadEmail)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userCredentials.Password), bcrypt.DefaultCost)
	// If password hash cannot be generated, return status code 500 Internal Service Error
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrPasswordEncryptFailed)
		return
	}
	userCredentials.Password = string(hashedPassword)

	user, err := a.UserDBHandler.CreateUser(userCredentials)
	// If username already exists, return status code 409 Status Conflict
	if err == gorm.ErrDuplicatedKey {
		helpers.OutputError(ctx, http.StatusConflict, ErrUsernameAlreadyExists)
		return
	}
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Saves session and sets a session cookie on the client's side; if unsuccessful, return
	// with status code 500 Internal Server Error
	if err := helpers.SaveSession(ctx, user); err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrCreateAccountNoCookie)
		return
	}
	helpers.OutputMessage(ctx, SuccessfulAccountCreationMsg)
}

// Deletes user - assuming Delete Account is implemented
func (a *APIEnv) DeleteUser(ctx *gin.Context) {
	userID := helpers.GetUserIDFromContext(ctx)
	err := a.UserDBHandler.DeleteUser(userID)
	// If user cannot be found in the database return status code 404 Status Not Found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		helpers.OutputError(ctx, http.StatusNotFound, ErrUserNotFound)
		return
	}
	// If user cannot be deleted for any other reason, return status code 501 Internal ServerError
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrCannotDeleteUser)
		return
	}
	// If unable to invalidate the session on the server side, return with a status code
	// 501 Internal Server Error
	session := sessions.Default(ctx)
	session.Clear()
	if err := session.Save(); err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrSessionClearFailed)
		return
	}
	helpers.OutputMessage(ctx, SuccessfulAccountDeleteMsg)
}

// Returns user's profile as seen by visitor
func (a *APIEnv) GetProfile(ctx *gin.Context) {
	username := helpers.GetUsernameFromContext(ctx)
	viewerID := helpers.GetUserIDFromContext(ctx)

	user, err := a.UserDBHandler.GetUserByUsername(username)
	// If cannot find user in database, return status code 404 Not Found
	if err != nil {
		helpers.OutputError(ctx, http.StatusNotFound, ErrUserNotFound)
		return
	}
	profile := user.GetUserView(viewerID)
	helpers.OutputData(ctx, profile)
}

// Returns user's own profile with private information
func (a *APIEnv) GetSelfProfile(ctx *gin.Context) {
	userID := helpers.GetUserIDFromContext(ctx)
	user, err := a.UserDBHandler.GetUserByID(userID)
	// If cannot find use in database, return status code 404 Not Found
	if err != nil {
		helpers.OutputError(ctx, http.StatusNotFound, ErrUserNotFound)
		return
	}
	helpers.OutputData(ctx, user)
}

// Updates user's profile
func (a *APIEnv) UpdateUser(ctx *gin.Context) {
	userID := helpers.GetUserIDFromContext(ctx)
	var inputUpdate models.User

	// If request is badly formatted, return status code 400 Bad Request
	if err := helpers.BindInput(ctx, &inputUpdate); err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrBadBinding)
		return
	}

	user, err := a.UserDBHandler.UpdateUser(&inputUpdate, userID)
	// If user cannot be found in the database return status code 404 Status Not Found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		helpers.OutputError(ctx, http.StatusNotFound, ErrUserNotFound)
		return
	}
	// If user cannot be updated for any other reason, return status code 500 Internal Server Error
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrCannotUpdateUser)
		return
	}
	helpers.OutputData(ctx, user)
}
