package helpers

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/models"
	"golang.org/x/crypto/bcrypt"
)

const (
	IdKey             = "userID"
	RouteIfSuccessful = "/posts"
)

func IsEmptyUserPass(user *models.UserCredentials) bool {
	return strings.Trim(user.Username, " ") == "" || strings.Trim(user.Password, " ") == ""
}

func IsValidSession(session SessionGetter) bool {
	userID := session.Get(IdKey)
	fmt.Println(userID)
	return userID != nil
}

type SessionGetter interface {
	Get(interface{}) interface{}
}

func ExtractUserCredentials(ctx *gin.Context) *models.UserCredentials {
	const usernameKey = "username"
	const passwordKey = "password"
	username := ctx.PostForm(usernameKey)
	password := ctx.PostForm(passwordKey)
	return &models.UserCredentials{
		Username: username,
		Password: password,
	}
}

func SaveSession(ctx *gin.Context, user *models.User) error {
	session := sessions.Default(ctx)
	session.Set(IdKey, user.ID)
	log.Printf("Saving userID: %v", user.ID)
	if err := session.Save(); err != nil {
		return err
	}
	return nil
}

func CheckHashEqualsPassword(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GenerateHashFromPassword(password string) (hash []byte, err error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
