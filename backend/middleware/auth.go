package middleware

import (
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
func AuthRequired(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if !helpers.IsValidSession(session) {
		log.Println("UserID in session does not match value in Redis")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to retrieve session",
		})
		ctx.Abort()
	}
	userID := session.Get("userID")
	helpers.AddParamsToContext(ctx, helpers.UserIDKey, userID)
	ctx.Next()
}

type AuthContext interface {
}
