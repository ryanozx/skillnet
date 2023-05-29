package middleware

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired(context *gin.Context) {
	session := sessions.Default(context)
	sessionID := session.Get("session_id")
	if sessionID == nil {
		log.Println("User not logged in")
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized, please log in",
		})
		context.Abort()
	}
	context.Next()
}
