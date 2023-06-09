package controllers

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

const (
	invalidPostID  = "badpostid"
	negativePostID = "-234"
)

var (
	defaultLike = models.Like{
		UserID: testUserID,
		PostID: testPostID,
	}
	defaultCreateLikeUpdate = models.LikeUpdate{
		Like:      defaultLike,
		LikeCount: 1,
	}
	defaultDeleteLikeUpdate = models.LikeUpdate{
		LikeCount: 0,
	}
)

type LikeControllerTestSuite struct {
	suite.Suite
	api          APIEnv
	dbHandler    *LikeDBTestHandler
	cacheHandler *LikeTestCache
	store        *helpers.MockSessionStore
}

func (s *LikeControllerTestSuite) SetupSuite() {
	dbHandler := LikeDBTestHandler{}
	cacheHandler := LikeTestCache{}
	api := APIEnv{
		LikeDBHandler:     &dbHandler,
		LikesCacheHandler: &cacheHandler,
	}
	s.api = api
	s.dbHandler = &dbHandler
	s.cacheHandler = &cacheHandler
	s.store = helpers.MakeMockStore()
}

func (s *LikeControllerTestSuite) TearDownTest() {
	s.dbHandler.ResetFuncs()
	s.store.Reset()
}

func TestLikeControllerSuite(t *testing.T) {
	t.Setenv("CLIENT_HOST", "http://localhost")
	t.Setenv("CLIENT_PORT", "3000")
	t.Setenv("BACKEND_HOST", "http://localhost")
	t.Setenv("BACKEND_PORT", "8080")
	helpers.SetModelClientAddress()
	helpers.SetModelBackendAddress()
	suite.Run(t, new(LikeControllerTestSuite))
}

// Mock DB Handler
type LikeDBTestHandler struct {
	CreateLikeFunc   func(*models.Like) (*models.Like, error)
	DeleteLikeFunc   func(string, string) error
	GetLikeCountFunc func(string) (uint64, error)
}

func (h *LikeDBTestHandler) CreateLike(newLike *models.Like) (*models.Like, error) {
	return h.CreateLikeFunc(newLike)
}

func (h *LikeDBTestHandler) DeleteLike(userID string, postID string) error {
	return h.DeleteLikeFunc(userID, postID)
}

func (h *LikeDBTestHandler) GetLikeCount(postID string) (uint64, error) {
	return h.GetLikeCountFunc(postID)
}

func (h *LikeDBTestHandler) SetMockCreateLikeFunc(like *models.Like, err error) {
	h.CreateLikeFunc = func(newLike *models.Like) (*models.Like, error) {
		return like, err
	}
}

func (h *LikeDBTestHandler) SetMockDeleteLikeFunc(err error) {
	h.DeleteLikeFunc = func(userID string, postID string) error {
		return err
	}
}

func (h *LikeDBTestHandler) SetMockGetLikeCountFunc(count uint64, err error) {
	h.GetLikeCountFunc = func(postID string) (uint64, error) {
		return count, err
	}
}

func (h *LikeDBTestHandler) ResetFuncs() {
	h.CreateLikeFunc = nil
	h.DeleteLikeFunc = nil
}

type LikeTestCache struct {
	GetCacheCountFunc func(context.Context, string) (uint64, error)
	SetCacheCountFunc func(context.Context, string) (uint64, error)
}

func (c *LikeTestCache) GetCacheCount(ctx context.Context, postID string) (uint64, error) {
	return c.GetCacheCountFunc(ctx, postID)
}

func (c *LikeTestCache) SetCacheCount(ctx context.Context, postID string) (uint64, error) {
	return c.SetCacheCountFunc(ctx, postID)
}

func (c *LikeTestCache) SetMockGetCacheCountFunc(count uint64, err error) {
	c.GetCacheCountFunc = func(ctx context.Context, postID string) (uint64, error) {
		return count, err
	}
}

func (c *LikeTestCache) SetMockSetCacheCountFunc(count uint64, err error) {
	c.SetCacheCountFunc = func(ctx context.Context, postID string) (uint64, error) {
		return count, err
	}
}

func (c *LikeTestCache) ResetFuncs() {
	c.GetCacheCountFunc = nil
	c.SetCacheCountFunc = nil
}

// Create Likes
func (s *LikeControllerTestSuite) Test_CreateLike_OK() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIdKey, testPostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPost, nil)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockCreateLikeFunc(&defaultLike, nil)
	s.cacheHandler.SetMockSetCacheCountFunc(1, nil)
	s.api.PostLike(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedDataEqualsActual(m, defaultCreateLikeUpdate); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *LikeControllerTestSuite) Test_CreateLike_InvalidPostIDString() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIdKey, invalidPostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPost, nil)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockCreateLikeFunc(&defaultLike, nil)
	s.api.PostLike(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusBadRequest, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrPostNotFound); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *LikeControllerTestSuite) Test_CreateLike_NegativePostID() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIdKey, negativePostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPost, nil)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockCreateLikeFunc(&defaultLike, nil)
	s.api.PostLike(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusBadRequest, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrPostNotFound); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *LikeControllerTestSuite) Test_CreateLike_DuplicateKey() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIdKey, testPostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPost, nil)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockCreateLikeFunc(&defaultLike, gorm.ErrDuplicatedKey)
	s.api.PostLike(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusBadRequest, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrAlreadyLiked); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *LikeControllerTestSuite) Test_CreateLike_CannotCreate() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIdKey, testPostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPost, nil)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockCreateLikeFunc(&defaultLike, ErrTest)
	s.api.PostLike(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusInternalServerError, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrLikeNotRegistered); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *LikeControllerTestSuite) Test_CreateLike_CannotSetCount() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIdKey, testPostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPost, nil)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockCreateLikeFunc(&defaultLike, nil)
	s.cacheHandler.SetMockSetCacheCountFunc(1, ErrTest)
	s.api.PostLike(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusInternalServerError, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrTest); !isEqual {
		s.T().Error(errStr)
	}
}

// Delete Likes
func (s *LikeControllerTestSuite) Test_DeleteLike_OK() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIdKey, testPostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodDelete, nil)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockDeleteLikeFunc(nil)
	s.cacheHandler.SetMockSetCacheCountFunc(0, nil)
	s.api.DeleteLike(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedDataEqualsActual(m, defaultDeleteLikeUpdate); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *LikeControllerTestSuite) Test_DeleteLike_InvalidPostIDString() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIdKey, invalidPostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodDelete, nil)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockDeleteLikeFunc(nil)
	s.api.DeleteLike(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusBadRequest, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrPostNotFound); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *LikeControllerTestSuite) Test_DeleteLike_NegativePostID() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIdKey, negativePostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodDelete, nil)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockDeleteLikeFunc(nil)
	s.api.PostLike(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusBadRequest, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrPostNotFound); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *LikeControllerTestSuite) Test_DeleteLike_NotFound() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIdKey, testPostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodDelete, nil)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockDeleteLikeFunc(gorm.ErrRecordNotFound)
	s.api.DeleteLike(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusBadRequest, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrPostNotFound); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *LikeControllerTestSuite) Test_DeleteLike_CannotDelete() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIdKey, testPostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodDelete, nil)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockDeleteLikeFunc(ErrTest)
	s.api.DeleteLike(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusInternalServerError, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrUnlikeFailed); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *LikeControllerTestSuite) Test_DeleteLike_CannotSetCount() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIdKey, testPostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodDelete, nil)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockDeleteLikeFunc(nil)
	s.cacheHandler.SetMockSetCacheCountFunc(0, ErrTest)
	s.api.DeleteLike(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusInternalServerError, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrTest); !isEqual {
		s.T().Error(errStr)
	}
}
