package controllers

import (
	"io"
	"net/http"
	"testing"

	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
)

const (
	testPostID           = 1
	testDiffCutoffPostID = 10
	invalidPostID        = "badpostid"
	negativePostID       = -234
	testProjectID        = 1
	invalidProjectID     = "badprojectid"
)

var (
	defaultPost = models.Post{
		Model: gorm.Model{
			ID: testPostID,
		},
		UserID:  testUserID,
		Content: "Hello world!",
		User:    defaultUser,
	}
	diffCutoffPost = models.Post{
		Model: gorm.Model{
			ID: testDiffCutoffPostID,
		},
		Content: "Hello world!",
		UserID:  testUserID,
		User:    defaultUser,
	}
	newTestPost = models.Post{
		Content: "Hello world!",
	}
)

type PostDBTestHandler struct {
	CreatePostFunc  func(*models.Post) (*models.Post, error)
	DeletePostFunc  func(uint, string) error
	GetPostsFunc    func(*helpers.NullableUint, *helpers.NullableUint, *helpers.NullableUint, string) ([]models.Post, error)
	GetPostByIDFunc func(uint, string) (*models.Post, error)
	UpdatePostFunc  func(*models.Post, uint, string) (*models.Post, error)
}

func (h *PostDBTestHandler) CreatePost(newPost *models.Post) (*models.Post, error) {
	return h.CreatePostFunc(newPost)
}

func (h *PostDBTestHandler) DeletePost(postID uint, userid string) error {
	return h.DeletePostFunc(postID, userid)
}

func (h *PostDBTestHandler) GetPosts(cutoff *helpers.NullableUint, communityID *helpers.NullableUint,
	projectID *helpers.NullableUint, userID string) ([]models.Post, error) {
	return h.GetPostsFunc(cutoff, communityID, projectID, userID)
}

func (h *PostDBTestHandler) GetPostByID(postID uint, userID string) (*models.Post, error) {
	return h.GetPostByIDFunc(postID, userID)
}

func (h *PostDBTestHandler) UpdatePost(post *models.Post, postID uint, userID string) (*models.Post, error) {
	return h.UpdatePostFunc(post, postID, userID)
}

func (h *PostDBTestHandler) SetMockCreatePostFunc(post *models.Post, err error) {
	h.CreatePostFunc = func(newPost *models.Post) (*models.Post, error) {
		return post, err
	}
}

func (h *PostDBTestHandler) SetMockDeletePostFunc(err error) {
	h.DeletePostFunc = func(postID uint, userid string) error {
		return err
	}
}

func (h *PostDBTestHandler) SetMockGetPostsFunc(posts []models.Post, err error) {
	h.GetPostsFunc = func(cutoff *helpers.NullableUint, communityID *helpers.NullableUint,
		projectID *helpers.NullableUint, userID string) ([]models.Post, error) {
		return posts, err
	}
}

func (h *PostDBTestHandler) SetMockGetPostByIDFunc(post *models.Post, err error) {
	h.GetPostByIDFunc = func(postID uint, userID string) (*models.Post, error) {
		return post, err
	}
}

func (h *PostDBTestHandler) SetMockUpdatePostFunc(outputPost *models.Post, err error) {
	h.UpdatePostFunc = func(post *models.Post, postID uint, userID string) (*models.Post, error) {
		return outputPost, err
	}
}

func TestAPIEnv_InitialisePostHandler(t *testing.T) {
	type fields struct {
		DB *gorm.DB
	}
	tests := []struct {
		name          string
		fields        fields
		expectedEmpty bool
	}{
		{
			"Initialise Post DB OK",
			fields{
				DB: &gorm.DB{},
			},
			false,
		},
		{
			"No DB OK",
			fields{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIEnv{
				DB: tt.fields.DB,
			}
			a.InitialisePostHandler()
			if postDB, ok := a.PostDBHandler.(*database.PostDB); ok {
				if tt.expectedEmpty && postDB.DB != nil {
					t.Error("Post DB contains unexpected DB instance")
				} else if !tt.expectedEmpty && postDB.DB != tt.fields.DB {
					t.Error("PostDBHandler not initialised correctly")
				}
			} else {
				t.Error("PostDBHandler is nil!")
			}
		})
	}
}

func TestAPIEnv_CreatePost(t *testing.T) {
	helpers.SetEnvVars(t)
	type args struct {
		ContextParams map[string]interface{}
		QueryParams   map[string]interface{}
		PostData      *models.Post
		PostDBOutput  *models.Post
		PostDBError   error
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.PostView]
	}{
		{
			"Create Post - OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				PostData:     &newTestPost,
				PostDBOutput: &defaultPost,
				PostDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.PostView]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data: defaultPost.PostView(&models.PostViewParams{
					UserID: testUserID,
				}),
			},
		},
		{
			"Create Post - Bad Request",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
			},
			helpers.ExpectedJSONOutput[models.PostView]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrBadBinding,
			},
		},
		{
			"Create Post - Cannot Create",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				PostData:     &newTestPost,
				PostDBOutput: &defaultPost,
				PostDBError:  ErrTest,
			},
			helpers.ExpectedJSONOutput[models.PostView]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrCannotCreatePost,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbTestHandler := &PostDBTestHandler{}
			likesCacheTestHandler := &helpers.TestCache{}
			commentsCacheTestHandler := &helpers.TestCache{}
			a := &APIEnv{
				PostDBHandler:        dbTestHandler,
				LikesCacheHandler:    likesCacheTestHandler,
				CommentsCacheHandler: commentsCacheTestHandler,
			}

			c, w := helpers.CreateTestContextAndRecorder()

			for paramKey, paramVal := range tt.args.ContextParams {
				helpers.AddParamsToContext(c, paramKey, paramVal)
			}

			if tt.args.PostData != nil {
				req, err := helpers.GenerateHttpJSONRequest(http.MethodPost, tt.args.PostData)
				if err != nil {
					t.Error(err)
				}

				for paramKey, paramVal := range tt.args.QueryParams {
					helpers.AddParamsToQuery(req, paramKey, paramVal)
				}

				c.Request = req
			}

			dbTestHandler.SetMockCreatePostFunc(tt.args.PostDBOutput, tt.args.PostDBError)
			a.CreatePost(c)

			b, _ := io.ReadAll(w.Body)
			if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(tt.expected.StatusCode, w.Code); !isEqual {
				t.Error(errStr)
			}

			m, err := helpers.ParseJSONString(b)
			if err != nil {
				t.Error(err)
			}

			if errStr, isEqual := helpers.CheckExpectedJSONEqualsActual(m, tt.expected); !isEqual {
				t.Error(errStr)
			}
		})
	}
}

func TestAPIEnv_DeletePost(t *testing.T) {
	type args struct {
		ContextParams map[string]interface{}
		PostDBError   error
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.Post]
	}{
		{
			"Delete Posts OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: testPostID,
				},
				PostDBError: nil,
			},
			helpers.ExpectedJSONOutput[models.Post]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedMessage,
				Message:    PostDeletedMsg,
			},
		},
		{
			"Delete Posts - Not found",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: testPostID,
				},
				PostDBError: gorm.ErrRecordNotFound,
			},
			helpers.ExpectedJSONOutput[models.Post]{
				StatusCode: http.StatusNotFound,
				JSONType:   helpers.ExpectedError,
				Error:      ErrPostNotFound,
			},
		},
		{
			"Delete Posts - Not owner",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: testPostID,
				},
				PostDBError: helpers.ErrNotOwner,
			},
			helpers.ExpectedJSONOutput[models.Post]{
				StatusCode: http.StatusForbidden,
				JSONType:   helpers.ExpectedError,
				Error:      helpers.ErrNotOwner,
			},
		},
		{
			"Delete Posts - Cannot delete",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: testPostID,
				},
				PostDBError: ErrTest,
			},
			helpers.ExpectedJSONOutput[models.Post]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrCannotDeletePost,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbTestHandler := &PostDBTestHandler{}
			a := &APIEnv{
				PostDBHandler: dbTestHandler,
			}

			c, w := helpers.CreateTestContextAndRecorder()

			for paramKey, paramVal := range tt.args.ContextParams {
				helpers.AddParamsToContext(c, paramKey, paramVal)
			}

			req, err := helpers.GenerateHttpJSONRequest(http.MethodDelete, nil)
			if err != nil {
				t.Error(err)
			}

			c.Request = req

			dbTestHandler.SetMockDeletePostFunc(tt.args.PostDBError)
			a.DeletePost(c)

			b, _ := io.ReadAll(w.Body)
			if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(tt.expected.StatusCode, w.Code); !isEqual {
				t.Error(errStr)
			}

			m, err := helpers.ParseJSONString(b)
			if err != nil {
				t.Error(err)
			}

			if errStr, isEqual := helpers.CheckExpectedJSONEqualsActual(m, tt.expected); !isEqual {
				t.Error(errStr)
			}
		})
	}
}

func TestAPIEnv_GetPosts(t *testing.T) {
	helpers.SetEnvVars(t)

	type args struct {
		ContextParams      map[string]interface{}
		QueryParams        map[string]interface{}
		PostDBOutput       []models.Post
		PostDBError        error
		LikesCacheVal      uint64
		LikesCacheError    error
		CommentsCacheVal   uint64
		CommentsCacheError error
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.PostViewArray]
	}{
		{
			"Get posts OK - global feed",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				PostDBOutput:       []models.Post{defaultPost},
				PostDBError:        nil,
				LikesCacheVal:      1,
				LikesCacheError:    nil,
				CommentsCacheVal:   2,
				CommentsCacheError: nil,
			},
			helpers.ExpectedJSONOutput[models.PostViewArray]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data: &models.PostViewArray{
					Posts: []models.PostView{*defaultPost.PostView(&models.PostViewParams{
						UserID:       testUserID,
						LikeCount:    1,
						CommentCount: 2,
					})},
					NextPageURL: helpers.GeneratePostNextPageURL(models.BackendAddress, testPostID, map[string]interface{}{}),
				},
			},
		},
		{
			"Get posts OK - different cutoff",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				PostDBOutput:       []models.Post{diffCutoffPost},
				PostDBError:        nil,
				LikesCacheVal:      1,
				LikesCacheError:    nil,
				CommentsCacheVal:   2,
				CommentsCacheError: nil,
			},
			helpers.ExpectedJSONOutput[models.PostViewArray]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data: &models.PostViewArray{
					Posts: []models.PostView{*diffCutoffPost.PostView(&models.PostViewParams{
						UserID:       testUserID,
						LikeCount:    1,
						CommentCount: 2,
					})},
					NextPageURL: helpers.GeneratePostNextPageURL(models.BackendAddress, testDiffCutoffPostID, map[string]interface{}{}),
				},
			},
		},
		{
			"Get posts OK - community feed",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				QueryParams: map[string]interface{}{
					helpers.CommunityIDQueryKey: testCommunityID,
				},
				PostDBOutput:       []models.Post{defaultPost},
				PostDBError:        nil,
				LikesCacheVal:      1,
				LikesCacheError:    nil,
				CommentsCacheVal:   2,
				CommentsCacheError: nil,
			},
			helpers.ExpectedJSONOutput[models.PostViewArray]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data: &models.PostViewArray{
					Posts: []models.PostView{*defaultPost.PostView(&models.PostViewParams{
						UserID:       testUserID,
						LikeCount:    1,
						CommentCount: 2,
					})},
					NextPageURL: helpers.GeneratePostNextPageURL(models.BackendAddress, testPostID, map[string]interface{}{
						helpers.CommunityIDQueryKey: testCommunityID,
					}),
				},
			},
		},
		{
			"Get posts OK - project feed",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				QueryParams: map[string]interface{}{
					helpers.ProjectIDQueryKey: testProjectID,
				},
				PostDBOutput:       []models.Post{defaultPost},
				PostDBError:        nil,
				LikesCacheVal:      1,
				LikesCacheError:    nil,
				CommentsCacheVal:   2,
				CommentsCacheError: nil,
			},
			helpers.ExpectedJSONOutput[models.PostViewArray]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data: &models.PostViewArray{
					Posts: []models.PostView{*defaultPost.PostView(&models.PostViewParams{
						UserID:       testUserID,
						LikeCount:    1,
						CommentCount: 2,
					})},
					NextPageURL: helpers.GeneratePostNextPageURL(models.BackendAddress, testPostID, map[string]interface{}{
						helpers.ProjectIDQueryKey: testProjectID,
					}),
				},
			},
		},
		{
			"Get posts OK - multiple posts",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				PostDBOutput:       []models.Post{diffCutoffPost, defaultPost},
				PostDBError:        nil,
				LikesCacheVal:      1,
				LikesCacheError:    nil,
				CommentsCacheVal:   2,
				CommentsCacheError: nil,
			},
			helpers.ExpectedJSONOutput[models.PostViewArray]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data: &models.PostViewArray{
					Posts: []models.PostView{
						*diffCutoffPost.PostView(&models.PostViewParams{
							UserID:       testUserID,
							LikeCount:    1,
							CommentCount: 2,
						}),
						*defaultPost.PostView(&models.PostViewParams{
							UserID:       testUserID,
							LikeCount:    1,
							CommentCount: 2,
						}),
					},
					NextPageURL: helpers.GeneratePostNextPageURL(models.BackendAddress, testPostID, map[string]interface{}{}),
				},
			},
		},
		{
			"Get posts - invalid cutoff",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				QueryParams: map[string]interface{}{
					helpers.CutoffKey: invalidCutoff,
				},
			},
			helpers.ExpectedJSONOutput[models.PostViewArray]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrBadBinding,
			},
		},
		{
			"Get posts - invalid community ID",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				QueryParams: map[string]interface{}{
					helpers.CommunityIDQueryKey: invalidCommunityID,
				},
			},
			helpers.ExpectedJSONOutput[models.PostViewArray]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrBadBinding,
			},
		},
		{
			"Get posts - invalid project ID",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				QueryParams: map[string]interface{}{
					helpers.ProjectIDQueryKey: invalidProjectID,
				},
			},
			helpers.ExpectedJSONOutput[models.PostViewArray]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrBadBinding,
			},
		},
		{
			"Get posts - not found",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				PostDBOutput: []models.Post{defaultPost},
				PostDBError:  ErrTest,
			},
			helpers.ExpectedJSONOutput[models.PostViewArray]{
				StatusCode: http.StatusNotFound,
				JSONType:   helpers.ExpectedError,
				Error:      ErrPostNotFound,
			},
		},
		{
			"Get posts likes cache error OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				PostDBOutput:    []models.Post{defaultPost},
				PostDBError:     nil,
				LikesCacheVal:   1,
				LikesCacheError: ErrTest,
			},
			helpers.ExpectedJSONOutput[models.PostViewArray]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data: &models.PostViewArray{
					Posts:       []models.PostView{},
					NextPageURL: helpers.GeneratePostNextPageURL(models.BackendAddress, testPostID, map[string]interface{}{}),
				},
			},
		},
		{
			"Get posts comments cache error OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				PostDBOutput:       []models.Post{defaultPost},
				PostDBError:        nil,
				LikesCacheVal:      1,
				LikesCacheError:    nil,
				CommentsCacheVal:   2,
				CommentsCacheError: ErrTest,
			},
			helpers.ExpectedJSONOutput[models.PostViewArray]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data: &models.PostViewArray{
					Posts:       []models.PostView{},
					NextPageURL: helpers.GeneratePostNextPageURL(models.BackendAddress, testPostID, map[string]interface{}{}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbTestHandler := &PostDBTestHandler{}
			likesCacheTestHandler := &helpers.TestCache{}
			commentsCacheTestHandler := &helpers.TestCache{}
			a := &APIEnv{
				PostDBHandler:        dbTestHandler,
				LikesCacheHandler:    likesCacheTestHandler,
				CommentsCacheHandler: commentsCacheTestHandler,
			}

			c, w := helpers.CreateTestContextAndRecorder()

			for paramKey, paramVal := range tt.args.ContextParams {
				helpers.AddParamsToContext(c, paramKey, paramVal)
			}

			req, err := helpers.GenerateHttpJSONRequest(http.MethodGet, nil)
			if err != nil {
				t.Error(err)
			}

			for paramKey, paramVal := range tt.args.QueryParams {
				helpers.AddParamsToQuery(req, paramKey, paramVal)
			}

			c.Request = req

			dbTestHandler.SetMockGetPostsFunc(tt.args.PostDBOutput, tt.args.PostDBError)
			likesCacheTestHandler.SetMockGetCacheValFunc(tt.args.LikesCacheVal, tt.args.LikesCacheError)
			commentsCacheTestHandler.SetMockGetCacheValFunc(tt.args.CommentsCacheVal, tt.args.CommentsCacheError)
			a.GetPosts(c)

			b, _ := io.ReadAll(w.Body)
			if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(tt.expected.StatusCode, w.Code); !isEqual {
				t.Error(errStr)
			}

			m, err := helpers.ParseJSONString(b)
			if err != nil {
				t.Error(err)
			}

			if errStr, isEqual := helpers.CheckExpectedJSONEqualsActual(m, tt.expected); !isEqual {
				t.Error(errStr)
			}
		})
	}
}

func TestAPIEnv_GetPostByID(t *testing.T) {
	type args struct {
		ContextParams      map[string]interface{}
		PostDBOutput       *models.Post
		PostDBError        error
		LikesCacheVal      uint64
		LikesCacheError    error
		CommentsCacheVal   uint64
		CommentsCacheError error
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.PostView]
	}{
		{
			"Get Post By ID - OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: testPostID,
				},
				PostDBOutput:       &defaultPost,
				PostDBError:        nil,
				LikesCacheVal:      1,
				LikesCacheError:    nil,
				CommentsCacheVal:   2,
				CommentsCacheError: nil,
			},
			helpers.ExpectedJSONOutput[models.PostView]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data: defaultPost.PostView(&models.PostViewParams{
					UserID:       testUserID,
					LikeCount:    1,
					CommentCount: 2,
				}),
			},
		},
		{
			"Get Post By ID - Invalid Post ID",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: invalidPostID,
				},
			},
			helpers.ExpectedJSONOutput[models.PostView]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrPostNotFound,
			},
		},
		{
			"Get Post By ID - Cannot get post",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: testPostID,
				},
				PostDBOutput: nil,
				PostDBError:  ErrTest,
			},
			helpers.ExpectedJSONOutput[models.PostView]{
				StatusCode: http.StatusNotFound,
				JSONType:   helpers.ExpectedError,
				Error:      ErrPostNotFound,
			},
		},
		{
			"Get Post By ID - Cannot get likes count",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: testPostID,
				},
				PostDBOutput:    &defaultPost,
				PostDBError:     nil,
				LikesCacheVal:   0,
				LikesCacheError: ErrTest,
			},
			helpers.ExpectedJSONOutput[models.PostView]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrTest,
			},
		},
		{
			"Get Post By ID - Cannot get comments count",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: testPostID,
				},
				PostDBOutput:       &defaultPost,
				PostDBError:        nil,
				LikesCacheVal:      1,
				LikesCacheError:    nil,
				CommentsCacheVal:   0,
				CommentsCacheError: ErrTest,
			},
			helpers.ExpectedJSONOutput[models.PostView]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrTest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbTestHandler := &PostDBTestHandler{}
			likesCacheTestHandler := &helpers.TestCache{}
			commentsCacheTestHandler := &helpers.TestCache{}
			a := &APIEnv{
				PostDBHandler:        dbTestHandler,
				LikesCacheHandler:    likesCacheTestHandler,
				CommentsCacheHandler: commentsCacheTestHandler,
			}

			c, w := helpers.CreateTestContextAndRecorder()

			for paramKey, paramVal := range tt.args.ContextParams {
				helpers.AddParamsToContext(c, paramKey, paramVal)
			}

			req, err := helpers.GenerateHttpJSONRequest(http.MethodGet, nil)
			if err != nil {
				t.Error(err)
			}

			c.Request = req

			dbTestHandler.SetMockGetPostByIDFunc(tt.args.PostDBOutput, tt.args.PostDBError)
			likesCacheTestHandler.SetMockGetCacheValFunc(tt.args.LikesCacheVal, tt.args.LikesCacheError)
			commentsCacheTestHandler.SetMockGetCacheValFunc(tt.args.CommentsCacheVal, tt.args.CommentsCacheError)
			a.GetPostByID(c)

			b, _ := io.ReadAll(w.Body)
			if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(tt.expected.StatusCode, w.Code); !isEqual {
				t.Error(errStr)
			}

			m, err := helpers.ParseJSONString(b)
			if err != nil {
				t.Error(err)
			}

			if errStr, isEqual := helpers.CheckExpectedJSONEqualsActual(m, tt.expected); !isEqual {
				t.Error(errStr)
			}
		})
	}
}

func TestAPIEnv_UpdatePost(t *testing.T) {
	type args struct {
		ContextParams      map[string]interface{}
		PostData           *models.Post
		PostDBOutput       *models.Post
		PostDBError        error
		LikesCacheVal      uint64
		LikesCacheError    error
		CommentsCacheVal   uint64
		CommentsCacheError error
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.PostView]
	}{
		{
			"Update Post OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: testPostID,
				},
				PostData:           &newTestPost,
				PostDBOutput:       &defaultPost,
				PostDBError:        nil,
				LikesCacheVal:      1,
				LikesCacheError:    nil,
				CommentsCacheVal:   2,
				CommentsCacheError: nil,
			},
			helpers.ExpectedJSONOutput[models.PostView]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data: defaultPost.PostView(&models.PostViewParams{
					UserID:       testUserID,
					LikeCount:    1,
					CommentCount: 2,
				}),
			},
		},
		{
			"Update Post - Invalid Post ID",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: invalidPostID,
				},
			},
			helpers.ExpectedJSONOutput[models.PostView]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrPostNotFound,
			},
		},
		{
			"Update Post - Bad Binding",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: testPostID,
				},
			},
			helpers.ExpectedJSONOutput[models.PostView]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrBadBinding,
			},
		},
		{
			"Update Post - Post Not Found",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: testPostID,
				},
				PostData:     &newTestPost,
				PostDBOutput: nil,
				PostDBError:  gorm.ErrRecordNotFound,
			},
			helpers.ExpectedJSONOutput[models.PostView]{
				StatusCode: http.StatusNotFound,
				JSONType:   helpers.ExpectedError,
				Error:      ErrPostNotFound,
			},
		},
		{
			"Update Post - Not owner",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: diffUserID,
					helpers.PostIDKey: testPostID,
				},
				PostData:     &newTestPost,
				PostDBOutput: nil,
				PostDBError:  helpers.ErrNotOwner,
			},
			helpers.ExpectedJSONOutput[models.PostView]{
				StatusCode: http.StatusForbidden,
				JSONType:   helpers.ExpectedError,
				Error:      helpers.ErrNotOwner,
			},
		},
		{
			"Update Post - Cannot update",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: diffUserID,
					helpers.PostIDKey: testPostID,
				},
				PostData:     &newTestPost,
				PostDBOutput: nil,
				PostDBError:  ErrTest,
			},
			helpers.ExpectedJSONOutput[models.PostView]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrCannotUpdatePost,
			},
		},
		{
			"Update Post - Cannot get likes count",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: diffUserID,
					helpers.PostIDKey: testPostID,
				},
				PostData:        &newTestPost,
				PostDBOutput:    &defaultPost,
				PostDBError:     nil,
				LikesCacheVal:   0,
				LikesCacheError: ErrTest,
			},
			helpers.ExpectedJSONOutput[models.PostView]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrTest,
			},
		},
		{
			"Update Post - Cannot get comments count",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: diffUserID,
					helpers.PostIDKey: testPostID,
				},
				PostData:           &newTestPost,
				PostDBOutput:       &defaultPost,
				PostDBError:        nil,
				LikesCacheVal:      1,
				LikesCacheError:    nil,
				CommentsCacheVal:   0,
				CommentsCacheError: ErrTest,
			},
			helpers.ExpectedJSONOutput[models.PostView]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrTest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbTestHandler := &PostDBTestHandler{}
			likesCacheTestHandler := &helpers.TestCache{}
			commentsCacheTestHandler := &helpers.TestCache{}
			a := &APIEnv{
				PostDBHandler:        dbTestHandler,
				LikesCacheHandler:    likesCacheTestHandler,
				CommentsCacheHandler: commentsCacheTestHandler,
			}

			c, w := helpers.CreateTestContextAndRecorder()

			for paramKey, paramVal := range tt.args.ContextParams {
				helpers.AddParamsToContext(c, paramKey, paramVal)
			}

			if tt.args.PostData != nil {
				req, err := helpers.GenerateHttpJSONRequest(http.MethodPatch, tt.args.PostData)
				if err != nil {
					t.Error(err)
				}

				c.Request = req
			}

			dbTestHandler.SetMockUpdatePostFunc(tt.args.PostDBOutput, tt.args.PostDBError)
			likesCacheTestHandler.SetMockGetCacheValFunc(tt.args.LikesCacheVal, tt.args.LikesCacheError)
			commentsCacheTestHandler.SetMockGetCacheValFunc(tt.args.CommentsCacheVal, tt.args.CommentsCacheError)
			a.UpdatePost(c)

			b, _ := io.ReadAll(w.Body)
			if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(tt.expected.StatusCode, w.Code); !isEqual {
				t.Error(errStr)
			}

			m, err := helpers.ParseJSONString(b)
			if err != nil {
				t.Error(err)
			}

			if errStr, isEqual := helpers.CheckExpectedJSONEqualsActual(m, tt.expected); !isEqual {
				t.Error(errStr)
			}
		})
	}
}
