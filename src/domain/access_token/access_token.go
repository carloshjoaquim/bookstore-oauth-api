package access_token

import (
	"github.com/carloshjoaquim/bookstore-users-api/utils/errors"
	"strings"
	"time"
)

const (
	expirationTime = 24
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
	ClientId    int64  `json:"client_id"`
	Expires int64 `json:"expires"`
}

func (at *AccessToken) Validate() *errors.RestErr {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if len(at.AccessToken) == 0 {
		return  errors.NewBadRequestError("invalid access token id")
	}
	if at.UserId <= 0 {
		return  errors.NewBadRequestError("invalid user_id")
	}
	if at.ClientId <= 0 {
		return  errors.NewBadRequestError("invalid client_id")
	}
	if at.Expires <= 0 {
		return  errors.NewBadRequestError("invalid expires")
	}
	return nil
}

func GetNewAccessToken() *AccessToken {
	return &AccessToken{
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at AccessToken) isExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}
