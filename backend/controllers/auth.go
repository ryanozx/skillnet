/*
Contains controllers for authentication.
*/
package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
)

const (
	MissingUserCredentialsErrMsg  = "Missing username or password"
	IncorrectUserCredentialErrMsg = "Incorrect username or password"
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
		ctx.Redirect(http.StatusPermanentRedirect, helpers.RouteIfSuccessful)
		return
	}
	helpers.OutputMessage(ctx, "OK")
}

// Handles user login
func (a *APIEnv) PostLogin(ctx *gin.Context) {
	session := sessions.Default(ctx)

	// If user is already logged in (there is an active session), return with status code 400 Bad Request
	if helpers.IsValidSession(session) {
		helpers.OutputError(ctx, http.StatusBadRequest, "Already logged in")
		return
	}

	// If request is badly formatted, do not process this request, return with status code 400 Bad Request
	userCredentials := helpers.ExtractUserCredentials(ctx)
	if helpers.IsEmptyUserPass(userCredentials) {
		helpers.OutputError(ctx, http.StatusBadRequest, MissingUserCredentialsErrMsg)
		return
	}

	// If username in request does not refer to a valid user in the database, return with
	// status code 401 Unauthorised
	dbUser, err := a.AuthDBHandler.GetUserByUsername(userCredentials.Username)
	if err != nil {
		helpers.OutputError(ctx, http.StatusUnauthorized, IncorrectUserCredentialErrMsg)
		return
	}

	// If hashed password does not match hash in database, return with status code 401
	// Unauthorised
	if err := helpers.CheckHashEqualsPassword(dbUser.Password, userCredentials.Password); err != nil {
		helpers.OutputError(ctx, http.StatusUnauthorized, IncorrectUserCredentialErrMsg)
		return
	}

	// Saves session and sets a session cookie on the client's side; if unsuccessful, return
	// with status code 500 Internal Server Error
	if err := helpers.SaveSession(ctx, dbUser); err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Login successful
	helpers.OutputMessage(ctx, "Logged in")
}

// Handles user logout
func (a *APIEnv) GetLogout(ctx *gin.Context) {
	session := sessions.Default(ctx)

	// There is no valid session to log out from - return with a status code 500
	// Bad Request
	if !helpers.IsValidSession(session) {
		helpers.OutputError(ctx, http.StatusBadRequest, "No valid session")
		return
	}

	session.Clear()
	// If unable to invalidate the session on the server side, return with a status code
	// 501 Internal Server Error
	if err := session.Save(); err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, "Failed to clear session")
		return
	}
	// Session invalidated, successful logout
	helpers.OutputMessage(ctx, "Logged out successfully")
}
