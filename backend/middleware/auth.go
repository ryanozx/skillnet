package middleware

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/helpers"
)

const loginRoute = "/login"

/*
If user does not have a valid session, the user will be automatically
redirected to the login gateway
*/
func AuthRequired(context *gin.Context) {
	session := sessions.Default(context)
	sessionID := session.Get("session_id")
	if !helpers.IsValidSession(sessionID) {
		log.Println("User does not have a valid session")
		context.Redirect(http.StatusPermanentRedirect, loginRoute)
		context.Abort()
	}
	context.Next()
}
