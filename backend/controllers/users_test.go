package controllers

import (
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
	testUserID   = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	testUsername = "testuser"
	testPassword = "12345"
	testEmail    = "abc@def.com"
	diffUserID   = "6ba7b812-9dad-11d1-80b4-00c04fd430c"
)

var (
	defaultCreds = models.UserCredentials{
		Username: testUsername,
		Password: testPassword,
	}
	defaultUser = models.User{
		UserCredentials: defaultCreds,
		UserView: models.UserView{
			Title:       null.NewString("Tester", true),
			UserMinimal: defaultUserMinimal,
		},
	}
	defaultUserMinimal = models.UserMinimal{
		Name: "Test User",
		URL:  "http://localhost:3000/profile/testuser",
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
	t.Setenv("CLIENT_HOST", "http://localhost")
	t.Setenv("CLIENT_PORT", "3000")
	t.Setenv("BACKEND_HOST", "http://localhost")
	t.Setenv("BACKEND_PORT", "8080")
	helpers.SetModelClientAddress()
	helpers.SetModelBackendAddress()
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

func (h *UserDBTestHandler) QueryUser(searchTerm string, limit int) ([]models.SearchResult, error) {
	return nil, nil
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

// DeleteUser
func (s *UserControllerTestSuite) Test_DeleteUser_OK() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)

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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)

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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)

	s.dbHandler.SetMockDeleteUserFunc(ErrTest)
	s.api.DeleteUser(c)

	b, _ := io.ReadAll(w.Body)
	if errStr, isEqual := helpers.CheckExpectedStatusCodeEqualsActual(http.StatusBadRequest, w.Code); !isEqual {
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

func (s *UserControllerTestSuite) Test_DeleteUser_CannotClearSession() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	s.store.SetSaveError(ErrTest)
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)

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
	if errStr, isEqual := helpers.CheckExpectedDataEqualsActual[models.UserView](m, defaultUserView); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *UserControllerTestSuite) Test_GetProfile_NotFound() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)

	s.dbHandler.SetMockGetUserByUsernameFunc(&defaultUser, ErrTest)
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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)

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
	if errStr, isEqual := helpers.CheckExpectedDataEqualsActual[models.User](m, &defaultUser); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *UserControllerTestSuite) Test_GetSelfProfile_NotFound() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)

	s.dbHandler.SetMockGetUserByIDFunc(&defaultUser, ErrTest)
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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)

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
	if errStr, isEqual := helpers.CheckExpectedDataEqualsActual[models.User](m, &defaultUser); !isEqual {
		s.T().Error(errStr)
	}
}

func (s *UserControllerTestSuite) Test_UpdateUser_BadRequest() {
	c, w := helpers.CreateTestContextAndRecorder()
	helpers.AddStoreToContext(c, s.store)
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)

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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)

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
	helpers.AddParamsToContext(c, helpers.UserIDKey, testUserID)

	req, err := helpers.GenerateHttpJSONRequest(http.MethodPatch, defaultUser)
	if err != nil {
		s.T().Error(err)
	}
	c.Request = req

	s.dbHandler.SetMockUpdateUserFunc(&defaultUser, ErrTest)
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

func TestAPIEnv_CreateUser(t *testing.T) {
	type args struct {
		Username     string
		Password     string
		Email        string
		UserDBOutput *models.User
		UserDBError  error
		StoreError   error
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.User]
	}{
		{
			"Create user OK",
			args{
				Username:     testUsername,
				Password:     testPassword,
				Email:        testEmail,
				UserDBOutput: &defaultUser,
				UserDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.User]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedMessage,
				Message:    SuccessfulAccountCreationMsg,
			},
		},
		{
			"Create user missing credentials",
			args{
				UserDBOutput: &defaultUser,
				UserDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.User]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrMissingSignupCredentials,
			},
		},
		{
			"Create user username conflict",
			args{
				Username:     testUsername,
				Password:     testPassword,
				Email:        testEmail,
				UserDBOutput: &defaultUser,
				UserDBError:  gorm.ErrDuplicatedKey,
			},
			helpers.ExpectedJSONOutput[models.User]{
				StatusCode: http.StatusConflict,
				JSONType:   helpers.ExpectedError,
				Error:      ErrUsernameAlreadyExists,
			},
		},
		{
			"Create user DB throws error",
			args{
				Username:     testUsername,
				Password:     testPassword,
				Email:        testEmail,
				UserDBOutput: &defaultUser,
				UserDBError:  ErrTest,
			},
			helpers.ExpectedJSONOutput[models.User]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrTest,
			},
		},
		{
			"Create user DB cannot save session",
			args{
				Username:     testUsername,
				Password:     testPassword,
				Email:        testEmail,
				UserDBOutput: &defaultUser,
				UserDBError:  nil,
				StoreError:   ErrTest,
			},
			helpers.ExpectedJSONOutput[models.User]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrCreateAccountNoCookie,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbTestHandler := &UserDBTestHandler{}
			a := &APIEnv{
				UserDBHandler: dbTestHandler,
			}

			c, w := helpers.CreateTestContextAndRecorder()
			sessionStore := helpers.MakeMockStore()
			helpers.AddStoreToContext(c, sessionStore)

			c.Request = helpers.GenerateHttpFormDataRequest(http.MethodPost, struct {
				Username string
				Password string
				Email    string
			}{tt.args.Username, tt.args.Password, tt.args.Email})
			dbTestHandler.SetMockCreateUserFunc(tt.args.UserDBOutput, tt.args.UserDBError)
			sessionStore.SetSaveError(tt.args.StoreError)
			a.CreateUser(c)

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
