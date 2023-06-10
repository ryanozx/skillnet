package helpers

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/models"
	"golang.org/x/crypto/bcrypt"
)

const (
	sessionKey        = "userID"
	RouteIfSuccessful = "/auth/test"
)

func IsEmptyUserPass(user *models.UserCredentials) bool {
	return strings.Trim(user.Username, " ") == "" || strings.Trim(user.Password, " ") == ""
}

func IsValidSession(session sessions.Session) bool {
	userID := session.Get(sessionKey)
	return userID != nil
}

func BindUserCredentials(context *gin.Context) (userCreds *models.UserCredentials, successfulBind bool) {
	userCredentials := extractUserCredentials(context)
	if IsEmptyUserPass(userCredentials) {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Missing username or password"})
		return userCredentials, false
	}
	return userCredentials, true
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

func SaveSession(context *gin.Context, user *models.User) {
	session := sessions.Default(context)
	session.Set("userID", user.ID)
	if sessionSaveErr := session.Save(); sessionSaveErr != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": sessionSaveErr.Error(),
		})
		return
	}
	context.Redirect(http.StatusMovedPermanently, RouteIfSuccessful)
}

func CheckPassword(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
