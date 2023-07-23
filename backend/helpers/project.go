package helpers

const (
	ProjectPath       = "/projects"
	ProjectIDQueryKey = "project"
	ProjectIDKey      = "projectid"
)

func GetProjectIDFromContext(ctx ParamGetter) (uint, error) {
	return getUnsignedValFromContext(ctx, ProjectIDKey)
}

func GetProjectIDFromQuery(ctx DefaultQueryer) (*NullableUint, error) {
	return validateUnsignedOrEmptyQuery(ctx, ProjectIDQueryKey)
}

func GenerateProjectsNextPageURL(backendURL string, newCutoff uint, communityID *NullableUint, username string) string {
	params := make(map[string]interface{})
	if !communityID.IsNull() {
		communityIDVal, _ := communityID.GetValue()
		params[CommunityIDQueryKey] = communityIDVal
	}
	if username != "" {
		params[UsernameQueryKey] = username
	}
	return generateNextPageURL(backendURL, ProjectPath, newCutoff, params)
}
