package access_token

import (
	"fmt"
	"strings"
	"time"

	"github.com/shyam0507/go-ms-bookstore-oauth/src/utils/crypto_utils"
	"github.com/shyam0507/go-ms-bookstore-utils/rest_errors"
)

const expirationTime = 24
const grantTypePassword = "password"
const grantTypeClientCredentials = "client_credentials"

type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`
	Scope     string `json:"scope"`

	// For scope: password
	Username string `json:"username"`
	Password string `json:"password"`

	// For scope client_credentials
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (at *AccessTokenRequest) Validate() *rest_errors.RestErr {
	switch at.GrantType {
	case grantTypePassword:
		break
	case grantTypeClientCredentials:
		break
	default:
		return rest_errors.NewBadRequestError("invalid grant_type parameter")
	}
	return nil
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
	ClientId    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

func (at *AccessToken) Validate() *rest_errors.RestErr {
	at.AccessToken = strings.TrimSpace(at.AccessToken)

	if at.AccessToken == "" {
		return rest_errors.NewBadRequestError("Invalid Access Token Id")
	}
	if at.ClientId < 0 {
		return rest_errors.NewBadRequestError("Invalid Client Id")
	}

	if at.UserId < 0 {
		return rest_errors.NewBadRequestError("Invalid Client Id")
	}

	if at.Expires < 0 {
		return rest_errors.NewBadRequestError("Invalid Client Id")
	}
	return nil
}

func GetNewAccessToken(userId int64) AccessToken {
	return AccessToken{
		UserId:  userId,
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at AccessToken) IsExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}

func (at *AccessToken) Generate() {
	at.AccessToken = crypto_utils.GetMd5(fmt.Sprintf("at-%d-%d-ran", at.UserId, at.Expires))
}
