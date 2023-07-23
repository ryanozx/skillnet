package controllers

import (
	"io"
	"net/http"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
)

const (
	likeCountLiked   = 1
	likeCountUnliked = 0
)

var (
	defaultLike = models.Like{
		UserID: testUserID,
		PostID: testPostID,
	}
	defaultCreateLikeUpdate = models.LikeUpdate{
		Like: models.Like{
			PostID: testPostID,
		},
		LikeCount: likeCountLiked,
	}
	defaultDeleteLikeUpdate = models.LikeUpdate{
		LikeCount: likeCountUnliked,
	}
)

// Mock DB Handler
type LikeDBTestHandler struct {
	CreateLikeFunc  func(*models.Like) (*models.Like, error)
	DeleteLikeFunc  func(string, uint) error
	GetCountFunc    func(uint) (uint64, error)
	GetLikeByIDFunc func(string) (*models.Like, error)
}

func (h *LikeDBTestHandler) CreateLike(newLike *models.Like) (*models.Like, error) {
	return h.CreateLikeFunc(newLike)
}

func (h *LikeDBTestHandler) DeleteLike(userID string, postID uint) error {
	return h.DeleteLikeFunc(userID, postID)
}

func (h *LikeDBTestHandler) GetValue(postID uint) (uint64, error) {
	return h.GetCountFunc(postID)
}

func (h *LikeDBTestHandler) GetLikeByID(likeID string) (*models.Like, error) {
	return h.GetLikeByIDFunc(likeID)
}

func (h *LikeDBTestHandler) SetMockCreateLikeFunc(like *models.Like, err error) {
	h.CreateLikeFunc = func(newLike *models.Like) (*models.Like, error) {
		return like, err
	}
}

func (h *LikeDBTestHandler) SetMockDeleteLikeFunc(err error) {
	h.DeleteLikeFunc = func(userID string, postID uint) error {
		return err
	}
}

func (h *LikeDBTestHandler) SetMockGetLikeCountFunc(count uint64, err error) {
	h.GetCountFunc = func(postID uint) (uint64, error) {
		return count, err
	}
}

func (h *LikeDBTestHandler) SetMockGetLikeByIDFunc(like *models.Like, err error) {
	h.GetLikeByIDFunc = func(likeID string) (*models.Like, error) {
		return like, err
	}
}

func TestAPIEnv_InitialiseLikeHandler(t *testing.T) {
	type fields struct {
		DB     *gorm.DB
		client *redis.Client
	}
	tests := []struct {
		name               string
		fields             fields
		expectedDBEmpty    bool
		expectedCacheEmpty bool
	}{
		{
			"Initialise Likes DB and Cache OK",
			fields{
				DB:     &gorm.DB{},
				client: &redis.Client{},
			},
			false,
			false,
		},
		{
			"Initialise Likes DB OK, cache empty",
			fields{
				DB: &gorm.DB{},
			},
			false,
			true,
		},
		{
			"Initialise Likes DB empty, cache OK",
			fields{
				client: &redis.Client{},
			},
			true,
			false,
		},
		{
			"Initialise Likes DB empty, cache empty",
			fields{},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIEnv{
				DB: tt.fields.DB,
			}
			a.InitialiseLikeHandler(tt.fields.client)
			likesDB, ok := a.LikeDBHandler.(*database.LikeDB)
			if ok {
				if tt.expectedDBEmpty && likesDB.DB != nil {
					t.Error("Likes DB contains unexpected DB instance")
				} else if !tt.expectedDBEmpty && likesDB.DB != tt.fields.DB {
					t.Error("LikesDBHandler not initialised correctly")
				}
			} else {
				t.Error("LikesDBHandler is nil!")
			}

			if likesCache, ok := a.LikesCacheHandler.(*Cache); ok {
				if tt.expectedCacheEmpty && likesCache.redisDB != nil {
					t.Error("Likes cache contains unexpected cache instance")
				} else if !tt.expectedCacheEmpty && (likesCache.redisDB != tt.fields.client || likesCache.DBHandler != a.LikeDBHandler) {
					t.Error("Likes cache not initialised correctly")
				}
			} else {
				t.Error("Likes cache is nil!")
			}
		})
	}
}

func TestAPIEnv_PostLike(t *testing.T) {
	type args struct {
		ContextParams     map[string]interface{}
		LikeDBOutput      *models.Like
		LikeDBError       error
		LikeCacheOutput   uint64
		LikeCacheError    error
		NotificationError error
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.LikeUpdate]
	}{
		{
			"Create Like OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: testPostID,
				},
				LikeDBOutput:      &defaultLike,
				LikeDBError:       nil,
				LikeCacheOutput:   likeCountLiked,
				LikeCacheError:    nil,
				NotificationError: nil,
			},
			helpers.ExpectedJSONOutput[models.LikeUpdate]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data:       &defaultCreateLikeUpdate,
			},
		},
		{
			"Create Like invalid Post ID",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: invalidPostID,
				},
			},
			helpers.ExpectedJSONOutput[models.LikeUpdate]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrPostNotFound,
			},
		},
		{
			"Create Like negative Post ID",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: negativePostID,
				},
			},
			helpers.ExpectedJSONOutput[models.LikeUpdate]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrPostNotFound,
			},
		},
		{
			"Create Like already liked",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: testPostID,
				},
				LikeDBOutput: &defaultLike,
				LikeDBError:  gorm.ErrDuplicatedKey,
			},
			helpers.ExpectedJSONOutput[models.LikeUpdate]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrAlreadyLiked,
			},
		},
		{
			"Create Like cannot create like",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: testPostID,
				},
				LikeDBOutput: &defaultLike,
				LikeDBError:  ErrTest,
			},
			helpers.ExpectedJSONOutput[models.LikeUpdate]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrLikeNotRegistered,
			},
		},
		{
			"Create Like cannot set cache count",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: testPostID,
				},
				LikeDBOutput:    &defaultLike,
				LikeDBError:     nil,
				LikeCacheOutput: likeCountLiked,
				LikeCacheError:  ErrTest,
			},
			helpers.ExpectedJSONOutput[models.LikeUpdate]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrUpdateLikeCountFailed,
			},
		},
		{
			// Errors from the notification API should not throw an error back to the client
			"Create Like notification error OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: testPostID,
				},
				LikeDBOutput:      &defaultLike,
				LikeDBError:       nil,
				LikeCacheOutput:   likeCountLiked,
				LikeCacheError:    nil,
				NotificationError: ErrTest,
			},
			helpers.ExpectedJSONOutput[models.LikeUpdate]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data:       &defaultCreateLikeUpdate,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbTestHandler := &LikeDBTestHandler{}
			cacheTestHandler := &helpers.TestCache{}
			notifPoster := &helpers.TestNotificationCreator{}
			a := &APIEnv{
				LikeDBHandler:      dbTestHandler,
				LikesCacheHandler:  cacheTestHandler,
				NotificationPoster: notifPoster,
			}

			c, w := helpers.CreateTestContextAndRecorder()

			for paramKey, paramVal := range tt.args.ContextParams {
				helpers.AddParamsToContext(c, paramKey, paramVal)
			}

			req, err := helpers.GenerateHttpJSONRequest(http.MethodPost, nil)
			if err != nil {
				t.Error(err)
			}
			c.Request = req

			dbTestHandler.SetMockCreateLikeFunc(tt.args.LikeDBOutput, tt.args.LikeDBError)
			cacheTestHandler.SetMockSetCacheValFunc(tt.args.LikeCacheOutput, tt.args.LikeCacheError)
			notifPoster.SetMockPostNotificationFromEventFunc(tt.args.NotificationError)
			a.PostLike(c)

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

func TestAPIEnv_DeleteLike(t *testing.T) {
	type args struct {
		ContextParams   map[string]interface{}
		LikeDBError     error
		LikeCacheOutput uint64
		LikeCacheError  error
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.LikeUpdate]
	}{
		{
			"Delete Like OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: testPostID,
				},
				LikeDBError:     nil,
				LikeCacheOutput: likeCountUnliked,
				LikeCacheError:  nil,
			},
			helpers.ExpectedJSONOutput[models.LikeUpdate]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data:       &defaultDeleteLikeUpdate,
			},
		},
		{
			"Delete Like invalid PostID",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: invalidPostID,
				},
			},
			helpers.ExpectedJSONOutput[models.LikeUpdate]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrPostNotFound,
			},
		},
		{
			"Create Like negative Post ID",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: negativePostID,
				},
			},
			helpers.ExpectedJSONOutput[models.LikeUpdate]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrPostNotFound,
			},
		},
		{
			"Delete Like not found",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: testPostID,
				},
				LikeDBError: gorm.ErrRecordNotFound,
			},
			helpers.ExpectedJSONOutput[models.LikeUpdate]{
				StatusCode: http.StatusNotFound,
				JSONType:   helpers.ExpectedError,
				Error:      ErrLikeNotFound,
			},
		},
		{
			"Delete Like cannot delete like",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: testPostID,
				},
				LikeDBError: ErrTest,
			},
			helpers.ExpectedJSONOutput[models.LikeUpdate]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrUnlikeFailed,
			},
		},
		{
			"Delete Like cannot set cache count",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
					helpers.PostIDKey: testPostID,
				},
				LikeDBError:     nil,
				LikeCacheOutput: likeCountUnliked,
				LikeCacheError:  ErrTest,
			},
			helpers.ExpectedJSONOutput[models.LikeUpdate]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrUpdateLikeCountFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbTestHandler := &LikeDBTestHandler{}
			cacheTestHandler := &helpers.TestCache{}
			a := &APIEnv{
				LikeDBHandler:     dbTestHandler,
				LikesCacheHandler: cacheTestHandler,
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

			dbTestHandler.SetMockDeleteLikeFunc(tt.args.LikeDBError)
			cacheTestHandler.SetMockSetCacheValFunc(tt.args.LikeCacheOutput, tt.args.LikeCacheError)
			a.DeleteLike(c)

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
