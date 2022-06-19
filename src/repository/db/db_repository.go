package db

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/shyam0507/go-ms-bookstore-oauth/src/clients/cassandra"
	"github.com/shyam0507/go-ms-bookstore-oauth/src/domain/access_token"
	"github.com/shyam0507/go-ms-bookstore-utils/rest_errors"
)

const (
	queuyGetAccessToken    = "SELECT access_token, user_id, client_id, expires From access_tokens where access_token=?;"
	queuyCreateAccessToken = "INSERT into access_tokens (access_token, user_id, client_id, expires) VALUES (?, ?, ?, ?);"
	queuyUpdateAccessToken = "UPDATE access_tokens SET expires =? WHERE  access_token=?;"
)

func New() DbRepository {
	return &dbRepository{}
}

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, *rest_errors.RestErr)
	Create(access_token.AccessToken) *rest_errors.RestErr
	UpdateExpirationTime(access_token.AccessToken, string) *rest_errors.RestErr
}

type dbRepository struct {
}

func (r *dbRepository) GetById(id string) (*access_token.AccessToken, *rest_errors.RestErr) {
	var result access_token.AccessToken
	if err := cassandra.GetSession().Query(queuyGetAccessToken, id).Scan(&result.AccessToken, &result.UserId, &result.ClientId, &result.Expires); err != nil {

		if err == gocql.ErrNotFound {
			return nil, rest_errors.NewNotFoundError("Access Token not found for the given id")
		}
		return nil, rest_errors.NewInternalServerError(err.Error(), err)
	}
	return &result, nil
}

func (r *dbRepository) Create(at access_token.AccessToken) *rest_errors.RestErr {
	if err := cassandra.GetSession().Query(queuyCreateAccessToken, at.AccessToken, at.UserId, at.ClientId, at.Expires).Exec(); err != nil {
		return rest_errors.NewInternalServerError("Error while storing the data to cassandra", err)
	}
	return nil
}

func (r *dbRepository) UpdateExpirationTime(at access_token.AccessToken, accessTokenId string) *rest_errors.RestErr {
	if err := cassandra.GetSession().Query(queuyUpdateAccessToken, at.Expires, accessTokenId).Exec(); err != nil {
		fmt.Println(err)
		return rest_errors.NewInternalServerError("Error while updating the data to cassandra", err)
	}
	return nil
}
