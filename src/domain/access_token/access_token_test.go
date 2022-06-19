package access_token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccessTokenConstant(t *testing.T) {
	assert.EqualValues(t, 24, expirationTime, "Expiration time should be 24 hours!")

}
func TestGetNewAccessToken(t *testing.T) {
	at := GetNewAccessToken()

	assert.False(t, at.IsExpired(), "brand new access token should not be expired")

}
