package helpers

func GenerateLikeID(userID string, postID string) string {
	return userID + postID
}
