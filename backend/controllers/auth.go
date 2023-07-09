/*
Contains controllers for authentication.
*/
package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
)

// Messages
const (
	GetLoginOKMsg       = "OK"
	LoginSuccessfulMsg  = "Logged in"
	SuccessfulLogoutMsg = "Logged out successfully"
)

// Errors
var (
	ErrAlreadyLoggedIn          = errors.New("already logged in")
	ErrIncorrectUserCredentials = errors.New("incorrect username or password")
	ErrMissingUserCredentials   = errors.New("missing username or password")
)

func (a *APIEnv) InitialiseAuthHandler() {
	a.AuthDBHandler = &database.UserDB{
		DB: a.DB,
	}
}

// If user already has a valid sessionID, the user is redirected, otherwise
// proceed on to login page (200 status code returned)
func (a *APIEnv) GetLogin(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if helpers.IsValidSession(session) {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrAlreadyLoggedIn)
		return
	}
	helpers.OutputMessage(ctx, GetLoginOKMsg)
}

// Handles user login
func (a *APIEnv) PostLogin(ctx *gin.Context) {
	session := sessions.Default(ctx)

	// If user is already logged in (there is an active session), return with status code 400 Bad Request
	if helpers.IsValidSession(session) {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrAlreadyLoggedIn)
		return
	}

	// If request is badly formatted, do not process this request, return with status code 400 Bad Request
	userCredentials := helpers.ExtractUserCredentials(ctx)
	if helpers.IsEmptyUserPass(userCredentials) {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrMissingUserCredentials)
		return
	}

	// If username in request does not refer to a valid user in the database, return with
	// status code 401 Unauthorised
	dbUser, err := a.AuthDBHandler.GetUserByUsername(userCredentials.Username)
	if err != nil {
		helpers.OutputError(ctx, http.StatusUnauthorized, ErrIncorrectUserCredentials)
		return
	}

	// If hashed password does not match hash in database, return with status code 401
	// Unauthorised
	if err := helpers.CheckHashEqualsPassword(dbUser.Password, userCredentials.Password); err != nil {
		helpers.OutputError(ctx, http.StatusUnauthorized, ErrIncorrectUserCredentials)
		return
	}

	// Saves session and sets a session cookie on the client's side; if unsuccessful, return
	// with status code 500 Internal Server Error
	if err := helpers.SaveSession(ctx, dbUser); err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrCookieSaveFail)
		return
	}

	// Login successful
	helpers.OutputMessage(ctx, LoginSuccessfulMsg)
}

// Handles user logout
func (a *APIEnv) PostLogout(ctx *gin.Context) {
	session := sessions.Default(ctx)

	// There is no valid session to log out from - return with a status code 500
	// Bad Request
	if !helpers.IsValidSession(session) {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrNoValidSession)
		return
	}

	session.Clear()
	// If unable to invalidate the session on the server side, return with a status code
	// 501 Internal Server Error
	if err := session.Save(); err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrSessionClearFailed)
		return
	}
	// Session invalidated, successful logout
	helpers.OutputMessage(ctx, SuccessfulLogoutMsg)
}
