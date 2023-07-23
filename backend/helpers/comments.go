package helpers

const (
	CommentPath  = "/comments"
	CommentIDKey = "commentid"
)

// Retrieves commentID from context; the commentID is inserted into the context
// by the router when parsing ("/comments/:commentid")
func GetCommentIDFromContext(ctx ParamGetter) (uint, error) {
	return getUnsignedValFromContext(ctx, CommentIDKey)
}

func GenerateCommentNextPageURL(backendURL string, postID uint, newCutoff uint) string {
	return generateNextPageURL(backendURL, CommentPath, newCutoff, map[string]interface{}{
		PostIDQueryKey: postID,
	})
}
