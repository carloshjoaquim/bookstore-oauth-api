package rest

import (
	"errors"
	"github.com/go-resty/resty"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

var (
	client = resty.New()
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestLoginUserTimeoutFromApi(t *testing.T) {
	httpmock.ActivateNonDefault(GetRestClient())
	defer httpmock.DeactivateAndReset()
	repository := usersRepository{}

	// mock to add a new article
	httpmock.RegisterResponder("POST", "https://api.bookstore.com/users/login",
			httpmock.NewErrorResponder(errors.New("timeout")))

	user, err := repository.LoginUser("email@gmail.com", "mypassword")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid restClient response when trying to login user", err.Message)
}

func TestLoginUserInvalidErrorInterface(t *testing.T) {

}

func TestLoginUserInvalidCredentials(t *testing.T) {

}

func TestLoginUserInvalidUserJsonResponse(t *testing.T) {

}

func TestLoginUserNoError(t *testing.T) {

}
