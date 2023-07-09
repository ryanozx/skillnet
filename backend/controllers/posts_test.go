package controllers

import (
	"io"
	"net/http"
	"testing"

	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

const (
	testPostID = 1
)

var (
	defaultPostView = models.PostView{
		Post: models.Post{
			Content: "Hello world!",
		},
		UserMinimal: *defaultUserView.GetUserMinimal(),
	}
	newTestPost = models.Post{
		Content: "Hello world!",
	}
)

type PostControllerTestSuite struct {
	suite.Suite
	api       APIEnv
	dbHandler *PostDBTestHandler
	store     *helpers.MockSessionStore
}

func (s *PostControllerTestSuite) SetupSuite() {
	dbHandler := PostDBTestHandler{}
	api := APIEnv{
		PostDBHandler: &dbHandler,
	}
	s.api = api
	s.dbHandler = &dbHandler
	s.store = helpers.MakeMockStore()
}

func (s *PostControllerTestSuite) TearDownTest() {
	s.dbHandler.ResetFuncs()
	s.store.Reset()
}

func TestPostControllerSuite(t *testing.T) {
	t.Setenv("CLIENT_HOST", "http://localhost")
	t.Setenv("CLIENT_PORT", "3000")
	t.Setenv("BACKEND_HOST", "http://localhost")
	t.Setenv("BACKEND_PORT", "8080")
	helpers.SetModelClientAddress()
	helpers.SetModelBackendAddress()
	suite.Run(t, new(PostControllerTestSuite))
}

type PostDBTestHandler struct {
	CreatePostFunc  func(*models.Post) (*models.PostView, error)
	DeletePostFunc  func(string, string) error
	GetPostsFunc    func(string, string) ([]models.PostView, uint, error)
	GetPostByIDFunc func(string, string) (*models.PostView, error)
	UpdatePostFunc  func(*models.Post, string, string) (*models.PostView, error)
}

func (h *PostDBTestHandler) CreatePost(newPost *models.Post) (*models.PostView, error) {
	post, err := h.CreatePostFunc(newPost)
	return post, err
}

func (h *PostDBTestHandler) DeletePost(id string, userid string) error {
	err := h.DeletePostFunc(id, userid)
	return err
}

func (h *PostDBTestHandler) GetPosts(cutoff string, userID string) ([]models.PostView, uint, error) {
	posts, newCutoff, err := h.GetPostsFunc(cutoff, userID)
	return posts, newCutoff, err
}

func (h *PostDBTestHandler) GetPostByID(id string, userID string) (*models.PostView, error) {
	post, err := h.GetPostByIDFunc(id, userID)
	return post, err
}

func (h *PostDBTestHandler) UpdatePost(post *models.Post, postID string, userID string) (*models.PostView, error) {
	postView, err := h.UpdatePostFunc(post, postID, userID)
	return postView, err
}

func (h *PostDBTestHandler) SetMockCreatePostFunc(post *models.PostView, err error) {
	h.CreatePostFunc = func(newPost *models.Post) (*models.PostView, error) {
		return post, err
	}
}

func (h *PostDBTestHandler) SetMockDeletePostFunc(err error) {
	h.DeletePostFunc = func(id string, userid string) error {
		return err
	}
}

func (h *PostDBTestHandler) SetMockGetPostsFunc(posts []models.PostView, newCutoff uint, err error) {
	h.GetPostsFunc = func(cutoff string, userID string) ([]models.PostView, uint, error) {
		return posts, newCutoff, err
	}
}

func (h *PostDBTestHandler) SetMockGetPostByIDFunc(post *models.PostView, err error) {
	h.GetPostByIDFunc = func(id string, userID string) (*models.PostView, error) {
		return post, err
	}
}

func (h *PostDBTestHandler) SetMockUpdatePostFunc(postView *models.PostView, err error) {
	h.UpdatePostFunc = func(post *models.Post, postID string, userID string) (*models.PostView, error) {
		return postView, err
	}
}

func (h *PostDBTestHandler) ResetFuncs() {
	h.CreatePostFunc = nil
	h.DeletePostFunc = nil
	h.GetPostsFunc = nil
	h.GetPostByIDFunc = nil
	h.UpdatePostFunc = nil
}

// CreatePosts
func (s *PostControllerTestSuite) Test_CreatePost_OK() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPost, newTestPost)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockCreatePostFunc(&defaultPostView, nil)
	s.api.CreatePost(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedDataEqualsActual(m, defaultPostView); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *PostControllerTestSuite) Test_CreatePost_BadRequest() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)

	s.dbHandler.SetMockCreatePostFunc(&defaultPostView, nil)
	s.api.CreatePost(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusBadRequest, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrBadBinding); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *PostControllerTestSuite) Test_CreatePost_CannotCreate() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPost, newTestPost)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockCreatePostFunc(&defaultPostView, ErrTest)
	s.api.CreatePost(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusInternalServerError, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrCannotCreatePost); !isEqual {
		s.T().Error(errStr)
	}
}

// Delete Post
func (s *PostControllerTestSuite) Test_DeletePost_OK() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)
	helpers.AddParamsToContext(c, "postid", testPostID)

	s.dbHandler.SetMockDeletePostFunc(nil)
	s.api.DeletePost(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedMessageEqualsActual(m, PostDeletedMsg); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *PostControllerTestSuite) Test_DeletePost_CannotFindPost() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)
	helpers.AddParamsToContext(c, "postid", testPostID)

	s.dbHandler.SetMockDeletePostFunc(gorm.ErrRecordNotFound)
	s.api.DeletePost(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusNotFound, w.Code); !isEqual {
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

func (s *PostControllerTestSuite) Test_DeletePost_NotOwner() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)
	helpers.AddParamsToContext(c, "postid", testPostID)

	s.dbHandler.SetMockDeletePostFunc(database.ErrNotOwner)
	s.api.DeletePost(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusUnauthorized, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, database.ErrNotOwner); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *PostControllerTestSuite) Test_DeletePost_CannotDelete() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)
	helpers.AddParamsToContext(c, "postid", testPostID)

	s.dbHandler.SetMockDeletePostFunc(ErrTest)
	s.api.DeletePost(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusBadRequest, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrCannotDeletePost); !isEqual {
		s.T().Error(errStr)
	}
}

// GetPosts
func (s *PostControllerTestSuite) Test_GetPosts_OK() {
	helpers.SetModelBackendAddress()
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)

	expectedPostViews := []models.PostView{defaultPostView}
	expectedArr := models.PostViewArray{
		Posts:       expectedPostViews,
		NextPageURL: "http://localhost:8080/auth/posts?cutoff=0",
	}

	s.dbHandler.SetMockGetPostsFunc(expectedPostViews, 0, nil)
	s.api.GetPosts(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedDataEqualsActual(m, expectedArr); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *PostControllerTestSuite) Test_GetPosts_DiffCutoff() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)

	expectedPostViews := []models.PostView{defaultPostView}
	expectedArr := models.PostViewArray{
		Posts:       expectedPostViews,
		NextPageURL: "http://localhost:8080/auth/posts?cutoff=10",
	}

	s.dbHandler.SetMockGetPostsFunc(expectedPostViews, 10, nil)
	s.api.GetPosts(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedDataEqualsActual(m, expectedArr); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *PostControllerTestSuite) Test_GetPosts_NotFound() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)

	expected := []models.PostView{defaultPostView}

	s.dbHandler.SetMockGetPostsFunc(expected, 0, ErrTest)
	s.api.GetPosts(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusNotFound, w.Code); !isEqual {
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

// GetPostByID
func (s *PostControllerTestSuite) Test_GetPostByID_OK() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, "postid", testPostID)

	s.dbHandler.SetMockGetPostByIDFunc(&defaultPostView, nil)
	s.api.GetPostByID(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedDataEqualsActual(m, defaultPostView); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *PostControllerTestSuite) Test_GetPostByID_NotFound() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, "postid", testPostID)

	s.dbHandler.SetMockGetPostByIDFunc(&defaultPostView, ErrTest)
	s.api.GetPostByID(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusNotFound, w.Code); !isEqual {
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

// Update Posts
func (s *PostControllerTestSuite) Test_UpdatePost_OK() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)
	helpers.AddParamsToContext(c, "postid", testPostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPatch, newTestPost)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockUpdatePostFunc(&defaultPostView, nil)
	s.api.UpdatePost(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedDataEqualsActual(m, defaultPostView); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *PostControllerTestSuite) Test_UpdatePost_BadRequest() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)
	helpers.AddParamsToContext(c, "postid", testPostID)

	s.dbHandler.SetMockUpdatePostFunc(&defaultPostView, nil)
	s.api.UpdatePost(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusBadRequest, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrBadBinding); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *PostControllerTestSuite) Test_UpdatePost_NotFound() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)
	helpers.AddParamsToContext(c, "postid", testPostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPatch, newTestPost)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockUpdatePostFunc(&defaultPostView, gorm.ErrRecordNotFound)
	s.api.UpdatePost(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusNotFound, w.Code); !isEqual {
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

func (s *PostControllerTestSuite) Test_UpdatePost_NotOwner() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)
	helpers.AddParamsToContext(c, "postid", testPostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPatch, newTestPost)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockUpdatePostFunc(&defaultPostView, database.ErrNotOwner)
	s.api.UpdatePost(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusUnauthorized, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, database.ErrNotOwner); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *PostControllerTestSuite) Test_UpdatePost_CannotUpdate() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIdKey, testUserID)
	helpers.AddParamsToContext(c, "postid", testPostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPatch, newTestPost)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockUpdatePostFunc(&defaultPostView, ErrTest)
	s.api.UpdatePost(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusBadRequest, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrCannotUpdatePost); !isEqual {
		s.T().Error(errStr)
	}
}
