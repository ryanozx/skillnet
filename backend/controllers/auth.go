package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

/*
If user already has a valid sessionID, the user is redirected, otherwise proceed on to login page
*/
func (a *APIEnv) GetLogin(context *gin.Context) {
	session := sessions.Default(context)
	if helpers.IsValidSession(session) {
		context.Redirect(http.StatusPermanentRedirect, helpers.RouteIfSuccessful)
		return
	}
	context.JSON(http.StatusOK, gin.H{})
}

func (a *APIEnv) PostLogin(context *gin.Context) {
	session := sessions.Default(context)

	if helpers.IsValidSession(session) {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Please logout first",
		})
		return
	}

	userCredentials, bindSucceed := helpers.BindUserCredentials(context)
	if !bindSucceed {
		return
	}
	dbUser, err := database.GetUserByUsername(a.DB, userCredentials.Username)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Incorrect username or password"})
		return
	}
	if err := helpers.CheckPassword(dbUser.Password, userCredentials.Password); err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Incorrect username or password"})
		return
	}

	helpers.SaveSession(context, dbUser)
}

func (a *APIEnv) GetLogout(context *gin.Context) {
	session := sessions.Default(context)
	if !helpers.IsValidSession(session) {
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

func (a *APIEnv) CreateUser(context *gin.Context) {
	userCredentials, bindSucceed := helpers.BindUserCredentials(context)
	if !bindSucceed {
		return
	}
	hashedPassword, passwordErr := bcrypt.GenerateFromPassword([]byte(userCredentials.Password), bcrypt.DefaultCost)
	if passwordErr != nil {
		fmt.Println(passwordErr)
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Password Encryption failed",
		})
		return
	}

	userCredentials.Password = string(hashedPassword)

	user, err := database.CreateUser(a.DB, userCredentials)
	if err != nil {
		context.JSON(http.StatusConflict, gin.H{
			"message": err.Error(),
		})
		return
	}
	helpers.SaveSession(context, user)
}

func (a *APIEnv) UpdateUser(context *gin.Context) {
	userID := context.Param("userID")
	var inputUpdate models.User
	if bindErr := bindInput(context, &inputUpdate); bindErr != nil {
		return
	}
	user, err := database.UpdateUser(a.DB, &inputUpdate, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, user)
}

func (a *APIEnv) GetProfile(context *gin.Context) {
	username := context.Param("username")
	user, err := database.GetUserByUsername(a.DB, username)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}
	profile := user.Profile
	context.JSON(http.StatusOK, profile)
}

func (a *APIEnv) GetProfileSettings(context *gin.Context) {
	userID := context.Param("userID")
	user, err := database.GetUserByID(a.DB, userID)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "Profile not found",
		})
		return
	}
	context.JSON(http.StatusOK, user)
}
