package helpers

const (
	PostPath       = "/posts"
	PostIDKey      = "postid"
	PostIDQueryKey = "post"
)

// Retrieves postID from context; the postID is inserted into the context
// by the router when parsing ("/posts/:postid")
func GetPostIDFromContext(ctx ParamGetter) (uint, error) {
	return getUnsignedValFromContext(ctx, PostIDKey)
}

func GetPostIDFromQuery(ctx DefaultQueryer) (uint, error) {
	return getUnsignedValFromQuery(ctx, PostIDQueryKey)
}

func GeneratePostNextPageURL(backendURL string, newCutoff uint, additionalParams map[string]interface{}) string {
	return generateNextPageURL(backendURL, PostPath, newCutoff, additionalParams)
}
