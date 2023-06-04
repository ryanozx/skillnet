package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"golang.org/x/crypto/bcrypt"
)

const sessionKey = "session_id"
const routeIfSuccessful = "/auth/test"

/*
If user already has a valid sessionID, the user is redirected
to "/posts", otherwise proceed on to login page
*/
func GetLogin(context *gin.Context) {
	session := sessions.Default(context)
	sessionID := session.Get(sessionKey)
	if helpers.IsValidSession(sessionID) {
		context.Redirect(http.StatusPermanentRedirect, routeIfSuccessful)
		return
	}
	context.JSON(http.StatusOK, gin.H{})
}

func PostLogin(context *gin.Context) {
	session := sessions.Default(context)
	sessionID := session.Get(sessionKey)

	if helpers.IsValidSession(sessionID) {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Please logout first",
		})
		return
	}

	userCredentials := extractUserCredentials(context)

	if helpers.EmptyUserPass(userCredentials) {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Missing username or password"})
		return
	}
	if !helpers.CheckUserPass(userCredentials) {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Incorrect username or password"})
		return
	}

	saveSession(context)
}

func GetLogout(context *gin.Context) {
	session := sessions.Default(context)
	sessionID := session.Get(sessionKey)
	if !helpers.IsValidSession(sessionID) {
		log.Println("Invalid session token")
		return
	}
	session.Clear()
	if sessionErr := session.Save(); sessionErr != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to clear session",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

func CreateUser(context *gin.Context) {
	user := extractUserCredentials(context)

	if helpers.EmptyUserPass(user) {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Missing username or password."})
		return
	}

	hashedPassword, passwordErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if passwordErr != nil {
		fmt.Println(passwordErr)
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Password Encryption failed",
		})
		return
	}

	user.Password = string(hashedPassword)
	userDBEntry := user.ConvertToUser()

	if err := database.Database.Where("username = ?", userDBEntry.Username).First(&userDBEntry); err != nil {
		context.JSON(http.StatusConflict, gin.H{
			"message": "Username already exists",
		})
		return
	}

	database.Database.Create(userDBEntry)
	saveSession(context)
}

func extractUserCredentials(context *gin.Context) *models.UserCredentials {
	const usernameKey = "username"
	const passwordKey = "password"
	username := context.PostForm(usernameKey)
	password := context.PostForm(passwordKey)
	return &models.UserCredentials{
		Username: username,
		Password: password,
	}
}

func saveSession(context *gin.Context) {
	session := sessions.Default(context)
	newSessionID := uuid.NewString()
	session.Set(sessionKey, newSessionID)
	// TODO: Register new session ID on redis
	if sessionSaveErr := session.Save(); sessionSaveErr != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to clear session",
		})
		return
	}
	context.Redirect(http.StatusMovedPermanently, routeIfSuccessful)
}
