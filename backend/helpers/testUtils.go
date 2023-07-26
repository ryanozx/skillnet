package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"

	gcs "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/models"
)

type MockSessionStore struct {
	Values    map[interface{}]interface{}
	SaveError error
}

func (s *MockSessionStore) ID() string {
	return ""
}

func (s *MockSessionStore) Get(key interface{}) interface{} {
	return s.Values[key]
}

func (s *MockSessionStore) Set(key interface{}, val interface{}) {
	s.Values[key] = val
}

func (s *MockSessionStore) Delete(key interface{}) {
	delete(s.Values, key)
}

func (s *MockSessionStore) Clear() {
	for key := range s.Values {
		s.Delete(key)
	}
}

func (s *MockSessionStore) AddFlash(value interface{}, vars ...string) {

}

func (s *MockSessionStore) Flashes(vars ...string) []interface{} {
	return nil
}

func (s *MockSessionStore) Options(gcs.Options) {

}

func (s *MockSessionStore) Save() error {
	return s.SaveError
}

func (s *MockSessionStore) SetSaveError(err error) {
	s.SaveError = err
}

func (s *MockSessionStore) Reset() {
	s.Clear()
	s.SaveError = nil
}

func MakeMockStore() *MockSessionStore {
	store := MockSessionStore{
		Values: make(map[interface{}]interface{}),
	}
	return &store
}

func AddStoreToContext(ctx *gin.Context, s gcs.Session) {
	ctx.Set(gcs.DefaultKey, s)
}

func GenerateHttpFormDataRequest(method string, s interface{}) *http.Request {
	v := reflect.ValueOf(s)
	typeOfS := v.Type()

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	for i := 0; i < v.NumField(); i++ {
		bodyWriter.WriteField(strings.ToLower(typeOfS.Field(i).Name), v.Field(i).Interface().(string))
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	req, _ := http.NewRequest(http.MethodPost, "", bodyBuf)
	req.Header.Add("Content-Type", contentType)
	return req
}

func GenerateHttpJSONRequest(method string, s interface{}) (*http.Request, error) {
	jsonData, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest(method, "", bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func CreateTestContextAndRecorder() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func ParseJSONString(b []byte) (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	return m, err
}

func AddParamsToQuery(req *http.Request, paramKey string, paramVal interface{}) {
	query := req.URL.Query()
	query.Add(paramKey, fmt.Sprintf("%v", paramVal))
	req.URL.RawQuery = query.Encode()
}

func CheckExpectedJSONEqualsActual[T interface{}](m map[string]interface{}, expected ExpectedJSONOutput[T]) (errorStr string, isSame bool) {
	switch expectedType := expected.JSONType; expectedType {
	case ExpectedData:
		return CheckExpectedDataEqualsActual[T](m, expected.Data)
	case ExpectedError:
		return CheckExpectedErrorEqualsActual(m, expected.Error)
	case ExpectedMessage:
		return CheckExpectedMessageEqualsActual(m, expected.Message)
	default:
		return "Invalid expected type", false
	}
}

type TestFormatter[T interface{}] interface {
	TestFormat() *T
}

type ExpectedJSONOutput[T interface{}] struct {
	StatusCode int
	JSONType   ExpectedJSONType
	Data       TestFormatter[T]
	Error      error
	Message    string
}

type ExpectedJSONType int64

const (
	ExpectedData    ExpectedJSONType = 1
	ExpectedError   ExpectedJSONType = 2
	ExpectedMessage ExpectedJSONType = 3
)

func CheckExpectedDataEqualsActual[T interface{}](m map[string]interface{}, expected TestFormatter[T]) (errorStr string, isSame bool) {
	var model T
	expectedData := expected.TestFormat()
	jsonData, _ := json.Marshal(m["data"])
	if err := json.Unmarshal(jsonData, &model); err != nil {
		return err.Error(), false
	}
	if !reflect.DeepEqual(expectedData, &model) {
		return fmt.Sprintf("Data Error:\nExpected: %v\nActual: %v\n", expectedData, &model), false
	}
	return "", true
}

func CheckExpectedErrorEqualsActual(m map[string]interface{}, err error) (errorStr string, isSame bool) {
	if m["error"] != err.Error() {
		return fmt.Sprintf("Message Error:\nExpected: %s\nActual: %s\n", err.Error(), m["error"]), false
	}
	return "", true
}

func CheckExpectedMessageEqualsActual(m map[string]interface{}, msg string) (errorStr string, isSame bool) {
	if m["message"] != msg {
		return fmt.Sprintf("Message Error:\nExpected: %s\nActual: %s\n", msg, m["message"]), false
	}
	return "", true
}

func CheckExpectedStatusCodeEqualsActual(expected int, actual int) (errorStr string, isSame bool) {
	if expected != actual {
		return fmt.Sprintf("Status Code Error - Expected: %v, Actual: %v\n", expected, actual), false
	}
	return "", true
}

type TestCache struct {
	GetCacheValFunc func(context.Context, uint) (uint64, error)
	SetCacheValFunc func(context.Context, uint) (uint64, error)
}

func (c *TestCache) GetCacheVal(ctx context.Context, postID uint) (uint64, error) {
	return c.GetCacheValFunc(ctx, postID)
}

func (c *TestCache) SetCacheVal(ctx context.Context, postID uint) (uint64, error) {
	return c.SetCacheValFunc(ctx, postID)
}

func (c *TestCache) SetMockGetCacheValFunc(count uint64, err error) {
	c.GetCacheValFunc = func(ctx context.Context, postID uint) (uint64, error) {
		return count, err
	}
}

func (c *TestCache) SetMockSetCacheValFunc(count uint64, err error) {
	c.SetCacheValFunc = func(ctx context.Context, postID uint) (uint64, error) {
		return count, err
	}
}

func (c *TestCache) ResetFuncs() {
	c.GetCacheValFunc = nil
	c.SetCacheValFunc = nil
}

type TestNotificationCreator struct {
	PostNotificationFromEventFunc func(context *gin.Context, notif *models.Notification) error
}

func (nc *TestNotificationCreator) PostNotificationFromEvent(context *gin.Context, notif *models.Notification) error {
	return nc.PostNotificationFromEventFunc(context, notif)
}

func (nc *TestNotificationCreator) SetMockPostNotificationFromEventFunc(err error) {
	nc.PostNotificationFromEventFunc = func(context *gin.Context, notif *models.Notification) error {
		return err
	}
}

func (nc *TestNotificationCreator) ResetFuncs() {
	nc.PostNotificationFromEventFunc = nil
}
