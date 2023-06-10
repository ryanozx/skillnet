package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/helpers"
)

/*
If user does not have a valid session, the user will be automatically
redirected to the login gateway
*/
func AuthRequired(context *gin.Context) {
	session := sessions.Default(context)
	if !helpers.IsValidSession(session) {
		log.Println("UserID in session does not match value in Redis")
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to retrieve session",
		})
		context.Abort()
	}
	userID := session.Get("userID")
	context.AddParam("userID", fmt.Sprintf("%v", userID))
	context.Next()
}
