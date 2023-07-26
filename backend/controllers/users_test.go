package controllers

import (
	"io"
	"net/http"
	"testing"

	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"gopkg.in/guregu/null.v3"
	"gorm.io/gorm"
)

const (
	testUserID   = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	testUsername = "testuser"
	testPassword = "testPassword123!"
	testEmail    = "abc@def.com"
	diffUserID   = "6ba7b812-9dad-11d1-80b4-00c04fd430c"
)

var (
	defaultUser = models.User{
		ID:              testUserID,
		UserView:        defaultUserView,
		UserCredentials: defaultCreds,
		Email:           testEmail,
	}
	defaultUserView = models.UserView{
		UserMinimal: defaultUserMinimal,
		Title:       null.NewString("tester", true),
		AboutMe:     null.NewString("testing", true),
		ShowTitle:   true,
		ShowAboutMe: false,
	}
	defaultUserMinimal = models.UserMinimal{
		Name: null.NewString("Test User", true),
		URL:  "http://localhost:3000/profile/testuser",
	}
)

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

func TestAPIEnv_InitialiseUserHandler(t *testing.T) {
	type fields struct {
		DB *gorm.DB
	}
	tests := []struct {
		name          string
		fields        fields
		expectedEmpty bool
	}{
		{
			"Initialise User DB OK",
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
			a.InitialiseUserHandler()
			if userDB, ok := a.UserDBHandler.(*database.UserDB); ok {
				if tt.expectedEmpty && userDB.DB != nil {
					t.Error("User DB contains unexpected DB instance")
				} else if !tt.expectedEmpty && userDB.DB != tt.fields.DB {
					t.Error("UserDBHandler not initialised correctly")
				}
			} else {
				t.Error("AuthDBHandler is nil!")
			}
		})
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

func TestAPIEnv_DeleteUser(t *testing.T) {
	type args struct {
		StoreParams map[string]interface{}
		UserDBError error
		StoreError  error
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.User]
	}{
		{
			"Delete User OK",
			args{
				StoreParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				UserDBError: nil,
				StoreError:  nil,
			},
			helpers.ExpectedJSONOutput[models.User]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedMessage,
				Message:    SuccessfulAccountDeleteMsg,
			},
		},
		{
			"Delete User not found",
			args{
				StoreParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				UserDBError: gorm.ErrRecordNotFound,
			},
			helpers.ExpectedJSONOutput[models.User]{
				StatusCode: http.StatusNotFound,
				JSONType:   helpers.ExpectedError,
				Error:      ErrUserNotFound,
			},
		},
		{
			"Delete User cannot delete",
			args{
				StoreParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				UserDBError: ErrTest,
			},
			helpers.ExpectedJSONOutput[models.User]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrCannotDeleteUser,
			},
		},
		{
			"Delete User cannot clear session",
			args{
				StoreParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				UserDBError: nil,
				StoreError:  ErrTest,
			},
			helpers.ExpectedJSONOutput[models.User]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrSessionClearFailed,
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
			store := helpers.MakeMockStore()
			helpers.AddStoreToContext(c, store)

			for paramKey, paramVal := range tt.args.StoreParams {
				store.Set(paramKey, paramVal)
			}

			store.SetSaveError(tt.args.StoreError)

			dbTestHandler.SetMockDeleteUserFunc(tt.args.UserDBError)
			a.DeleteUser(c)

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

func TestAPIEnv_GetProfile(t *testing.T) {
	type args struct {
		ContextParams map[string]interface{}
		UserDBOutput  *models.User
		UserDBError   error
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.UserView]
	}{
		{
			"Get Profile OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:   testUserID,
					helpers.UsernameKey: testUsername,
				},
				UserDBOutput: &defaultUser,
				UserDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.UserView]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data:       defaultUser.GetUserView(testUserID, testClientAddress),
			},
		},
		{
			"Get other Profile OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:   diffUserID,
					helpers.UsernameKey: testUsername,
				},
				UserDBOutput: &defaultUser,
				UserDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.UserView]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data:       defaultUser.GetUserView(diffUserID, testClientAddress),
			},
		},
		{
			"Get Profile not found",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:   testUserID,
					helpers.UsernameKey: testUsername,
				},
				UserDBOutput: &defaultUser,
				UserDBError:  ErrTest,
			},
			helpers.ExpectedJSONOutput[models.UserView]{
				StatusCode: http.StatusNotFound,
				JSONType:   helpers.ExpectedError,
				Error:      ErrUserNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbTestHandler := &UserDBTestHandler{}
			a := &APIEnv{
				UserDBHandler: dbTestHandler,
				ClientAddress: testClientAddress,
			}

			c, w := helpers.CreateTestContextAndRecorder()
			for paramKey, paramVal := range tt.args.ContextParams {
				helpers.AddParamsToContext(c, paramKey, paramVal)
			}

			dbTestHandler.SetMockGetUserByUsernameFunc(tt.args.UserDBOutput, tt.args.UserDBError)
			a.GetProfile(c)

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

func TestAPIEnv_GetSelfProfile(t *testing.T) {
	type args struct {
		ContextParams map[string]interface{}
		UserDBOutput  *models.User
		UserDBError   error
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.User]
	}{
		{
			"Get own profile OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				UserDBOutput: &defaultUser,
				UserDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.User]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data:       &defaultUser,
			},
		},
		{
			"Get own profile not found",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				UserDBOutput: &defaultUser,
				UserDBError:  ErrTest,
			},
			helpers.ExpectedJSONOutput[models.User]{
				StatusCode: http.StatusNotFound,
				JSONType:   helpers.ExpectedError,
				Error:      ErrUserNotFound,
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

			for paramKey, paramVal := range tt.args.ContextParams {
				helpers.AddParamsToContext(c, paramKey, paramVal)
			}

			dbTestHandler.SetMockGetUserByIDFunc(tt.args.UserDBOutput, tt.args.UserDBError)
			a.GetSelfProfile(c)

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

func TestAPIEnv_UpdateUser(t *testing.T) {
	type args struct {
		ContextParams map[string]interface{}
		UserUpdate    *models.User
		UserDBOutput  *models.User
		UserDBError   error
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.User]
	}{
		{
			"Update User OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				UserUpdate:   &defaultUser,
				UserDBOutput: &defaultUser,
				UserDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.User]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data:       &defaultUser,
			},
		},
		{
			"Update User bad request",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
			},
			helpers.ExpectedJSONOutput[models.User]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrBadBinding,
			},
		},
		{
			"Update User not found",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				UserUpdate:   &defaultUser,
				UserDBOutput: &defaultUser,
				UserDBError:  gorm.ErrRecordNotFound,
			},
			helpers.ExpectedJSONOutput[models.User]{
				StatusCode: http.StatusNotFound,
				JSONType:   helpers.ExpectedError,
				Error:      ErrUserNotFound,
			},
		},
		{
			"Update User cannot update",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				UserUpdate:   &defaultUser,
				UserDBOutput: &defaultUser,
				UserDBError:  ErrTest,
			},
			helpers.ExpectedJSONOutput[models.User]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrCannotUpdateUser,
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
			for paramKey, paramVal := range tt.args.ContextParams {
				helpers.AddParamsToContext(c, paramKey, paramVal)
			}

			if tt.args.UserUpdate != nil {
				req, err := helpers.GenerateHttpJSONRequest(http.MethodPatch, tt.args.UserUpdate)
				if err != nil {
					t.Error(err)
				}
				c.Request = req
			}

			dbTestHandler.SetMockUpdateUserFunc(tt.args.UserDBOutput, tt.args.UserDBError)
			a.UpdateUser(c)

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
