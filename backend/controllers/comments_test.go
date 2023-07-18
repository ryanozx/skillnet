package controllers

import (
	"io"
	"net/http"
	"testing"

	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
)

const (
	testCommentID     = 1
	diffCommentID     = 10
	invalidCommentID  = "badCommentID"
	negativeCommentID = -1
)

var (
	newTestComment = models.Comment{
		Text: "Test comment",
	}
	defaultComment = models.Comment{
		Model: gorm.Model{
			ID: testCommentID,
		},
		UserID: testUserID,
		PostID: testPostID,
		Text:   "TextComment",
		User:   defaultUser,
		Post:   defaultPost,
	}
	diffCutoffComment = models.Comment{
		Model: gorm.Model{
			ID: diffCommentID,
		},
		UserID: testUserID,
		PostID: testPostID,
		Text:   "Diff ID",
		User:   defaultUser,
		Post:   defaultPost,
	}
)

type CommentsDBTestHandler struct {
	CreateCommentFunc  func(*models.Comment) (*models.Comment, error)
	DeleteCommentFunc  func(uint, string) (uint, error)
	GetCommentsFunc    func(uint, *helpers.NullableUint) ([]models.Comment, error)
	GetCommentByIDFunc func(uint) (*models.Comment, error)
	UpdateCommentFunc  func(*models.Comment, uint, string) (*models.Comment, error)
	GetValueFunc       func(uint) (uint64, error)
}

func (h *CommentsDBTestHandler) CreateComment(comment *models.Comment) (*models.Comment, error) {
	return h.CreateCommentFunc(comment)
}

func (h *CommentsDBTestHandler) DeleteComment(commentID uint, userID string) (uint, error) {
	return h.DeleteCommentFunc(commentID, userID)
}

func (h *CommentsDBTestHandler) GetComments(postID uint, cutoff *helpers.NullableUint) ([]models.Comment, error) {
	return h.GetCommentsFunc(postID, cutoff)
}

func (h *CommentsDBTestHandler) GetCommentByID(commentID uint) (*models.Comment, error) {
	return h.GetCommentByIDFunc(commentID)
}

func (h *CommentsDBTestHandler) UpdateComment(comment *models.Comment, commentID uint, userID string) (*models.Comment, error) {
	return h.UpdateCommentFunc(comment, commentID, userID)
}

func (h *CommentsDBTestHandler) GetValue(postID uint) (uint64, error) {
	return h.GetValueFunc(postID)
}

func (h *CommentsDBTestHandler) SetMockCreateCommentFunc(newComment *models.Comment, err error) {
	h.CreateCommentFunc = func(comment *models.Comment) (*models.Comment, error) {
		return newComment, err
	}
}

func (h *CommentsDBTestHandler) SetMockDeleteCommentFunc(postID uint, err error) {
	h.DeleteCommentFunc = func(commentID uint, userID string) (uint, error) {
		return postID, err
	}
}

func (h *CommentsDBTestHandler) SetMockGetCommentsFunc(comments []models.Comment, err error) {
	h.GetCommentsFunc = func(postID uint, cutoff *helpers.NullableUint) ([]models.Comment, error) {
		return comments, err
	}
}

func (h *CommentsDBTestHandler) SetMockGetCommentByIDFunc(comment *models.Comment, err error) {
	h.GetCommentByIDFunc = func(commentID uint) (*models.Comment, error) {
		return comment, err
	}
}

func (h *CommentsDBTestHandler) SetMockUpdateCommentFunc(updatedComment *models.Comment, err error) {
	h.UpdateCommentFunc = func(comment *models.Comment, commentID uint, userID string) (*models.Comment, error) {
		return updatedComment, err
	}
}

func (h *CommentsDBTestHandler) SetMockGetValueFunc(count uint64, err error) {
	h.GetValueFunc = func(postID uint) (uint64, error) {
		return count, err
	}
}

func TestAPIEnv_CreateComment(t *testing.T) {
	helpers.SetEnvVars(t)

	type args struct {
		ContextParams      map[string]interface{}
		QueryParams        map[string]interface{}
		CommentData        *models.Comment
		CommentDBOutput    *models.Comment
		CommentDBError     error
		CommentCacheOutput uint64
		CommentCacheError  error
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.CommentUpdate]
	}{
		{
			"Create comment OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				QueryParams: map[string]interface{}{
					helpers.PostIDQueryKey: testPostID,
				},
				CommentData:        &newTestComment,
				CommentDBOutput:    &defaultComment,
				CommentDBError:     nil,
				CommentCacheOutput: 1,
				CommentCacheError:  nil,
			},
			helpers.ExpectedJSONOutput[models.CommentUpdate]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data: &models.CommentUpdate{
					Comment:      *defaultComment.CommentView(testUserID),
					CommentCount: 1,
				},
			},
		},
		{
			"Create comment bad binding",
			args{},
			helpers.ExpectedJSONOutput[models.CommentUpdate]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrBadBinding,
			},
		},
		{
			"Create comment no Post ID",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				CommentData: &newTestComment,
			},
			helpers.ExpectedJSONOutput[models.CommentUpdate]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrPostNotFound,
			},
		},
		{
			"Create comment empty Post ID string",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				QueryParams: map[string]interface{}{
					helpers.PostIDQueryKey: "",
				},
				CommentData: &newTestComment,
			},
			helpers.ExpectedJSONOutput[models.CommentUpdate]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrPostNotFound,
			},
		},
		{
			"Create comment invalid Post ID",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				QueryParams: map[string]interface{}{
					helpers.PostIDQueryKey: -1,
				},
				CommentData: &newTestComment,
			},
			helpers.ExpectedJSONOutput[models.CommentUpdate]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrPostNotFound,
			},
		},
		{
			"Create comment DB throws error",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				QueryParams: map[string]interface{}{
					helpers.PostIDQueryKey: testPostID,
				},
				CommentData:     &newTestComment,
				CommentDBOutput: &defaultComment,
				CommentDBError:  ErrTest,
			},
			helpers.ExpectedJSONOutput[models.CommentUpdate]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrCannotCreateComment,
			},
		},
		{
			"Create comment fail to update cache count",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				QueryParams: map[string]interface{}{
					helpers.PostIDQueryKey: testPostID,
				},
				CommentData:        &newTestComment,
				CommentDBOutput:    &defaultComment,
				CommentDBError:     nil,
				CommentCacheOutput: 0,
				CommentCacheError:  ErrTest,
			},
			helpers.ExpectedJSONOutput[models.CommentUpdate]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrTest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbTestHandler := &CommentsDBTestHandler{}
			cacheTestHandler := &helpers.TestCache{}
			notifPoster := &helpers.TestNotificationCreator{}
			a := &APIEnv{
				CommentDBHandler:     dbTestHandler,
				CommentsCacheHandler: cacheTestHandler,
				NotificationPoster:   notifPoster,
			}

			c, w := helpers.CreateTestContextAndRecorder()
			helpers.AddStoreToContext(c, helpers.MakeMockStore())

			for paramKey, paramVal := range tt.args.ContextParams {
				helpers.AddParamsToContext(c, paramKey, paramVal)
			}

			if tt.args.CommentData != nil {
				req, err := helpers.GenerateHttpJSONRequest(http.MethodPost, tt.args.CommentData)
				if err != nil {
					t.Error(err)
				}
				c.Request = req

				for paramKey, paramVal := range tt.args.QueryParams {
					helpers.AddParamsToQuery(req, paramKey, paramVal)
				}
			}
			dbTestHandler.SetMockCreateCommentFunc(tt.args.CommentDBOutput, tt.args.CommentDBError)
			cacheTestHandler.SetMockSetCacheValFunc(tt.args.CommentCacheOutput, tt.args.CommentCacheError)
			notifPoster.SetMockPostNotificationFromEventFunc(nil)
			a.CreateComment(c)

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

func TestAPIEnv_DeleteComment(t *testing.T) {
	helpers.SetEnvVars(t)

	type args struct {
		ContextParams      map[string]interface{}
		CommentDBOutput    uint
		CommentDBError     error
		CommentCacheOutput uint64
		CommentCacheError  error
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.CommentUpdate]
	}{
		{
			"Delete comment OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:    testUserID,
					helpers.CommentIDKey: testCommentID,
				},
				CommentDBOutput:    1,
				CommentDBError:     nil,
				CommentCacheOutput: 0,
				CommentCacheError:  nil,
			},
			helpers.ExpectedJSONOutput[models.CommentUpdate]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data: &models.CommentUpdate{
					CommentCount: 0,
				},
			},
		},
		{
			"Delete comment invalid Post ID",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:    testUserID,
					helpers.CommentIDKey: -1,
				},
			},
			helpers.ExpectedJSONOutput[models.CommentUpdate]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrCommentNotFound,
			},
		},
		{
			"Delete comment DB not found",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:    testUserID,
					helpers.CommentIDKey: testCommentID,
				},
				CommentDBOutput: 0,
				CommentDBError:  gorm.ErrRecordNotFound,
			},
			helpers.ExpectedJSONOutput[models.CommentUpdate]{
				StatusCode: http.StatusNotFound,
				JSONType:   helpers.ExpectedError,
				Error:      ErrCommentNotFound,
			},
		},
		{
			"Delete comment DB not owner",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:    testUserID,
					helpers.CommentIDKey: testCommentID,
				},
				CommentDBOutput: 0,
				CommentDBError:  helpers.ErrNotOwner,
			},
			helpers.ExpectedJSONOutput[models.CommentUpdate]{
				StatusCode: http.StatusForbidden,
				JSONType:   helpers.ExpectedError,
				Error:      helpers.ErrNotOwner,
			},
		},
		{
			"Delete comment DB throws error",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:    testUserID,
					helpers.CommentIDKey: testCommentID,
				},
				CommentDBOutput: 0,
				CommentDBError:  ErrTest,
			},
			helpers.ExpectedJSONOutput[models.CommentUpdate]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrCannotDeleteComment,
			},
		},
		{
			"Delete comment fail to update cache count",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:    testUserID,
					helpers.CommentIDKey: testCommentID,
				},
				CommentDBOutput:    1,
				CommentDBError:     nil,
				CommentCacheOutput: 0,
				CommentCacheError:  ErrTest,
			},
			helpers.ExpectedJSONOutput[models.CommentUpdate]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrTest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbTestHandler := &CommentsDBTestHandler{}
			cacheTestHandler := &helpers.TestCache{}
			a := &APIEnv{
				CommentDBHandler:     dbTestHandler,
				CommentsCacheHandler: cacheTestHandler,
			}

			c, w := helpers.CreateTestContextAndRecorder()
			helpers.AddStoreToContext(c, helpers.MakeMockStore())

			for paramKey, paramVal := range tt.args.ContextParams {
				helpers.AddParamsToContext(c, paramKey, paramVal)
			}

			dbTestHandler.SetMockDeleteCommentFunc(tt.args.CommentDBOutput, tt.args.CommentDBError)
			cacheTestHandler.SetMockSetCacheValFunc(tt.args.CommentCacheOutput, tt.args.CommentCacheError)
			a.DeleteComment(c)

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

func TestAPIEnv_GetComments(t *testing.T) {
	helpers.SetEnvVars(t)

	type args struct {
		ContextParams   map[string]interface{}
		QueryParams     map[string]interface{}
		CommentDBOutput []models.Comment
		CommentDBError  error
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.CommentViewsArray]
	}{
		{
			"Get comments no cutoff OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				QueryParams: map[string]interface{}{
					helpers.PostIDQueryKey: testPostID,
				},
				CommentDBOutput: []models.Comment{defaultComment},
				CommentDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.CommentViewsArray]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data: &models.CommentViewsArray{
					Comments:    []models.CommentView{*defaultComment.CommentView(testUserID)},
					NextPageURL: helpers.GenerateCommentNextPageURL(models.BackendAddress, testPostID, testCommentID),
				},
			},
		},
		{
			"Get comments no comments",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				QueryParams: map[string]interface{}{
					helpers.PostIDQueryKey: testPostID,
				},
				CommentDBOutput: []models.Comment{},
				CommentDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.CommentViewsArray]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data: &models.CommentViewsArray{
					Comments:    []models.CommentView{},
					NextPageURL: helpers.GenerateCommentNextPageURL(models.BackendAddress, testPostID, 0),
				},
			},
		},
		{
			"Get comments different comment cutoff OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				QueryParams: map[string]interface{}{
					helpers.PostIDQueryKey: testPostID,
				},
				CommentDBOutput: []models.Comment{diffCutoffComment},
				CommentDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.CommentViewsArray]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data: &models.CommentViewsArray{
					Comments:    []models.CommentView{*diffCutoffComment.CommentView(testUserID)},
					NextPageURL: helpers.GenerateCommentNextPageURL(models.BackendAddress, testPostID, 10),
				},
			},
		},
		{
			"Get comments not owner OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: diffUserID,
				},
				QueryParams: map[string]interface{}{
					helpers.PostIDQueryKey: testPostID,
				},
				CommentDBOutput: []models.Comment{defaultComment},
				CommentDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.CommentViewsArray]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data: &models.CommentViewsArray{
					Comments:    []models.CommentView{*defaultComment.CommentView(diffUserID)},
					NextPageURL: helpers.GenerateCommentNextPageURL(models.BackendAddress, testPostID, testCommentID),
				},
			},
		},
		{
			"Get comments different query cutoff OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				QueryParams: map[string]interface{}{
					helpers.PostIDQueryKey: testPostID,
					helpers.CutoffKey:      validCutoff,
				},
				CommentDBOutput: []models.Comment{defaultComment},
				CommentDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.CommentViewsArray]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data: &models.CommentViewsArray{
					Comments:    []models.CommentView{*defaultComment.CommentView(testUserID)},
					NextPageURL: helpers.GenerateCommentNextPageURL(models.BackendAddress, testPostID, testCommentID),
				},
			},
		},
		{
			"Get comments invalid cutoff",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				QueryParams: map[string]interface{}{
					helpers.PostIDQueryKey: testPostID,
					helpers.CutoffKey:      invalidCutoff,
				},
				CommentDBOutput: []models.Comment{defaultComment},
				CommentDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.CommentViewsArray]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrBadBinding,
			},
		},
		{
			"Get comments invalid PostID",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				QueryParams: map[string]interface{}{
					helpers.PostIDQueryKey: invalidPostID,
				},
				CommentDBOutput: []models.Comment{defaultComment},
				CommentDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.CommentViewsArray]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrPostNotFound,
			},
		},
		{
			"Get comments negative PostID",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				QueryParams: map[string]interface{}{
					helpers.PostIDQueryKey: negativePostID,
				},
				CommentDBOutput: []models.Comment{defaultComment},
				CommentDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.CommentViewsArray]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrPostNotFound,
			},
		},
		{
			"Get comments DB throws error",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				QueryParams: map[string]interface{}{
					helpers.PostIDQueryKey: testPostID,
				},
				CommentDBOutput: []models.Comment{defaultComment},
				CommentDBError:  ErrTest,
			},
			helpers.ExpectedJSONOutput[models.CommentViewsArray]{
				StatusCode: http.StatusNotFound,
				JSONType:   helpers.ExpectedError,
				Error:      ErrCommentNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbTestHandler := &CommentsDBTestHandler{}
			a := &APIEnv{
				CommentDBHandler: dbTestHandler,
			}

			c, w := helpers.CreateTestContextAndRecorder()
			helpers.AddStoreToContext(c, helpers.MakeMockStore())

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

			dbTestHandler.SetMockGetCommentsFunc(tt.args.CommentDBOutput, tt.args.CommentDBError)
			a.GetComments(c)

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

func TestAPIEnv_UpdateComment(t *testing.T) {
	helpers.SetEnvVars(t)

	type args struct {
		ContextParams   map[string]interface{}
		CommentUpdate   *models.Comment
		CommentDBOutput *models.Comment
		CommentDBError  error
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.CommentView]
	}{
		{
			"Update comment OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:    testUserID,
					helpers.CommentIDKey: testCommentID,
				},
				CommentUpdate:   &newTestComment,
				CommentDBOutput: &defaultComment,
				CommentDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.CommentView]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data:       defaultComment.CommentView(testUserID),
			},
		},
		{
			"Update comment invalid comment ID",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:    testUserID,
					helpers.CommentIDKey: invalidCommentID,
				},
				CommentUpdate:   &newTestComment,
				CommentDBOutput: &defaultComment,
				CommentDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.CommentView]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrCommentNotFound,
			},
		},
		{
			"Update comment negative comment ID",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:    testUserID,
					helpers.CommentIDKey: negativeCommentID,
				},
				CommentUpdate:   &newTestComment,
				CommentDBOutput: &defaultComment,
				CommentDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.CommentView]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrCommentNotFound,
			},
		},
		{
			"Update comment bad binding",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:    testUserID,
					helpers.CommentIDKey: testCommentID,
				},
				CommentDBOutput: &defaultComment,
				CommentDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.CommentView]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrBadBinding,
			},
		},
		{
			"Update comment not found",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:    testUserID,
					helpers.CommentIDKey: testCommentID,
				},
				CommentUpdate:   &newTestComment,
				CommentDBOutput: &defaultComment,
				CommentDBError:  gorm.ErrRecordNotFound,
			},
			helpers.ExpectedJSONOutput[models.CommentView]{
				StatusCode: http.StatusNotFound,
				JSONType:   helpers.ExpectedError,
				Error:      ErrCommentNotFound,
			},
		},
		{
			"Update comment not owner",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:    diffUserID,
					helpers.CommentIDKey: testCommentID,
				},
				CommentUpdate:   &newTestComment,
				CommentDBOutput: &defaultComment,
				CommentDBError:  helpers.ErrNotOwner,
			},
			helpers.ExpectedJSONOutput[models.CommentView]{
				StatusCode: http.StatusForbidden,
				JSONType:   helpers.ExpectedError,
				Error:      helpers.ErrNotOwner,
			},
		},
		{
			"Update comment DB throws error",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:    diffUserID,
					helpers.CommentIDKey: testCommentID,
				},
				CommentUpdate:   &newTestComment,
				CommentDBOutput: &defaultComment,
				CommentDBError:  ErrTest,
			},
			helpers.ExpectedJSONOutput[models.CommentView]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrCannotUpdateComment,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbTestHandler := &CommentsDBTestHandler{}
			cacheTestHandler := &helpers.TestCache{}
			a := &APIEnv{
				CommentDBHandler:     dbTestHandler,
				CommentsCacheHandler: cacheTestHandler,
			}

			c, w := helpers.CreateTestContextAndRecorder()
			helpers.AddStoreToContext(c, helpers.MakeMockStore())

			for paramKey, paramVal := range tt.args.ContextParams {
				helpers.AddParamsToContext(c, paramKey, paramVal)
			}

			if tt.args.CommentUpdate != nil {
				req, err := helpers.GenerateHttpJSONRequest(http.MethodPost, tt.args.CommentUpdate)
				if err != nil {
					t.Error(err)
				}
				c.Request = req
			}
			dbTestHandler.SetMockUpdateCommentFunc(tt.args.CommentDBOutput, tt.args.CommentDBError)
			a.UpdateComment(c)

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
