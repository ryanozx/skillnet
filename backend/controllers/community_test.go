/*
Contains controllers for Community API.
*/
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
	testCommunityID    = 1
	diffCommunityID    = 10
	invalidCommunityID = "badcommid"
	testCommunityName  = "testCommunity"
	diffCommunityName  = "diffCommunityName"
)

var (
	newCommunity = models.Community{
		Name:  testCommunityName,
		About: null.NewString("Example About Text", true),
	}
	testCommunity = models.Community{
		Model: gorm.Model{
			ID: testCommunityID,
		},
		Name:    testCommunityName,
		OwnerID: testUserID,
		User:    defaultUser,
		About:   null.NewString("Example About Text", true),
	}
	diffCommunity = models.Community{
		Model: gorm.Model{
			ID: diffCommunityID,
		},
		Name:    diffCommunityName,
		OwnerID: testUserID,
		User:    defaultUser,
		About:   null.NewString("Lorem Ipsum", true),
	}
)

type CommunityDBTestHandler struct {
	CreateCommunityFunc    func(*models.Community) (*models.Community, error)
	GetCommunityByIDFunc   func(uint) (*models.Community, error)
	GetCommunityByNameFunc func(string) (*models.Community, error)
	GetCommunitiesFunc     func(*helpers.NullableUint) ([]models.Community, error)
	UpdateCommunityFunc    func(*models.Community, string, string) (*models.Community, error)
}

func (h *CommunityDBTestHandler) CreateCommunity(newCommunity *models.Community) (*models.Community, error) {
	return h.CreateCommunityFunc(newCommunity)
}

func (h *CommunityDBTestHandler) GetCommunityByID(communityID uint) (*models.Community, error) {
	return h.GetCommunityByIDFunc(communityID)
}

func (h *CommunityDBTestHandler) GetCommunityByName(communityName string) (*models.Community, error) {
	return h.GetCommunityByNameFunc(communityName)
}

func (h *CommunityDBTestHandler) GetCommunities(cutoff *helpers.NullableUint) ([]models.Community, error) {
	return h.GetCommunitiesFunc(cutoff)
}

func (h *CommunityDBTestHandler) UpdateCommunity(update *models.Community, communityName string, userID string) (*models.Community, error) {
	return h.UpdateCommunityFunc(update, communityName, userID)
}

func (h *CommunityDBTestHandler) QueryCommunity(queryString string, cutoff int) ([]models.SearchResult, error) {
	return nil, nil
}

func (h *CommunityDBTestHandler) SetMockCreateCommunityFunc(community *models.Community, err error) {
	h.CreateCommunityFunc = func(newCommunity *models.Community) (*models.Community, error) {
		return community, err
	}
}

func (h *CommunityDBTestHandler) SetMockGetCommunityByIDFunc(community *models.Community, err error) {
	h.GetCommunityByIDFunc = func(communityID uint) (*models.Community, error) {
		return community, err
	}
}

func (h *CommunityDBTestHandler) SetMockGetCommunityByNameFunc(community *models.Community, err error) {
	h.GetCommunityByNameFunc = func(communityName string) (*models.Community, error) {
		return community, err
	}
}

func (h *CommunityDBTestHandler) SetMockGetCommunitiesFunc(communities []models.Community, err error) {
	h.GetCommunitiesFunc = func(cutoff *helpers.NullableUint) ([]models.Community, error) {
		return communities, err
	}
}

func (h *CommunityDBTestHandler) SetMockUpdateCommunityFunc(community *models.Community, err error) {
	h.UpdateCommunityFunc = func(update *models.Community, communityName string, userID string) (*models.Community, error) {
		return community, err
	}
}

func TestAPIEnv_InitialiseCommunityHandler(t *testing.T) {
	type fields struct {
		DB *gorm.DB
	}
	tests := []struct {
		name          string
		fields        fields
		expectedEmpty bool
	}{
		{
			"Initialise Community DB OK",
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
			a.InitialiseCommunityHandler()
			if communityDB, ok := a.CommunityDBHandler.(*database.CommunityDB); ok {
				if tt.expectedEmpty && communityDB.DB != nil {
					t.Error("Community DB contains unexpected DB instance")
				} else if !tt.expectedEmpty && communityDB.DB != tt.fields.DB {
					t.Error("CommunityDBHandler not initialised correctly")
				}
			} else {
				t.Error("CommunityDBHandler is nil!")
			}
		})
	}
}

func TestAPIEnv_CreateCommunity(t *testing.T) {
	helpers.SetEnvVars(t)
	type args struct {
		ContextParams     map[string]interface{}
		CommunityData     *models.Community
		CommunityDBOutput *models.Community
		CommunityDBError  error
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.CommunityView]
	}{
		{
			"Create Community OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				CommunityData:     &newCommunity,
				CommunityDBOutput: &testCommunity,
				CommunityDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.CommunityView]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data:       testCommunity.CommunityView(testUserID),
			},
		},
		{
			"Create Community Bad Binding",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
			},
			helpers.ExpectedJSONOutput[models.CommunityView]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrBadBinding,
			},
		},
		{
			"Create Community cannot create",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				CommunityData:    &newCommunity,
				CommunityDBError: ErrTest,
			},
			helpers.ExpectedJSONOutput[models.CommunityView]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrCannotCreateCommunity,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbTestHandler := CommunityDBTestHandler{}
			a := &APIEnv{
				CommunityDBHandler: &dbTestHandler,
			}
			c, w := helpers.CreateTestContextAndRecorder()

			for paramKey, paramVal := range tt.args.ContextParams {
				helpers.AddParamsToContext(c, paramKey, paramVal)
			}

			if tt.args.CommunityData != nil {
				req, err := helpers.GenerateHttpJSONRequest(http.MethodPost, nil)
				if err != nil {
					t.Error(err)
				}

				c.Request = req
			}

			dbTestHandler.SetMockCreateCommunityFunc(tt.args.CommunityDBOutput, tt.args.CommunityDBError)
			a.CreateCommunity(c)

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

func TestAPIEnv_GetCommunities(t *testing.T) {
	helpers.SetEnvVars(t)
	type args struct {
		ContextParams     map[string]interface{}
		QueryParams       map[string]interface{}
		CommunityDBOutput []models.Community
		CommunityDBError  error
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.CommunityArray]
	}{
		{
			"Get Communities OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				CommunityDBOutput: []models.Community{testCommunity},
				CommunityDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.CommunityArray]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data: &models.CommunityArray{
					Communities: []models.CommunityView{*testCommunity.CommunityView(testUserID)},
					NextPageURL: helpers.GenerateCommunitiesNextPageURL(models.BackendAddress, testCommunityID),
				},
			},
		},
		{
			"Get Communities cutoff OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				QueryParams: map[string]interface{}{
					helpers.CutoffKey: testCommunityID,
				},
				CommunityDBOutput: []models.Community{testCommunity},
				CommunityDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.CommunityArray]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data: &models.CommunityArray{
					Communities: []models.CommunityView{*testCommunity.CommunityView(testUserID)},
					NextPageURL: helpers.GenerateCommunitiesNextPageURL(models.BackendAddress, testCommunityID),
				},
			},
		},
		{
			"Get Communities multiple communities OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				CommunityDBOutput: []models.Community{diffCommunity, testCommunity},
				CommunityDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.CommunityArray]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data: &models.CommunityArray{
					Communities: []models.CommunityView{*diffCommunity.CommunityView(testUserID), *testCommunity.CommunityView(testUserID)},
					NextPageURL: helpers.GenerateCommunitiesNextPageURL(models.BackendAddress, testCommunityID),
				},
			},
		},
		{
			"Get Communities bad cutoff",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				QueryParams: map[string]interface{}{
					helpers.CutoffKey: invalidCutoff,
				},
				CommunityDBOutput: []models.Community{testCommunity},
				CommunityDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.CommunityArray]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrBadBinding,
			},
		},
		{
			"Get Communities not found",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey: testUserID,
				},
				CommunityDBOutput: []models.Community{},
				CommunityDBError:  ErrTest,
			},
			helpers.ExpectedJSONOutput[models.CommunityArray]{
				StatusCode: http.StatusNotFound,
				JSONType:   helpers.ExpectedError,
				Error:      ErrCommunityNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbTestHandler := CommunityDBTestHandler{}
			a := &APIEnv{
				CommunityDBHandler: &dbTestHandler,
			}
			c, w := helpers.CreateTestContextAndRecorder()

			for paramKey, paramVal := range tt.args.ContextParams {
				helpers.AddParamsToContext(c, paramKey, paramVal)
			}

			req, err := helpers.GenerateHttpJSONRequest(http.MethodGet, nil)
			if err != nil {
				t.Error(err)
			}

			for paramKey, paramVal := range tt.args.QueryParams {
				helpers.AddParamsToQuery(req, paramKey, paramVal)
			}

			c.Request = req

			dbTestHandler.SetMockGetCommunitiesFunc(tt.args.CommunityDBOutput, tt.args.CommunityDBError)
			a.GetCommunities(c)

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

func TestAPIEnv_GetCommunityByName(t *testing.T) {
	helpers.SetEnvVars(t)
	type args struct {
		ContextParams     map[string]interface{}
		CommunityDBOutput *models.Community
		CommunityDBError  error
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.CommunityView]
	}{
		{
			"Get Community By ID OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:        testUserID,
					helpers.CommunityNameKey: testCommunityName,
				},
				CommunityDBOutput: &testCommunity,
				CommunityDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.CommunityView]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data:       testCommunity.CommunityView(testUserID),
			},
		},
		{
			"Get Community By ID community not found",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:        testUserID,
					helpers.CommunityNameKey: testCommunityName,
				},
				CommunityDBOutput: nil,
				CommunityDBError:  ErrTest,
			},
			helpers.ExpectedJSONOutput[models.CommunityView]{
				StatusCode: http.StatusNotFound,
				JSONType:   helpers.ExpectedError,
				Error:      ErrCommunityNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbTestHandler := CommunityDBTestHandler{}
			a := &APIEnv{
				CommunityDBHandler: &dbTestHandler,
			}
			c, w := helpers.CreateTestContextAndRecorder()

			for paramKey, paramVal := range tt.args.ContextParams {
				helpers.AddParamsToContext(c, paramKey, paramVal)
			}

			req, err := helpers.GenerateHttpJSONRequest(http.MethodGet, nil)
			if err != nil {
				t.Error(err)
			}

			c.Request = req

			dbTestHandler.SetMockGetCommunityByNameFunc(tt.args.CommunityDBOutput, tt.args.CommunityDBError)
			a.GetCommunityByName(c)

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

func TestAPIEnv_UpdateCommunity(t *testing.T) {
	helpers.SetEnvVars(t)
	type args struct {
		ContextParams     map[string]interface{}
		CommunityData     *models.Community
		CommunityDBOutput *models.Community
		CommunityDBError  error
	}
	tests := []struct {
		name     string
		args     args
		expected helpers.ExpectedJSONOutput[models.CommunityView]
	}{
		{
			"Update Communities OK",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:        testUserID,
					helpers.CommunityNameKey: testCommunityName,
				},
				CommunityData:     &newCommunity,
				CommunityDBOutput: &testCommunity,
				CommunityDBError:  nil,
			},
			helpers.ExpectedJSONOutput[models.CommunityView]{
				StatusCode: http.StatusOK,
				JSONType:   helpers.ExpectedData,
				Data:       testCommunity.CommunityView(testUserID),
			},
		},
		{
			"Update Communities bad binding",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:        testUserID,
					helpers.CommunityNameKey: testCommunityName,
				},
			},
			helpers.ExpectedJSONOutput[models.CommunityView]{
				StatusCode: http.StatusBadRequest,
				JSONType:   helpers.ExpectedError,
				Error:      ErrBadBinding,
			},
		},
		{
			"Update Communities not found",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:        testUserID,
					helpers.CommunityNameKey: testCommunityName,
				},
				CommunityData:     &newCommunity,
				CommunityDBOutput: &testCommunity,
				CommunityDBError:  gorm.ErrRecordNotFound,
			},
			helpers.ExpectedJSONOutput[models.CommunityView]{
				StatusCode: http.StatusNotFound,
				JSONType:   helpers.ExpectedError,
				Error:      ErrCommunityNotFound,
			},
		},
		{
			"Update Communities not owner",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:        testUserID,
					helpers.CommunityNameKey: testCommunityName,
				},
				CommunityData:     &newCommunity,
				CommunityDBOutput: &testCommunity,
				CommunityDBError:  helpers.ErrNotOwner,
			},
			helpers.ExpectedJSONOutput[models.CommunityView]{
				StatusCode: http.StatusForbidden,
				JSONType:   helpers.ExpectedError,
				Error:      helpers.ErrNotOwner,
			},
		},
		{
			"Update Communities cannot update",
			args{
				ContextParams: map[string]interface{}{
					helpers.UserIDKey:        testUserID,
					helpers.CommunityNameKey: testCommunityName,
				},
				CommunityData:     &newCommunity,
				CommunityDBOutput: &testCommunity,
				CommunityDBError:  ErrTest,
			},
			helpers.ExpectedJSONOutput[models.CommunityView]{
				StatusCode: http.StatusInternalServerError,
				JSONType:   helpers.ExpectedError,
				Error:      ErrCannotUpdateCommunity,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbTestHandler := CommunityDBTestHandler{}
			a := &APIEnv{
				CommunityDBHandler: &dbTestHandler,
			}
			c, w := helpers.CreateTestContextAndRecorder()

			for paramKey, paramVal := range tt.args.ContextParams {
				helpers.AddParamsToContext(c, paramKey, paramVal)
			}

			if tt.args.CommunityData != nil {
				req, err := helpers.GenerateHttpJSONRequest(http.MethodPost, nil)
				if err != nil {
					t.Error(err)
				}

				c.Request = req
			}

			dbTestHandler.SetMockUpdateCommunityFunc(tt.args.CommunityDBOutput, tt.args.CommunityDBError)
			a.UpdateCommunity(c)

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
