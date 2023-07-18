package helpers

import (
	"log"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/models"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserIDKey         = "userID"
	RouteIfSuccessful = "/posts"
)

func IsEmptyUserPass(user *models.UserCredentials) bool {
	return strings.Trim(user.Username, " ") == "" || strings.Trim(user.Password, " ") == ""
}

func IsSignupUserCredsEmpty(user *models.SignupUserCredentials) bool {
	return IsEmptyUserPass(&user.UserCredentials) || strings.Trim(user.Email, " ") == ""
}

func IsValidSession(session SessionGetter) bool {
	userID := session.Get(UserIDKey)
	return userID != nil
}

type SessionGetter interface {
	Get(interface{}) interface{}
}

func ExtractUserCredentials(ctx postFormer) *models.UserCredentials {
	const usernameKey = "username"
	const passwordKey = "password"
	username := ctx.PostForm(usernameKey)
	password := ctx.PostForm(passwordKey)
	return &models.UserCredentials{
		Username: username,
		Password: password,
	}
}

type postFormer interface {
	PostForm(string) string
}

func ExtractSignupUserCredentials(ctx postFormer) *models.SignupUserCredentials {
	const emailKey = "email"
	email := ctx.PostForm(emailKey)
	userCreds := ExtractUserCredentials(ctx)
	return &models.SignupUserCredentials{
		UserCredentials: *userCreds,
		Email:           email,
	}
}

func SaveSession(ctx *gin.Context, user *models.User) error {
	session := sessions.Default(ctx)
	session.Set(UserIDKey, user.ID)
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

// Retrieves userID from context; will be non-empty in private routes since
// AuthRequired adds userID as a parameter in the context
func GetUserIDFromContext(ctx ParamGetter) string {
	userID := getParamFromContext(ctx, UserIDKey)
	return userID
}
