package controllers

import (
	"io"
	"net/http"
	"testing"

	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

var (
	defaultLoginUserDBEntry = models.User{
		UserCredentials: defaultCreds,
	}
)

type AuthControllerTestSuite struct {
	suite.Suite
	api       APIEnv
	dbHandler *UserDBTestHandler
	store     *helpers.MockSessionStore
}

func (s *AuthControllerTestSuite) SetupSuite() {
	dbHandler := UserDBTestHandler{}
	api := APIEnv{
		AuthDBHandler: &dbHandler,
	}
	s.api = api
	s.dbHandler = &dbHandler
	s.store = helpers.MakeMockStore()
}

func (s *AuthControllerTestSuite) TearDownTest() {
	s.dbHandler.ResetFuncs()
	s.store.Reset()
}

func TestAuthControllerSuite(t *testing.T) {
	t.Setenv("CLIENT_HOST", "http://localhost")
	t.Setenv("CLIENT_PORT", "3000")
	t.Setenv("BACKEND_HOST", "http://localhost")
	t.Setenv("BACKEND_PORT", "8080")
	helpers.SetModelClientAddress()
	helpers.SetModelBackendAddress()
	suite.Run(t, new(AuthControllerTestSuite))
}

func (s *AuthControllerTestSuite) Test_GetLogin_OK() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)

	s.api.GetLogin(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedMessageEqualsActual(m, GetLoginOKMsg); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *AuthControllerTestSuite) Test_GetLogin_ExistingSession() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	s.store.Set(helpers.IdKey, testUserID)

	s.api.GetLogin(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusBadRequest, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrAlreadyLoggedIn); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *AuthControllerTestSuite) Test_PostLogin_OK() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)

	expected := defaultLoginUserDBEntry
	hash, err := bcrypt.GenerateFromPassword([]byte(defaultCreds.Password), bcrypt.DefaultCost)
	if err != nil {
		s.T().Error(err)
	}
	expected.Password = string(hash)

	c.Request = helpers.GenerateHttpFormDataRequest(http.MethodPost, defaultCreds)
	s.dbHandler.SetMockGetUserByUsernameFunc(&expected, nil)
	s.api.PostLogin(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedMessageEqualsActual(m, LoginSuccessfulMsg); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *AuthControllerTestSuite) Test_PostLogin_AlreadyLoggedIn() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	s.store.Set(helpers.IdKey, testUserID)

	expected := defaultLoginUserDBEntry
	hash, err := bcrypt.GenerateFromPassword([]byte(defaultCreds.Password), bcrypt.DefaultCost)
	if err != nil {
		s.T().Error(err)
	}
	expected.Password = string(hash)

	c.Request = helpers.GenerateHttpFormDataRequest(http.MethodPost, defaultCreds)
	s.dbHandler.SetMockGetUserByUsernameFunc(&expected, nil)
	s.api.PostLogin(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusBadRequest, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrAlreadyLoggedIn); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *AuthControllerTestSuite) Test_PostLogin_EmptyUserPass() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)

	expected := defaultLoginUserDBEntry
	hash, err := bcrypt.GenerateFromPassword([]byte(defaultCreds.Password), bcrypt.DefaultCost)
	if err != nil {
		s.T().Error(err)
	}
	expected.Password = string(hash)

	userCreds := models.UserCredentials{
		Username: "",
		Password: "",
	}
	c.Request = helpers.GenerateHttpFormDataRequest(http.MethodPost, userCreds)

	s.dbHandler.SetMockGetUserByUsernameFunc(&expected, nil)
	s.api.PostLogin(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusBadRequest, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrMissingUserCredentials); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *AuthControllerTestSuite) Test_PostLogin_InvalidUser() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)

	expected := defaultLoginUserDBEntry
	hash, err := bcrypt.GenerateFromPassword([]byte(defaultCreds.Password), bcrypt.DefaultCost)
	if err != nil {
		s.T().Error(err)
	}
	expected.Password = string(hash)

	c.Request = helpers.GenerateHttpFormDataRequest(http.MethodPost, defaultCreds)
	s.dbHandler.SetMockGetUserByUsernameFunc(&expected, errTest)
	s.api.PostLogin(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusUnauthorized, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrIncorrectUserCredentials); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *AuthControllerTestSuite) Test_PostLogin_WrongPassword() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)

	expected := defaultLoginUserDBEntry
	hash, err := bcrypt.GenerateFromPassword([]byte(defaultCreds.Password), bcrypt.DefaultCost)
	if err != nil {
		s.T().Error(err)
	}
	expected.Password = string(hash)

	badUserCreds := models.UserCredentials{
		Username: defaultCreds.Username,
		Password: "1234",
	}

	c.Request = helpers.GenerateHttpFormDataRequest(http.MethodPost, badUserCreds)
	s.dbHandler.SetMockGetUserByUsernameFunc(&expected, errTest)
	s.api.PostLogin(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusUnauthorized, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrIncorrectUserCredentials); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *AuthControllerTestSuite) Test_PostLogin_CannotSave() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	s.store.SetSaveError(errTest)

	expected := defaultLoginUserDBEntry
	hash, err := bcrypt.GenerateFromPassword([]byte(defaultCreds.Password), bcrypt.DefaultCost)
	if err != nil {
		s.T().Error(err)
	}
	expected.Password = string(hash)

	c.Request = helpers.GenerateHttpFormDataRequest(http.MethodPost, defaultCreds)
	s.dbHandler.SetMockGetUserByUsernameFunc(&expected, nil)
	s.api.PostLogin(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusInternalServerError, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrCookieSaveFail); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *AuthControllerTestSuite) Test_PostLogout_OK() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	s.store.Set(helpers.IdKey, testUserID)

	s.api.PostLogout(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusOK, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedMessageEqualsActual(m, SuccessfulLogoutMsg); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *AuthControllerTestSuite) Test_PostLogout_InvalidSession() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)

	s.api.PostLogout(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusBadRequest, w.Code); !isEqual {
		s.T().Error(errStr)
	}
	m, err := helpers.ParseJSONString(b)
	if err != nil {
		s.T().Error(err)
	}
	if errStr, isEqual := helpers.CheckExpectedErrorEqualsActual(m, ErrNoValidSession); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *AuthControllerTestSuite) Test_PostLogout_CannotSave() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	s.store.Set(helpers.IdKey, testUserID)
	s.store.SetSaveError(errTest)

	s.api.PostLogout(c)

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
