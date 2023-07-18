package helpers

import "fmt"

const (
	PostIDKey           = "postid"
	PostIDQueryKey      = "post"
	ProjectIDQueryKey   = "project"
	CommunityIDQueryKey = "community"
)

// Retrieves postID from context; the postID is inserted into the context
// by the router when parsing ("/posts/:postid")
func GetPostIDFromContext(ctx ParamGetter) (uint, error) {
	return getUnsignedValFromContext(ctx, PostIDKey)
}

func GetPostIDFromQuery(ctx DefaultQueryer) (uint, error) {
	return getUnsignedValFromQuery(ctx, PostIDQueryKey)
}

func GetProjectIDFromQuery(ctx DefaultQueryer) (*NullableUint, error) {
	return validateUnsignedOrEmptyQuery(ctx, ProjectIDQueryKey)
}

func GetCommunityIDFromQuery(ctx DefaultQueryer) (*NullableUint, error) {
	return validateUnsignedOrEmptyQuery(ctx, CommunityIDQueryKey)
}

func GeneratePostNextPageURL(backendURL string, newCutoff uint, additionalParams map[string]interface{}) string {
	nextPageURL := fmt.Sprintf("%s/auth/posts?%s=%d", backendURL, CutoffKey, newCutoff)
	for paramKey, paramVal := range additionalParams {
		nextPageURL += fmt.Sprintf("&%s=%v", paramKey, paramVal)
	}
	return nextPageURL
}
