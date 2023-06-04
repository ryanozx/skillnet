package helpers

import (
	"strings"

	"github.com/ryanozx/skillnet/models"
)

func CheckUserPass(user *models.UserCredentials) bool {
	return true
}

func EmptyUserPass(user *models.UserCredentials) bool {
	return strings.Trim(user.Username, " ") == "" || strings.Trim(user.Password, " ") == ""
}

func IsValidSession(sessionID interface{}) bool {
	if sessionID == nil {
		return false
	} else {
		return true
	}
}
