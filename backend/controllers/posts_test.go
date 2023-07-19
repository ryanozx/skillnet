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

const (
	testPostID     = 1
	invalidPostID  = "badpostid"
	negativePostID = -234
)

var (
	defaultPostView = models.PostView{
		Post: models.Post{
			Content: "Hello world!",
			User:    defaultUser,
		},
		UserMinimal: defaultUserMinimal,
	}
	defaultPost = models.Post{
		Content: "Hello world!",
		User:    defaultUser,
	}
	diffCutoffPost = models.Post{
		Model: gorm.Model{
			ID: 10,
		},
		Content: "Hello world!",
		UserID:  testUserID,
		User:    defaultUser,
	}
	newTestPost = models.Post{
		Content: "Hello world!",
	}
)

type PostControllerTestSuite struct {
	suite.Suite
	api                  APIEnv
	dbHandler            *PostDBTestHandler
	likesCacheHandler    *helpers.TestCache
	commentsCacheHandler *helpers.TestCache
	store                *helpers.MockSessionStore
}

func (s *PostControllerTestSuite) SetupSuite() {
	dbHandler := PostDBTestHandler{}
	likesCacheHandler := helpers.TestCache{}
	commentsCacheHandler := helpers.TestCache{}
	api := APIEnv{
		PostDBHandler:        &dbHandler,
		LikesCacheHandler:    &likesCacheHandler,
		CommentsCacheHandler: &commentsCacheHandler,
	}
	s.api = api
	s.dbHandler = &dbHandler
	s.likesCacheHandler = &likesCacheHandler
	s.commentsCacheHandler = &commentsCacheHandler
	s.store = helpers.MakeMockStore()
}

func (s *PostControllerTestSuite) TearDownTest() {
	s.dbHandler.ResetFuncs()
	s.likesCacheHandler.ResetFuncs()
	s.store.Reset()
}

func TestPostControllerSuite(t *testing.T) {
	suite.Run(t, new(PostControllerTestSuite))
}

func (s *PostControllerTestSuite) SetupTest() {
	s.T().Setenv("CLIENT_HOST", "http://localhost")
	s.T().Setenv("CLIENT_PORT", "3000")
	s.T().Setenv("BACKEND_HOST", "http://localhost")
	s.T().Setenv("BACKEND_PORT", "8080")
	helpers.SetModelClientAddress()
	helpers.SetModelBackendAddress()
}

type PostDBTestHandler struct {
	CreatePostFunc  func(*models.Post) (*models.Post, error)
	DeletePostFunc  func(uint, string) error
	GetPostsFunc    func(*helpers.NullableUint, string) ([]models.Post, error)
	GetPostByIDFunc func(uint, string) (*models.Post, error)
	UpdatePostFunc  func(*models.Post, uint, string) (*models.Post, error)
}

func (h *PostDBTestHandler) CreatePost(newPost *models.Post) (*models.Post, error) {
	post, err := h.CreatePostFunc(newPost)
	return post, err
}

func (h *PostDBTestHandler) DeletePost(postID uint, userid string) error {
	err := h.DeletePostFunc(postID, userid)
	return err
}

func (h *PostDBTestHandler) GetPosts(cutoff *helpers.NullableUint, userID string) ([]models.Post, error) {
	posts, err := h.GetPostsFunc(cutoff, userID)
	return posts, err
}

func (h *PostDBTestHandler) GetPostByID(postID uint, userID string) (*models.Post, error) {
	post, err := h.GetPostByIDFunc(postID, userID)
	return post, err
}

func (h *PostDBTestHandler) UpdatePost(post *models.Post, postID uint, userID string) (*models.Post, error) {
	updatedPost, err := h.UpdatePostFunc(post, postID, userID)
	return updatedPost, err
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
	h.GetPostsFunc = func(cutoff *helpers.NullableUint, userID string) ([]models.Post, error) {
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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPost, newTestPost)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockCreatePostFunc(&defaultPost, nil)
	s.api.CreatePost(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedDataEqualsActual[models.PostView](m, &defaultPostView); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *PostControllerTestSuite) Test_CreatePost_BadRequest() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)

	s.dbHandler.SetMockCreatePostFunc(&defaultPost, nil)
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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPost, newTestPost)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockCreatePostFunc(&defaultPost, ErrTest)
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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIDKey, testPostID)

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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIDKey, testPostID)

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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIDKey, testPostID)

	s.dbHandler.SetMockDeletePostFunc(helpers.ErrNotOwner)
	s.api.DeletePost(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusUnauthorized, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, helpers.ErrNotOwner); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *PostControllerTestSuite) Test_DeletePost_CannotDelete() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIDKey, testPostID)

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
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)

	expectedPosts := []models.Post{defaultPost}
	expectedArr := models.PostViewArray{
		Posts:       []models.PostView{defaultPostView},
		NextPageURL: "http://localhost:8080/auth/posts?cutoff=0",
	}

	s.dbHandler.SetMockGetPostsFunc(expectedPosts, nil)
	s.likesCacheHandler.SetMockGetCacheValFunc(0, nil)
	s.commentsCacheHandler.SetMockGetCacheValFunc(0, nil)
	s.api.GetPosts(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedDataEqualsActual[models.PostViewArray](m, &expectedArr); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *PostControllerTestSuite) Test_GetPosts_DiffCutoff() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)

	expectedPosts := []models.Post{diffCutoffPost}
	expectedArr := models.PostViewArray{
		Posts: []models.PostView{*diffCutoffPost.PostView(&models.PostViewParams{
			UserID:       testUserID,
			CommentCount: 0,
			LikeCount:    0,
		})},
		NextPageURL: "http://localhost:8080/auth/posts?cutoff=10",
	}

	s.dbHandler.SetMockGetPostsFunc(expectedPosts, nil)
	s.likesCacheHandler.SetMockGetCacheValFunc(0, nil)
	s.commentsCacheHandler.SetMockGetCacheValFunc(0, nil)
	s.api.GetPosts(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedDataEqualsActual[models.PostViewArray](m, &expectedArr); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *PostControllerTestSuite) Test_GetPosts_NotFound() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)

	expected := []models.Post{defaultPost}

	s.dbHandler.SetMockGetPostsFunc(expected, ErrTest)
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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIDKey, testPostID)

	s.dbHandler.SetMockGetPostByIDFunc(&defaultPost, nil)
	s.likesCacheHandler.SetMockGetCacheValFunc(0, nil)
	s.commentsCacheHandler.SetMockGetCacheValFunc(0, nil)
	s.api.GetPostByID(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedDataEqualsActual[models.PostView](m, &defaultPostView); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *PostControllerTestSuite) Test_GetPostByID_NotFound() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.PostIDKey, testPostID)

	s.dbHandler.SetMockGetPostByIDFunc(&defaultPost, ErrTest)
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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIDKey, testPostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPatch, newTestPost)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockUpdatePostFunc(&defaultPost, nil)
	s.likesCacheHandler.SetMockGetCacheValFunc(0, nil)
	s.api.UpdatePost(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedDataEqualsActual[models.PostView](m, &defaultPostView); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *PostControllerTestSuite) Test_UpdatePost_BadRequest() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIDKey, testPostID)

	s.dbHandler.SetMockUpdatePostFunc(&defaultPost, nil)
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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIDKey, testPostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPatch, newTestPost)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockUpdatePostFunc(&defaultPost, gorm.ErrRecordNotFound)
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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIDKey, testPostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPatch, newTestPost)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockUpdatePostFunc(&defaultPost, helpers.ErrNotOwner)
	s.api.UpdatePost(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusUnauthorized, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, helpers.ErrNotOwner); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *PostControllerTestSuite) Test_UpdatePost_CannotUpdate() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)
	helpers.AddParamsToContext(c, helpers.PostIDKey, testPostID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPatch, newTestPost)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req
	s.dbHandler.SetMockUpdatePostFunc(&defaultPost, ErrTest)
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
