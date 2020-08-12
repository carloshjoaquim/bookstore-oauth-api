package rest

import (
	"encoding/json"
	"github.com/carloshjoaquim/bookstore-oauth-api/src/domain/users"
	"github.com/carloshjoaquim/bookstore-oauth-api/src/utils/errors_utils"
	"github.com/go-resty/resty"
	"net/http"
	"time"
)

var (
	usersRestClient = resty.New().
		SetHostURL("https://api.bookstore.com").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetTimeout(1000 * time.Millisecond)
)

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, *errors_utils.RestErr)
}

type usersRepository struct {}

func GetRestClient() *http.Client {
	return usersRestClient.GetClient()
}
func NewRepository() RestUsersRepository {
	return &usersRepository{}
}

func (u *usersRepository) LoginUser(email string, password string) (*users.User, *errors_utils.RestErr) {
	request := users.UserLoginRequest{
		Email: email,
		Password: password,
	}

	response, err := usersRestClient.R().
		SetBody(request).
		Post("/users/login")

	if err != nil {
		return nil, errors_utils.NewInternalServerError("invalid restClient response when trying to login user")
	}

	if response.StatusCode()  > 299 {
		var restErr errors_utils.RestErr

 		err := json.Unmarshal(response.Body(), &restErr)
		if err != nil {
			return nil, errors_utils.NewInternalServerError("invalid error interface when trying to login user")
		}
		return nil, &restErr
	}

	var user users.User
	if err := json.Unmarshal(response.Body(), &user); err != nil {
		return nil, errors_utils.NewInternalServerError("error when trying to unmarshal users response")
	}

	return &user, nil
}