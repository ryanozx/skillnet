package controllers

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"github.com/stretchr/testify/suite"
	"gopkg.in/guregu/null.v3"
	"gorm.io/gorm"
)

const (
	testUserID = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
)

var (
	// errTest is a test error that can be used to simulate an unexpected error returned by the database helper functions
	errTest      = errors.New("test error")
	defaultCreds = models.UserCredentials{
		Username: "testuser",
		Password: "12345",
	}
	defaultUser = models.User{
		UserCredentials: models.UserCredentials{
			Username: "testuser",
		},
		UserView: models.UserView{
			Title: null.NewString("Tester", true),
			UserMinimal: models.UserMinimal{
				Name: "Test User",
				URL:  "localhost:8080/user/testuser",
			},
		},
	}
	defaultUserView = defaultUser.GetUserView()
)

type UserControllerTestSuite struct {
	suite.Suite
	api       APIEnv
	dbHandler *UserDBTestHandler
	store     *helpers.MockSessionStore
}

func (s *UserControllerTestSuite) SetupSuite() {
	dbHandler := UserDBTestHandler{}
	api := APIEnv{
		UserDBHandler: &dbHandler,
	}
	s.api = api
	s.dbHandler = &dbHandler
	s.store = helpers.MakeMockStore()
}

func (s *UserControllerTestSuite) TearDownTest() {
	s.dbHandler.ResetFuncs()
	s.store.Reset()
}

func TestUserControllerSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}

type UserDBTestHandler struct {
	CreateUserFunc        func(database.NewUser) (*models.User, error)
	DeleteUserFunc        func(string) error
	GetUserByIDFunc       func(string) (*models.User, error)
	GetUserByUsernameFunc func(string) (*models.User, error)
	UpdateUserFunc        func(*models.User, string) (*models.User, error)
}

func (h *UserDBTestHandler) CreateUser(newUser database.NewUser) (*models.User, error) {
	user, err := h.CreateUserFunc(newUser)
	return user, err
}

func (h *UserDBTestHandler) DeleteUser(id string) error {
	err := h.DeleteUserFunc(id)
	return err
}

func (h *UserDBTestHandler) GetUserByID(id string) (*models.User, error) {
	user, err := h.GetUserByIDFunc(id)
	return user, err
}

func (h *UserDBTestHandler) GetUserByUsername(username string) (*models.User, error) {
	user, err := h.GetUserByUsernameFunc(username)
	return user, err
}

func (h *UserDBTestHandler) UpdateUser(user *models.User, id string) (*models.User, error) {
	updatedUser, err := h.UpdateUserFunc(user, id)
	return updatedUser, err
}

func (h *UserDBTestHandler) SetMockCreateUserFunc(user *models.User, err error) {
	h.CreateUserFunc = func(newUser database.NewUser) (*models.User, error) {
		return user, err
	}
}

func (h *UserDBTestHandler) SetMockDeleteUserFunc(err error) {
	h.DeleteUserFunc = func(id string) error {
		return err
	}
}

func (h *UserDBTestHandler) SetMockGetUserByIDFunc(user *models.User, err error) {
	h.GetUserByIDFunc = func(id string) (*models.User, error) {
		return user, err
	}
}

func (h *UserDBTestHandler) SetMockGetUserByUsernameFunc(user *models.User, err error) {
	h.GetUserByUsernameFunc = func(username string) (*models.User, error) {
		return user, err
	}
}

func (h *UserDBTestHandler) SetMockUpdateUserFunc(user *models.User, err error) {
	h.UpdateUserFunc = func(u *models.User, id string) (*models.User, error) {
		return user, err
	}
}

func (h *UserDBTestHandler) ResetFuncs() {
	h.CreateUserFunc = nil
	h.DeleteUserFunc = nil
	h.GetUserByIDFunc = nil
	h.GetUserByUsernameFunc = nil
	h.UpdateUserFunc = nil
}

// CreateUser
func (s *UserControllerTestSuite) Test_CreateUser_OK() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)

	c.Request = helpers.GenerateHttpFormDataRequest(http.MethodPost, struct {
		Username string
		Password string
		Email    string
	}{"testuser", "12345", "abc@gmail.com"})
	s.dbHandler.SetMockCreateUserFunc(&defaultUser, nil)
	s.api.CreateUser(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedMessageEqualsActual(m, SuccessfulAccountCreationMsg); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *UserControllerTestSuite) Test_CreateUser_EmptyUsername() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)

	c.Request = helpers.GenerateHttpFormDataRequest(http.MethodPost, struct {
		Username string
		Password string
		Email    string
	}{"", "12345", "abc@gmail.com"})

	s.dbHandler.SetMockCreateUserFunc(&defaultUser, nil)
	s.api.CreateUser(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusBadRequest, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrMissingSignupCredentials); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *UserControllerTestSuite) Test_CreateUser_EmptyPassword() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)

	c.Request = helpers.GenerateHttpFormDataRequest(http.MethodPost, struct {
		Username string
		Password string
		Email    string
	}{"testuser", "", "abc@gmail.com"})

	s.dbHandler.SetMockCreateUserFunc(&defaultUser, nil)
	s.api.CreateUser(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusBadRequest, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrMissingSignupCredentials); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *UserControllerTestSuite) Test_CreateUser_EmptyEmail() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)

	c.Request = helpers.GenerateHttpFormDataRequest(http.MethodPost, struct {
		Username string
		Password string
		Email    string
	}{"testuser", "12345", ""})

	s.dbHandler.SetMockCreateUserFunc(&defaultUser, nil)
	s.api.CreateUser(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusBadRequest, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrMissingSignupCredentials); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *UserControllerTestSuite) Test_CreateUser_UsernameConflict() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)

	c.Request = helpers.GenerateHttpFormDataRequest(http.MethodPost, struct {
		Username string
		Password string
		Email    string
	}{"testuser", "12345", "abc@gmail.com"})
	s.dbHandler.SetMockCreateUserFunc(&defaultUser, gorm.ErrDuplicatedKey)
	s.api.CreateUser(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusConflict, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrUsernameAlreadyExists); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *UserControllerTestSuite) Test_CreateUser_CannotCreate() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)

	c.Request = helpers.GenerateHttpFormDataRequest(http.MethodPost, struct {
		Username string
		Password string
		Email    string
	}{"testuser", "12345", "abc@gmail.com"})
	s.dbHandler.SetMockCreateUserFunc(&defaultUser, errTest)
	s.api.CreateUser(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusInternalServerError, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, errTest); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *UserControllerTestSuite) Test_CreateUser_CannotSaveSession() {
	c, w := helpers.CreateTestContextAndRecorder()
	s.store.SetSaveError(errTest)
	helpers.AddStoreToContext(c, s.store)

	c.Request = helpers.GenerateHttpFormDataRequest(http.MethodPost, struct {
		Username string
		Password string
		Email    string
	}{"testuser", "12345", "abc@gmail.com"})
	s.dbHandler.SetMockCreateUserFunc(&defaultUser, nil)
	s.api.CreateUser(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusInternalServerError, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrCreateAccountNoCookie); !isEqual {
		s.T().Error(errStr)
	}
}

// DeleteUser
func (s *UserControllerTestSuite) Test_DeleteUser_OK() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.IdKey, testUserID)

	s.dbHandler.SetMockDeleteUserFunc(nil)
	s.api.DeleteUser(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedMessageEqualsActual(m, SuccessfulAccountDeleteMsg); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *UserControllerTestSuite) Test_DeleteUser_NotFound() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.IdKey, testUserID)

	s.dbHandler.SetMockDeleteUserFunc(gorm.ErrRecordNotFound)
	s.api.DeleteUser(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusNotFound, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrUserNotFound); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *UserControllerTestSuite) Test_DeleteUser_CannotDelete() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.IdKey, testUserID)

	s.dbHandler.SetMockDeleteUserFunc(errTest)
	s.api.DeleteUser(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusBadRequest, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, errTest); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *UserControllerTestSuite) Test_DeleteUser_CannotClearSession() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	s.store.SetSaveError(errTest)
	helpers.AddParamsToContext(c, helpers.IdKey, testUserID)

	s.dbHandler.SetMockDeleteUserFunc(nil)
	s.api.DeleteUser(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusInternalServerError, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrSessionClearFailed); !isEqual {
		s.T().Error(errStr)
	}
}

// GetProfile
func (s *UserControllerTestSuite) Test_GetProfile_OK() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)

	s.dbHandler.SetMockGetUserByUsernameFunc(&defaultUser, nil)
	s.api.GetProfile(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedDataEqualsActual[models.UserView](m, *defaultUserView); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *UserControllerTestSuite) Test_GetProfile_NotFound() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)

	s.dbHandler.SetMockGetUserByUsernameFunc(&defaultUser, errTest)
	s.api.GetProfile(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusNotFound, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrUserNotFound); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *UserControllerTestSuite) Test_GetSelfProfile_OK() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.IdKey, testUserID)

	s.dbHandler.SetMockGetUserByIDFunc(&defaultUser, nil)
	s.api.GetSelfProfile(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedDataEqualsActual[models.User](m, defaultUser); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *UserControllerTestSuite) Test_GetSelfProfile_NotFound() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.IdKey, testUserID)

	s.dbHandler.SetMockGetUserByIDFunc(&defaultUser, errTest)
	s.api.GetSelfProfile(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusNotFound, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrUserNotFound); !isEqual {
		s.T().Error(errStr)
	}
}

// UpdateUser
func (s *UserControllerTestSuite) Test_UpdateUser_OK() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.IdKey, testUserID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPatch, defaultUser)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req

	s.dbHandler.SetMockUpdateUserFunc(&defaultUser, nil)
	s.api.UpdateUser(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedDataEqualsActual[models.User](m, defaultUser); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *UserControllerTestSuite) Test_UpdateUser_BadRequest() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.IdKey, testUserID)

	s.dbHandler.SetMockUpdateUserFunc(&defaultUser, nil)
	s.api.UpdateUser(c)

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

func (s *UserControllerTestSuite) Test_UpdateUser_NotFound() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.IdKey, testUserID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPatch, defaultUser)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req

	s.dbHandler.SetMockUpdateUserFunc(&defaultUser, gorm.ErrRecordNotFound)
	s.api.UpdateUser(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusNotFound, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrUserNotFound); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *UserControllerTestSuite) Test_UpdateUser_CannotUpdate() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.IdKey, testUserID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPatch, defaultUser)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req

	s.dbHandler.SetMockUpdateUserFunc(&defaultUser, errTest)
	s.api.UpdateUser(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusInternalServerError, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrCannotUpdateUser); !isEqual {
		s.T().Error(errStr)
	}
}
