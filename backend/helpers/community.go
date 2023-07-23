package helpers

import "fmt"

const (
	CommunityPath       = "/community"
	CommunityIDQueryKey = "community"
	CommunityNameKey    = "communityid"
)

func GetCommunityNameFromContext(ctx ParamGetter) string {
	return getParamFromContext(ctx, CommunityNameKey)
}

func GetCommunityIDFromQuery(ctx DefaultQueryer) (*NullableUint, error) {
	return validateUnsignedOrEmptyQuery(ctx, CommunityIDQueryKey)
}

func GenerateCommunitiesNextPageURL(backendURL string, newCutoff uint) string {
	nextPageURL := fmt.Sprintf("%s/auth%s?%s=%d", backendURL, CommunityPath, CutoffKey, newCutoff)
	return nextPageURL
}
