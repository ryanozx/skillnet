package controllers

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	CreateAccountNoCookieErrMsg    = "Account successfully created but cookie not set"
	PasswordEncryptionFailedErrMsg = "Password encryption failed"
	SessionClearFailedErrMsg       = "Failed to clear session"
	UserNotFoundErrMsg             = "User not found"
	UsernameAlreadyExistsErrMsg    = "Username already exists"

	SuccessfulAccountCreationMsg = "Account successfully created and logged in"
	SuccessfulAccountDeleteMsg   = "User successfully deleted"
)

func (a *APIEnv) InitialiseUserHandler() {
	a.UserDBHandler = &database.UserDB{
		DB: a.DB,
	}
}

// Creates a new user (sign up)
func (a *APIEnv) CreateUser(ctx *gin.Context) {
	userCredentials := helpers.ExtractUserCredentials(ctx)

	// If request is badly formatted, return status code 400 Bad Request
	if helpers.IsEmptyUserPass(userCredentials) {
		helpers.OutputError(ctx, http.StatusBadRequest, MissingUserCredentialsErrMsg)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userCredentials.Password), bcrypt.DefaultCost)
	// If password hash cannot be generated, return status code 500 Internal Service Error
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, PasswordEncryptionFailedErrMsg)
		return
	}
	userCredentials.Password = string(hashedPassword)

	user, err := a.UserDBHandler.CreateUser(userCredentials)
	// If username already exists, return status code 409 Status Conflict
	if err == gorm.ErrDuplicatedKey {
		helpers.OutputError(ctx, http.StatusConflict, UsernameAlreadyExistsErrMsg)
		return
	}
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	// Saves session and sets a session cookie on the client's side; if unsuccessful, return
	// with status code 500 Internal Server Error
	if err := helpers.SaveSession(ctx, user); err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, CreateAccountNoCookieErrMsg)
		return
	}
	helpers.OutputMessage(ctx, SuccessfulAccountCreationMsg)
}

// Deletes user - assuming Delete Account is implemented
func (a *APIEnv) DeleteUser(ctx *gin.Context) {
	userID := helpers.GetUserIdFromContext(ctx)
	err := a.UserDBHandler.DeleteUser(userID)
	// If user cannot be found in the database return status code 404 Status Not Found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		helpers.OutputError(ctx, http.StatusNotFound, UserNotFoundErrMsg)
		return
	}
	// If user cannot be deleted for any other reason, return status code 403 Bad Request
	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	// If unable to invalidate the session on the server side, return with a status code
	// 501 Internal Server Error
	session := sessions.Default(ctx)
	session.Clear()
	if err := session.Save(); err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, SessionClearFailedErrMsg)
		return
	}
	helpers.OutputMessage(ctx, SuccessfulAccountDeleteMsg)
}

// Returns user's profile as seen by visitor
func (a *APIEnv) GetProfile(ctx *gin.Context) {
	username := helpers.GetUsernameFromContext(ctx)

	user, err := a.UserDBHandler.GetUserByUsername(username)
	// If cannot find user in database, return status code 404 Not Found
	if err != nil {
		helpers.OutputError(ctx, http.StatusNotFound, UserNotFoundErrMsg)
		return
	}
	profile := user.GetUserView()
	helpers.OutputData(ctx, profile)
}

// Returns user's own profile with private information
func (a *APIEnv) GetSelfProfile(ctx *gin.Context) {
	userID := helpers.GetUserIdFromContext(ctx)
	user, err := a.UserDBHandler.GetUserByID(userID)
	// If cannot find use in database, return status code 404 Not Found
	if err != nil {
		helpers.OutputError(ctx, http.StatusNotFound, UserNotFoundErrMsg)
		return
	}
	helpers.OutputData(ctx, user)
}

// Updates user's profile
func (a *APIEnv) UpdateUser(ctx *gin.Context) {
	userID := helpers.GetUserIdFromContext(ctx)
	var inputUpdate models.User

	// If request is badly formatted, return status code 400 Bad Request
	if err := helpers.BindInput(ctx, &inputUpdate); err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, err := a.UserDBHandler.UpdateUser(&inputUpdate, userID)
	// If user cannot be found in the database return status code 404 Status Not Found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		helpers.OutputError(ctx, http.StatusNotFound, UserNotFoundErrMsg)
		return
	}
	// If user cannot be updated for any other reason, return status code 500 Internal Server Error
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.OutputData(ctx, user)
}

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
