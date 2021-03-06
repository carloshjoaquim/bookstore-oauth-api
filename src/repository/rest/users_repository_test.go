package rest

import (
	"github.com/carloshjoaquim/bookstore-oauth-api/src/domain/users"
	"github.com/carloshjoaquim/bookstore-utils-go/rest_errors"
	"github.com/go-resty/resty/v2"
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

	httpmock.RegisterResponder("POST", "https://api.bookstore.com/users/login",
		httpmock.NewErrorResponder(rest_errors.NewError("timeout")))

	user, err := repository.LoginUser("email@gmail.com", "mypassword")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid restClient response when trying to login user", err.Message)
}

func TestLoginUserInvalidErrorInterface(t *testing.T) {
	httpmock.ActivateNonDefault(GetRestClient())
	defer httpmock.DeactivateAndReset()
	repository := usersRepository{}

	httpmock.RegisterResponder("POST", "https://api.bookstore.com/users/login",
		func(req *http.Request) (*http.Response, error) {

			resp, _ := httpmock.NewJsonResponse(500, `{}`)
			return resp, nil
		},
	)

	user, err := repository.LoginUser("email@gmail.com", "mypassword")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid error interface when trying to login user", err.Message)
}

func TestLoginUserInvalidCredentials(t *testing.T) {
	httpmock.ActivateNonDefault(GetRestClient())
	defer httpmock.DeactivateAndReset()
	repository := usersRepository{}

	httpmock.RegisterResponder("POST", "https://api.bookstore.com/users/login",
		func(req *http.Request) (*http.Response, error) {
			resp ,_ := httpmock.NewJsonResponse(404, rest_errors.RestErr{
				Status: 404,
				Message: "invalid user credentials",
				Error: "not found",
				Causes: nil,
			})
			return resp, nil
		},
	)

	user, err := repository.LoginUser("email@gmail.com", "mypassword")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status)
	assert.EqualValues(t, "invalid user credentials", err.Message)
}

func TestLoginUserInvalidUserJsonResponse(t *testing.T) {
	httpmock.ActivateNonDefault(GetRestClient())
	defer httpmock.DeactivateAndReset()
	repository := usersRepository{}

	httpmock.RegisterResponder("POST", "https://api.bookstore.com/users/login",
		func(req *http.Request) (*http.Response, error) {
			resp ,_ := httpmock.NewJsonResponse(200, `{}`)
			return resp, nil
		},
	)

	user, err := repository.LoginUser("email@gmail.com", "mypassword")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "error when trying to unmarshal users response", err.Message)
}

func TestLoginUserNoError(t *testing.T) {
	httpmock.ActivateNonDefault(GetRestClient())
	defer httpmock.DeactivateAndReset()
	repository := usersRepository{}

	httpmock.RegisterResponder("POST", "https://api.bookstore.com/users/login",
		func(req *http.Request) (*http.Response, error) {
			resp ,_ := httpmock.NewJsonResponse(200, users.User{
				Email: "email@gmail.com",
				FirstName: "Carlos",
				LastName: "Joaquim",
			})
			return resp, nil
		},
	)

	user, err := repository.LoginUser("email@gmail.com", "mypassword")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, "Carlos", user.FirstName)
	assert.EqualValues(t, "Joaquim", user.LastName)
	assert.EqualValues(t, "email@gmail.com", user.Email)
}
