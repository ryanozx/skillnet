package controllers

import (
	"io"
	"net/http"
	"testing"

	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
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
		LikeCount: 1,
	}
	defaultDeleteLikeUpdate = models.LikeUpdate{
		LikeCount: 0,
	}
)

type LikeControllerTestSuite struct {
	suite.Suite
	api                APIEnv
	dbHandler          *LikeDBTestHandler
	cacheHandler       *helpers.TestCache
	store              *helpers.MockSessionStore
	notificationPoster *helpers.TestNotificationCreator
}

func (s *LikeControllerTestSuite) SetupSuite() {
	dbHandler := LikeDBTestHandler{}
	cacheHandler := helpers.TestCache{}
	notificationCreator := helpers.TestNotificationCreator{}
	api := APIEnv{
		LikeDBHandler:      &dbHandler,
		LikesCacheHandler:  &cacheHandler,
		NotificationPoster: &notificationCreator,
	}
	s.api = api
	s.dbHandler = &dbHandler
	s.cacheHandler = &cacheHandler
	s.store = helpers.MakeMockStore()
	s.notificationPoster = &notificationCreator
}

func (s *LikeControllerTestSuite) TearDownTest() {
	s.dbHandler.ResetFuncs()
	s.store.Reset()
	s.notificationPoster.ResetFuncs()
}

func TestLikeControllerSuite(t *testing.T) {
	helpers.SetEnvVars(t)
	suite.Run(t, new(LikeControllerTestSuite))
}

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

func (h *LikeDBTestHandler) ResetFuncs() {
	h.CreateLikeFunc = nil
	h.DeleteLikeFunc = nil
}

// Create Likes
func (s *LikeControllerTestSuite) Test_CreateLike_OK() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIDKey, testPostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPost, nil)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockCreateLikeFunc(&defaultLike, nil)
	s.cacheHandler.SetMockSetCacheValFunc(1, nil)
	s.notificationPoster.SetMockPostNotificationFromEventFunc(nil)
	s.api.PostLike(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedDataEqualsActual[models.LikeUpdate](m, &defaultCreateLikeUpdate); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *LikeControllerTestSuite) Test_CreateLike_InvalidPostIDString() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIDKey, invalidPostID)

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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIDKey, negativePostID)

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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIDKey, testPostID)

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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIDKey, testPostID)

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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIDKey, testPostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPost, nil)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockCreateLikeFunc(&defaultLike, nil)
	s.cacheHandler.SetMockSetCacheValFunc(1, ErrTest)
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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIDKey, testPostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodDelete, nil)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockDeleteLikeFunc(nil)
	s.cacheHandler.SetMockSetCacheValFunc(0, nil)
	s.api.DeleteLike(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedDataEqualsActual[models.LikeUpdate](m, &defaultDeleteLikeUpdate); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *LikeControllerTestSuite) Test_DeleteLike_InvalidPostIDString() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIDKey, invalidPostID)

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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIDKey, negativePostID)

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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIDKey, testPostID)

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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIDKey, testPostID)

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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIDKey, testPostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodDelete, nil)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockDeleteLikeFunc(nil)
	s.cacheHandler.SetMockSetCacheValFunc(0, ErrTest)
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
