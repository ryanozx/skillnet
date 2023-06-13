package helpers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ryanozx/skillnet/models"
	"github.com/stretchr/testify/suite"
)

type AuthHelperTestSuite struct {
	suite.Suite
	store *MockSessionStore
}

func (s *AuthHelperTestSuite) SetupSuite() {
	s.store = &MockSessionStore{}
}

func (s *AuthHelperTestSuite) TearDownTest() {
	s.store.Clear()
}

func TestUserControllerSuite(t *testing.T) {
	suite.Run(t, new(AuthHelperTestSuite))
}

func (s *AuthHelperTestSuite) Test_ExtractUserCredentials_OK() {
	c, _ := CreateTestContextAndRecorder()

	expected := models.UserCredentials{
		Username: "testUser",
		Password: "12345",
	}

	req := GenerateHttpFormDataRequest(http.MethodPost, expected)
	c.Request = req

	actual := ExtractUserCredentials(c)
	if *actual != expected {
		s.T().Fail()
		fmt.Printf("Actual: %v\n", actual)
		fmt.Printf("Expected: %v\n", expected)
	}
}
