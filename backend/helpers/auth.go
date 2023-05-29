package helpers

import (
	"strings"
)

func CheckUserPass(username, password string) bool {
	return true
}

func EmptyUserPass(username, password string) bool {
	return strings.Trim(username, " ") == "" || strings.Trim(password, " ") == ""
}
