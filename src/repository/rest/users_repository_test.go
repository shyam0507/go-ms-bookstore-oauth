package rest

import (
	"net/http"
	"os"
	"testing"

	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestLoginUserTimeoutFromApi(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@gmail.com","password":"password"}`,
		RespHTTPCode: -1,
		RespBody:     `{}`,
	})

	repository := userRepository{}
	user, err := repository.LoginUser("email@gmail.com", "password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid restclient response when trying to login user", err.Message)
}

func TestLoginUserInvalidErrorInterface(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@gmail.com","password":"password"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message":"invalid login credentials","status":"404","error":"not_found"}`,
	})

	repository := userRepository{}
	user, err := repository.LoginUser("email@gmail.com", "password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid error interface when trying to login user", err.Message)
}

func TestLoginUserInvalidLoginCredentials(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@gmail.com","password":"password"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message":"invalid login credentials","status":404,"error":"not_found"}`,
	})

	repository := userRepository{}
	user, err := repository.LoginUser("email@gmail.com", "password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status)
	assert.EqualValues(t, "invalid login credentials", err.Message)
}

func TestLoginUserInvalidUserJsonResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@gmail.com","password":"password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody: `{
			"id": "34",
			"first_name": "Shyam",
			"last_name": "Nath",
			"email": "shyam114@gmail.com",
			"date_created": "2022-05-14T10:52:18Z",
			"status": "active"
		}`,
	})

	repository := userRepository{}
	user, err := repository.LoginUser("email@gmail.com", "password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "error when trying to unmarshal users response", err.Message)
}

func TestLoginUserNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@gmail.com","password":"password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody: `{
			"id": 34,
			"first_name": "Shyam",
			"last_name": "Nath",
			"email": "shyam114@gmail.com",
			"date_created": "2022-05-14T10:52:18Z",
			"status": "active"
		}`,
	})

	repository := userRepository{}
	user, err := repository.LoginUser("email@gmail.com", "password")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, 34, user.Id)
	assert.EqualValues(t, "Shyam", user.FirstName)
	assert.EqualValues(t, "Nath", user.LastName)
	assert.EqualValues(t, "shyam114@gmail.com", user.Email)
}
