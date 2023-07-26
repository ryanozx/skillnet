package controllers

import (
	"io"
	"net/http"
	"testing"

	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	defaultCreds = models.UserCredentials{
		Username: testUsername,
		Password: testPassword,
	}
	defaultLoginUserDBEntry = models.User{
		UserCredentials: defaultCreds,
	}
	emptyUserCreds = models.UserCredentials{
		Username: "",
		Password: "",
	}
	incorrectPasswordCreds = models.UserCredentials{
		Username: defaultCreds.Username,
		Password: "1234",
	}
)

func TestAPIEnv_InitialiseAuthHandler(t *testing.T) {
	type fields struct {
		DB *gorm.DB
	}
	tests := []struct {
		name          string
		fields        fields
		expectedEmpty bool
	}{
		{
			"Initialise Auth DB OK",
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
			a.InitialiseAuthHandler()
			if authDB, ok := a.AuthDBHandler.(*database.UserDB); ok {
				if tt.expectedEmpty && authDB.DB != nil {
					t.Error("User DB contains unexpected DB instance")
				} else if !tt.expectedEmpty && authDB.DB != tt.fields.DB {
					t.Error("AuthDBHandler not initialised correctly")
				}
			} else {
				t.Error("AuthDBHandler is nil!")
			}
		})
	}
}

func TestAPIEnv_GetLogin(t *testing.T) {
	type args struct {
		StoreParams map[string]interface{}
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.UserCredentials]
	}{
		{
			"Get Login OK",
			args{},
			helpers.ExpectedJSONOutput[models.UserCredentials]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedMessage,
				Message:    GetLoginOKMsg,
			},
		},
		{
			"Get Login already logged in",
			args{
				StoreParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
			},
			helpers.ExpectedJSONOutput[models.UserCredentials]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrAlreadyLoggedIn,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbTestHandler := UserDBTestHandler{}
			a := &APIEnv{
				AuthDBHandler: &dbTestHandler,
			}
			c, w := helpers.CreateTestContextAndRecorder()
			store := helpers.MakeMockStore()
			helpers.AddStoreToContext(c, store)

			for paramKey, paramVal := range tt.args.StoreParams {
				store.Set(paramKey, paramVal)
			}

			a.GetLogin(c)

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

func TestAPIEnv_PostLogin(t *testing.T) {
	type args struct {
		StoreParams  map[string]interface{}
		UserCreds    *models.UserCredentials
		UserDBOutput *models.User
		UserDBError  error
		SaveError    error
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.UserCredentials]
	}{
		{
			"Post Login OK",
			args{
				UserCreds:    &defaultCreds,
				UserDBOutput: &defaultLoginUserDBEntry,
				UserDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.UserCredentials]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedMessage,
				Message:    LoginSuccessfulMsg,
			},
		},
		{
			"Post Login already logged in",
			args{
				StoreParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				UserCreds:    &defaultCreds,
				UserDBOutput: &defaultLoginUserDBEntry,
				UserDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.UserCredentials]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrAlreadyLoggedIn,
			},
		},
		{
			"Post Login empty user password",
			args{
				UserCreds:    &emptyUserCreds,
				UserDBOutput: &defaultLoginUserDBEntry,
				UserDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.UserCredentials]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrMissingUserCredentials,
			},
		},
		{
			"Post Login cannot retrieve user",
			args{
				UserCreds:    &defaultCreds,
				UserDBOutput: &defaultLoginUserDBEntry,
				UserDBError:  ErrTest,
			},
			helpers.ExpectedJSONOutput[models.UserCredentials]{
				StatusCode: http.StatusUnauthorized,
				JSONType:   helpers.ExpectedError,
				Error:      ErrIncorrectUserCredentials,
			},
		},
		{
			"Post Login incorrect password",
			args{
				UserCreds:    &incorrectPasswordCreds,
				UserDBOutput: &defaultLoginUserDBEntry,
				UserDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.UserCredentials]{
				StatusCode: http.StatusUnauthorized,
				JSONType:   helpers.ExpectedError,
				Error:      ErrIncorrectUserCredentials,
			},
		},
		{
			"Post Login cannot save",
			args{
				UserCreds:    &defaultCreds,
				UserDBOutput: &defaultLoginUserDBEntry,
				UserDBError:  nil,
				SaveError:    ErrTest,
			},
			helpers.ExpectedJSONOutput[models.UserCredentials]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrCookieSaveFail,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbTestHandler := UserDBTestHandler{}
			a := &APIEnv{
				AuthDBHandler: &dbTestHandler,
			}
			c, w := helpers.CreateTestContextAndRecorder()
			store := helpers.MakeMockStore()
			helpers.AddStoreToContext(c, store)

			for paramKey, paramVal := range tt.args.StoreParams {
				store.Set(paramKey, paramVal)
			}

			store.SetSaveError(tt.args.SaveError)

			expectedUser := *tt.args.UserDBOutput
			hash, err := bcrypt.GenerateFromPassword([]byte(tt.args.UserDBOutput.Password), bcrypt.DefaultCost)
			if err != nil {
				t.Error(err)
			}
			expectedUser.Password = string(hash)

			c.Request = helpers.GenerateHttpFormDataRequest(http.MethodPost, *tt.args.UserCreds)
			dbTestHandler.SetMockGetUserByUsernameFunc(&expectedUser, tt.args.UserDBError)
			a.PostLogin(c)

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

func TestAPIEnv_PostLogout(t *testing.T) {
	type args struct {
		StoreParams map[string]interface{}
		SaveError   error
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.UserCredentials]
	}{
		{
			"Post Logout OK",
			args{
				StoreParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
			},
			helpers.ExpectedJSONOutput[models.UserCredentials]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedMessage,
				Message:    SuccessfulLogoutMsg,
			},
		},
		{
			"Post Logout Invalid Session",
			args{},
			helpers.ExpectedJSONOutput[models.UserCredentials]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrNoValidSession,
			},
		},
		{
			"Post Logout OK",
			args{
				StoreParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				SaveError: ErrTest,
			},
			helpers.ExpectedJSONOutput[models.UserCredentials]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrSessionClearFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbTestHandler := UserDBTestHandler{}
			a := &APIEnv{
				AuthDBHandler: &dbTestHandler,
			}
			c, w := helpers.CreateTestContextAndRecorder()
			store := helpers.MakeMockStore()
			helpers.AddStoreToContext(c, store)

			for paramKey, paramVal := range tt.args.StoreParams {
				store.Set(paramKey, paramVal)
			}

			store.SetSaveError(tt.args.SaveError)

			a.PostLogout(c)

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
