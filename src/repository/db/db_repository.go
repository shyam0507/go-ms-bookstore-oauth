package db

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/shyam0507/go-ms-bookstore-oauth/src/clients/cassandra"
	"github.com/shyam0507/go-ms-bookstore-oauth/src/domain/access_token"
	"github.com/shyam0507/go-ms-bookstore-oauth/src/utils/errors"
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
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
	Create(access_token.AccessToken) *errors.RestErr
	UpdateExpirationTime(access_token.AccessToken, string) *errors.RestErr
}

type dbRepository struct {
}

func (r *dbRepository) GetById(id string) (*access_token.AccessToken, *errors.RestErr) {
	var result access_token.AccessToken
	if err := cassandra.GetSession().Query(queuyGetAccessToken, id).Scan(&result.AccessToken, &result.UserId, &result.ClientId, &result.Expires); err != nil {

		if err == gocql.ErrNotFound {
			return nil, errors.NewNotFoundError("Access Token not found for the given id")
		}
		return nil, errors.NewInternalServerError(err.Error())
	}
	return &result, nil
}

func (r *dbRepository) Create(at access_token.AccessToken) *errors.RestErr {
	if err := cassandra.GetSession().Query(queuyCreateAccessToken, at.AccessToken, at.UserId, at.ClientId, at.Expires).Exec(); err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}

func (r *dbRepository) UpdateExpirationTime(at access_token.AccessToken, accessTokenId string) *errors.RestErr {
	if err := cassandra.GetSession().Query(queuyUpdateAccessToken, at.Expires, accessTokenId).Exec(); err != nil {
		fmt.Println(err)
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}
