package access_token

import (
	"fmt"
	"github.com/carloshjoaquim/bookstore-oauth-api/src/utils/crypto_utils"
	"github.com/carloshjoaquim/bookstore-oauth-api/src/utils/errors_utils"
	"strings"
	"time"
)

const (
	expirationTime             = 24
	grantTypePassword          = "password"
	grantTypeClientCredentials = "client_credentials"
)

type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`
	Scope     string `json:"scope"`

	// Used for password grand type
	Username string `json:"username"`
	Password string `json:"password"`

	// Used for client credentials grant type
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
	ClientId    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}


func (at *AccessTokenRequest) Validate() *errors_utils.RestErr {
	switch at.GrantType {
	case grantTypePassword: {
		if at.Password == "" || at.Username == "" {
			return errors_utils.NewBadRequestError("username or password cannot be empty for grant_type = password")
		}
	}
	case grantTypeClientCredentials: {
		if at.ClientId == "" || at.ClientSecret == "" {
			return errors_utils.NewBadRequestError("client_id or client_secret cannot be empty for grant_type = client_credentials")
		}
	}
	default:
		return errors_utils.NewBadRequestError("invalid grant_type.")
	}
	return nil
}

func (at *AccessToken) Validate() *errors_utils.RestErr {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if len(at.AccessToken) == 0 {
		return errors_utils.NewBadRequestError("invalid access token id")
	}
	if at.UserId <= 0 {
		return errors_utils.NewBadRequestError("invalid user_id")
	}
	if at.ClientId <= 0 {
		return errors_utils.NewBadRequestError("invalid client_id")
	}
	if at.Expires <= 0 {
		return errors_utils.NewBadRequestError("invalid expires")
	}
	return nil
}

func GetNewAccessToken(userId int64) AccessToken {
	return AccessToken{
		UserId:  userId,
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at AccessToken) isExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}

func (at *AccessToken) Generate() {
	at.AccessToken = crypto_utils.GetMd5(fmt.Sprintf("at-%d-%d-ran", at.UserId, at.Expires))
}