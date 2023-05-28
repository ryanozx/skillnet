package controllers

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ryanozx/skillnet/helpers"
)

func LoginGetHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		session := sessions.Default(context)
		sessionID := session.Get("session_id")
		if sessionID != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{
				"message": "Please logout first",
			})
			return
		}
		context.IndentedJSON(http.StatusOK, gin.H{})
	}
}

func LoginPostHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		session := sessions.Default(context)
		sessionID := session.Get("session_id")
		if sessionID != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{
				"message": "Please logout first",
			})
			return
		}

		username := context.PostForm("username")
		password := context.PostForm("password")

		if helpers.EmptyUserPass(username, password) {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing username or password"})
			return
		}

		if !helpers.CheckUserPass(username, password) {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Incorrect username or password"})
			return
		}

		session.Set("session_id", uuid.NewString())
		if sessionErr := session.Save(); sessionErr != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to save session",
			})
			return
		}

		context.Redirect(http.StatusMovedPermanently, "/auth/test")
	}
}

func LogoutGetHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		session := sessions.Default(context)
		sessionID := session.Get("session_id")
		if sessionID == nil {
			log.Println("Invalid session token")
			return
		}
		session.Clear()
		if err := session.Save(); err != nil {
			log.Println("Failed to save session:", err)
			return
		}

		context.IndentedJSON(http.StatusOK, gin.H{
			"message": "Logged out successfully",
		})
	}
}
