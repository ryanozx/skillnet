package controllers

/*
type TestSuiteEnv struct {
	suite.Suite
	api APIEnv
}

func (suite *TestSuiteEnv) SetupSuite() {
	db := database.ConnectTestDatabase()
	api := APIEnv{
		DB: db,
	}
	suite.api = api
}

// clears database
func (suite *TestSuiteEnv) TearDownTest() {
	suite.api.DB.Exec("DELETE FROM post_schemas")
	suite.api.DB.Exec("ALTER SEQUENCE post_schemas_id_seq RESTART WITH 1")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TestSuiteEnv))
}

type HTTPMethod string

const (
	GET    HTTPMethod = http.MethodGet
	POST   HTTPMethod = http.MethodPost
	PATCH  HTTPMethod = http.MethodPatch
	DELETE HTTPMethod = http.MethodDelete
)

func (suite *TestSuiteEnv) Test_GetPosts_EmptyResult() {
	test := setupCrud[[]models.Post](suite)
	test.generateResponse(suite.api.setGetPostsRouter, "/posts")

	test.assertHTTPMethod(GET)
	test.assertHTTPStatus(http.StatusOK)

	expected := []models.Post{}
	actual := test.generateOutput()
	test.assertOutput(expected, actual)
}

type crudTest[T interface{}] struct {
	request   *http.Request
	writer    *httptest.ResponseRecorder
	assertion *assert.Assertions
}

func setupCrud[T interface{}](suite *TestSuiteEnv) *crudTest[T] {
	test := crudTest[T]{
		assertion: suite.Assert(),
	}
	return &test
}

func (t *crudTest[T]) generateResponse(routerSetup func(string) (*http.Request, *httptest.ResponseRecorder, error), url string) {
	request, writer, err := routerSetup(url)
	if err != nil {
		t.assertion.Error(err)
	}
	t.request = request
	t.writer = writer
}

func (t *crudTest[T]) assertHTTPMethod(httpMethod HTTPMethod) {
	t.assertion.Equal(string(httpMethod), t.request.Method, "HTTP method request error")
}

func (t *crudTest[T]) assertHTTPStatus(statusCode int) {
	t.assertion.Equal(statusCode, t.writer.Code, "HTTP request status code error")
}

func (t *crudTest[T]) generateOutput() T {
	body, err := io.ReadAll(t.writer.Body)
	if err != nil {
		t.assertion.Error(err)
	}

	var actual T
	if err := json.Unmarshal(body, &actual); err != nil {
		t.assertion.Error(err)
	}
	return actual
}

func (t *crudTest[T]) assertOutput(expected interface{}, actual interface{}) {
	t.assertion.Equal(expected, actual)
}

type TestRequestHandler struct {
	handler     func(*gin.Context)
	routerRoute string
	method      HTTPMethod
}

func (api *APIEnv) setGetPostsRouter(url string) (*http.Request, *httptest.ResponseRecorder, error) {
	testHandler := TestRequestHandler{
		handler:     api.GetPosts,
		routerRoute: "/posts",
		method:      GET,
	}
	return testHandler.generateReqAndResponse(url, nil)
}

func (th *TestRequestHandler) generateReqAndResponse(url string, body io.Reader) (*http.Request, *httptest.ResponseRecorder, error) {
	router := gin.Default()
	methodString := string(th.method)
	router.Handle(methodString, th.routerRoute, th.handler)
	req, err := http.NewRequest(methodString, url, body)
	w := httptest.NewRecorder()
	if err != nil {
		return req, httptest.NewRecorder(), err
	}
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	return req, w, nil
}

func (suite *TestSuiteEnv) Test_GetPost_OK() {
	test := setupCrud[models.Post](suite)
	post, url := suite.insertTestPost("/posts")

	test.generateResponse(suite.api.setGetPostRouter, url)

	test.assertHTTPMethod(GET)
	test.assertHTTPStatus(http.StatusOK)

	expected := post
	actual := test.generateOutput()
	test.assertOutput(expected, actual)
}

var testUserID = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

// findPath is the prefix of the path to access this post
// e.g. to access the post at "/posts/2", where 2 is the post ID,
// findPath = "/posts"
func (suite *TestSuiteEnv) insertTestPost(findPath string) (models.Post, string) {
	postInput := testPostInput
	postSchema := postInput.ConvertToPostSchema(testUserID)
	result := suite.api.DB.Create(&postSchema)
	post := models.Post{
		PostSchema: *postSchema,
	}
	url := fmt.Sprintf("%s/%v", findPath, post.ID)
	if err := result.Error; err != nil {
		suite.Assertions.Error(err)
		return post, url
	}
	return post, url
}

func (api *APIEnv) setGetPostRouter(url string) (*http.Request, *httptest.ResponseRecorder, error) {
	testHandler := TestRequestHandler{
		handler:     api.GetPostByID,
		routerRoute: "/posts/:id",
		method:      GET,
	}
	return testHandler.generateReqAndResponse(url, nil)
}

var testPostInput = models.PostInput{
	Content: "Hello world!",
}

func (suite *TestSuiteEnv) Test_GetPost_NotFound() {
	test := setupCrud[models.Post](suite)

	test.generateResponse(suite.api.setGetPostRouter, "/posts/1")

	test.assertHTTPMethod(GET)
	test.assertHTTPStatus(http.StatusNotFound)
}

// TODO: incorporate authentication tests for CreatePost
func (suite *TestSuiteEnv) Test_CreatePost_OK() {
	test := setupCrud[models.Post](suite)
	test.generatePostResponse(suite.api.setCreatePostRouter, "/auth/posts", testPostInput)

	test.assertHTTPMethod(POST)
	test.assertHTTPStatus(http.StatusOK)

	expected := testPostInput.ConvertToPostSchema(testUserID)
	actual := test.generateOutput()
	// individual fields have to be compared
	test.assertOutput(expected.Content, actual.Content)
}

func (api *APIEnv) setCreatePostRouter(url string, body *bytes.Buffer) (*http.Request, *httptest.ResponseRecorder, error) {
	testHandler := TestRequestHandler{
		handler:     api.CreatePost,
		routerRoute: "/auth/posts",
		method:      POST,
	}
	return testHandler.generateReqAndResponse(url, body)
}

func (t *crudTest[T]) generatePostResponse(routerSetup func(string, *bytes.Buffer) (*http.Request, *httptest.ResponseRecorder, error), url string, content any) {
	reqBody, err := json.Marshal(content)
	if err != nil {
		t.assertion.Error(err)
	}
	request, writer, err := routerSetup(url, bytes.NewBuffer(reqBody))
	if err != nil {
		t.assertion.Error(err)
	}
	t.request = request
	t.writer = writer
}

func (suite *TestSuiteEnv) Test_UpdatePost_OK() {
	test := setupCrud[models.Post](suite)
	post, url := suite.insertTestPost("/auth/posts")

	test.generatePostResponse(suite.api.setUpdatePostRouter, url, testUpdatePostInput)

	test.assertHTTPMethod(PATCH)
	test.assertHTTPStatus(http.StatusOK)

	expected := post
	expected.Content = testUpdatePostInput.Content
	actual := test.generateOutput()
	test.assertOutput(expected.ID, actual.ID)
	test.assertOutput(expected.Content, actual.Content)
}

func (api *APIEnv) setUpdatePostRouter(url string, body *bytes.Buffer) (*http.Request, *httptest.ResponseRecorder, error) {
	testHandler := TestRequestHandler{
		handler:     api.UpdatePost,
		routerRoute: "/auth/posts/:id",
		method:      PATCH,
	}
	return testHandler.generateReqAndResponse(url, body)
}

var testUpdatePostInput = models.PostInput{
	Content: "this post is updated!",
}

func (suite *TestSuiteEnv) Test_UpdatePost_NotFound() {
	test := setupCrud[models.Post](suite)
	test.generatePostResponse(suite.api.setUpdatePostRouter, "/auth/posts/1", testUpdatePostInput)

	test.assertHTTPMethod(PATCH)
	test.assertHTTPStatus(http.StatusNotFound)
}

func (suite *TestSuiteEnv) Test_DeletePost_OK() {
	test := setupCrud[models.Post](suite)
	_, url := suite.insertTestPost("/auth/posts")

	test.generateResponse(suite.api.setDeletePostRouter, url)

	test.assertHTTPMethod(DELETE)
	test.assertHTTPStatus(http.StatusOK)
}

func (api *APIEnv) setDeletePostRouter(url string) (*http.Request, *httptest.ResponseRecorder, error) {
	testHandler := TestRequestHandler{
		handler:     api.DeletePost,
		routerRoute: "/auth/posts/:id",
		method:      DELETE,
	}
	return testHandler.generateReqAndResponse(url, nil)
}

func (suite *TestSuiteEnv) Test_DeletePost_DeleteTwice() {
	test := setupCrud[models.Post](suite)
	post, url := suite.insertTestPost("/auth/posts")
	suite.deleteTestPost(&post)

	test.generateResponse(suite.api.setDeletePostRouter, url)

	test.assertHTTPMethod(DELETE)
	test.assertHTTPStatus(http.StatusNotFound)
}

func (suite *TestSuiteEnv) deleteTestPost(post *models.Post) error {
	postID := fmt.Sprint(post.ID)
	err := database.DeletePost(suite.api.DB, postID, testUserID)
	if err != nil {
		suite.Assertions.Error(err)
	}
	return err
}

func (suite *TestSuiteEnv) Test_DeletePost_NotFound() {
	test := setupCrud[models.Post](suite)
	test.generateResponse(suite.api.setDeletePostRouter, "/auth/posts/1")

	test.assertHTTPMethod(DELETE)
	test.assertHTTPStatus(http.StatusNotFound)
}
*/
