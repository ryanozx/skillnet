package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"

	gcs "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
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

func CheckExpectedDataEqualsActual[T interface{}](m map[string]interface{}, expected T) (errorStr string, isSame bool) {
	var model T
	jsonData, _ := json.Marshal(m["data"])
	if err := json.Unmarshal(jsonData, &model); err != nil {
		return err.Error(), false
	}
	if !cmp.Equal(expected, model, cmp.AllowUnexported(model)) {
		return fmt.Sprintf("Data Error:\nExpected: %v\nActual: %v\n", expected, model), false
	}
	return "", true
}

func CheckExpectedErrorEqualsActual(m map[string]interface{}, errStr string) (errorStr string, isSame bool) {
	if m["error"] != errStr {
		return fmt.Sprintf("Message Error:\nExpected: %s\nActual: %s\n", errStr, m["error"]), false
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
