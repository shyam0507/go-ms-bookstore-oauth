package access_token

import (
	"github.com/shyam0507/go-ms-bookstore-oauth/src/domain/access_token"
	"github.com/shyam0507/go-ms-bookstore-oauth/src/repository/db"
	"github.com/shyam0507/go-ms-bookstore-oauth/src/repository/rest"
	"github.com/shyam0507/go-ms-bookstore-oauth/src/utils/errors"
)

type Repository interface {
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
	Create(access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RestErr)
	UpdateExpirationTime(access_token.AccessToken, string) *errors.RestErr
}

type Service interface {
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
	Create(access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RestErr)
	UpdateExpirationTime(access_token.AccessToken, string) *errors.RestErr
}

type service struct {
	restUsersRepo rest.RestUsersRepository
	dbRepo        db.DbRepository
}

func NewService(usersRepo rest.RestUsersRepository, dbRepo db.DbRepository) Service {
	return &service{
		restUsersRepo: usersRepo,
		dbRepo:        dbRepo,
	}
}
func (s *service) GetById(accessTokenId string) (*access_token.AccessToken, *errors.RestErr) {
	accessToken, err := s.dbRepo.GetById(accessTokenId)

	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (s *service) Create(request access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RestErr) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	user, err := s.restUsersRepo.LoginUser(request.Username, request.Password)

	if err != nil {
		return nil, err
	}

	at := access_token.GetNewAccessToken(user.Id)
	at.Generate()

	// Save the new access token in Cassandra:
	if err := s.dbRepo.Create(at); err != nil {
		return nil, err
	}
	return &at, nil
}

func (s *service) UpdateExpirationTime(at access_token.AccessToken, accessTokenId string) *errors.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}
	err := s.dbRepo.UpdateExpirationTime(at, accessTokenId)

	if err != nil {
		return err
	}
	return nil
}
