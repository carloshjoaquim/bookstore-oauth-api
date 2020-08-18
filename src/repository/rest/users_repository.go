package rest

import (
	"encoding/json"
	"github.com/carloshjoaquim/bookstore-oauth-api/src/domain/users"
	"github.com/carloshjoaquim/bookstore-utils-go/rest_errors"
	resty "github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

var (
	usersRestClient = resty.New().
		SetHostURL("http://localhost:8081").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetTimeout(100 * time.Millisecond)
)

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, *rest_errors.RestErr)
}

type usersRepository struct {}

func GetRestClient() *http.Client {
	return usersRestClient.GetClient()
}
func NewUsersRepository() RestUsersRepository {
	return &usersRepository{}
}

func (u *usersRepository) LoginUser(email string, password string) (*users.User, *rest_errors.RestErr) {
	request := users.UserLoginRequest{
		Email: email,
		Password: password,
	}

	response, err := usersRestClient.R().
		SetBody(request).
		Post("/users/login")

	if err != nil {
		return nil, rest_errors.NewInternalServerError(
			"invalid restClient response when trying to login user",
			rest_errors.NewError("rest client error"))
	}

	if response.StatusCode()  > 299 {
		var restErr rest_errors.RestErr

 		err := json.Unmarshal(response.Body(), &restErr)
		if err != nil {
			return nil, rest_errors.NewInternalServerError("invalid error interface when trying to login user",
				rest_errors.NewError("rest client error"))
		}
		return nil, &restErr
	}

	var user users.User
	if err := json.Unmarshal(response.Body(), &user); err != nil {
		return nil, rest_errors.NewInternalServerError("error when trying to unmarshal users response",
			rest_errors.NewError("rest client error"))
	}

	return &user, nil
}