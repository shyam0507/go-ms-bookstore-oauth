package rest

import (
	"encoding/json"
	"time"

	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/shyam0507/go-ms-bookstore-oauth/src/domain/users"
	"github.com/shyam0507/go-ms-bookstore-utils/rest_errors"
)

var userRestClient = rest.RequestBuilder{
	BaseURL: "http://localhost:8082",
	Timeout: 100 * time.Millisecond,
}

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, *rest_errors.RestErr)
}

type userRepository struct {
}

func NewRepository() RestUsersRepository {
	return &userRepository{}
}

func (u *userRepository) LoginUser(email string, password string) (*users.User, *rest_errors.RestErr) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}

	response := userRestClient.Post("/users/login", request)

	if response == nil || response.Response == nil {
		return nil, rest_errors.NewInternalServerError("invalid restclient response when trying to login user", response.Err)
	}

	if response.StatusCode > 299 {
		var restErr rest_errors.RestErr
		err := json.Unmarshal(response.Bytes(), &restErr)
		if err != nil {
			return nil, rest_errors.NewInternalServerError("invalid error interface when trying to login user", response.Err)
		}
		return nil, &restErr
	}
	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, rest_errors.NewInternalServerError("error when trying to unmarshal users response", response.Err)
	}
	return &user, nil
}
