package helpers

import "github.com/gin-gonic/gin"

const (
	UsernameQueryKey = "username"
	UsernameKey      = "username"
)

func GetUsernameFromQuery(ctx *gin.Context) string {
	usernameStr := ctx.DefaultQuery(UsernameQueryKey, "")
	return usernameStr
}

// Retrieves username from context; the username is inserted into the context
// by the router when parsing ("/users/:username")
func GetUsernameFromContext(ctx ParamGetter) string {
	username := getParamFromContext(ctx, UsernameKey)
	return username
}
