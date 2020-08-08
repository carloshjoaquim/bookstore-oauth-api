package access_token

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAccessTokenConstants(t *testing.T) {
	assert.EqualValues(t, 24, expirationTime, "expiration time should be 24 hours")
}

func TestGetNewAccessToken(t *testing.T) {
	at := GetNewAccessToken()
	assert.False(t, at.isExpired(), "brand new access token should not be expired")
	assert.EqualValues(t, "", at.AccessToken, "new access token should not have defined access token id")
	assert.True(t, at.UserId == 0, "new access token should not have an associated user id")
}

func TestAccessTokenExpired(t *testing.T) {
	at := AccessToken{}
	assert.True(t, at.isExpired(), "empty access token should be expired by default")
}

func TestAccessTokenNotExpired(t *testing.T) {
	at := AccessToken{}
	at.Expires = time.Now().UTC().Add(3 * time.Hour).Unix()
	assert.False(t, at.isExpired(), "access token expiring three hours from now should NOT be expired")
}