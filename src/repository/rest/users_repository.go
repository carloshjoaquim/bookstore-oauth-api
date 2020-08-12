package rest

import (
	"encoding/json"
	"github.com/carloshjoaquim/bookstore-oauth-api/src/domain/users"
	"github.com/carloshjoaquim/bookstore-oauth-api/src/utils/errors"
	"github.com/go-resty/resty"
	"net/http"
	"time"
)

var (
	usersRestClient = resty.New().
		SetHostURL("https://api.bookstore.com").
		SetTimeout(1000 * time.Millisecond)
)

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, *errors.RestErr)
}

type usersRepository struct {}

func GetRestClient() *http.Client {
	return usersRestClient.GetClient()
}
func NewRepository() RestUsersRepository {
	return &usersRepository{}
}

func (u *usersRepository) LoginUser(email string, password string) (*users.User, *errors.RestErr) {
	request := users.UserLoginRequest{
		Email: email,
		Password: password,
	}

	usersRestClient.SetHeader("Accept", "application/json")
	usersRestClient.SetHeaders(map[string]string{
		"Content-Type": "application/json",
	})

	response, err := usersRestClient.R().
		SetBody(request).
		SetResult(&users.User{}).
		Post("/users/login")

	if err != nil {
		return nil, errors.NewInternalServerError("invalid restClient response when trying to login user")
	}

	if response.StatusCode()  > 299 {
		var restErr errors.RestErr
		err := json.Unmarshal(response.Body(), &restErr)
		if err != nil {
			return nil, errors.NewInternalServerError("invalid error interface when trying to login user")
		}
		return nil, &restErr
	}

	var user users.User
	if err := json.Unmarshal(response.Body(), &user); err != nil {
		return nil, errors.NewInternalServerError("error when trying to unmarshal users response")
	}

	return &user, nil
}