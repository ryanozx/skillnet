package helpers

import "fmt"

func GenerateLikeID(userID string, postID uint) string {
	return fmt.Sprintf("%s%v", userID, postID)
}
