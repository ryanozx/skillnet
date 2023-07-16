package helpers

const (
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
