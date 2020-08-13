package access_token

import (
	"github.com/carloshjoaquim/bookstore-oauth-api/src/repository/db"
	"github.com/carloshjoaquim/bookstore-oauth-api/src/repository/rest"
	"github.com/carloshjoaquim/bookstore-oauth-api/src/utils/errors_utils"
	"strings"
)

type Service interface {
	GetById(string) (*AccessToken, *errors_utils.RestErr)
	Create(request AccessTokenRequest) (*AccessToken, *errors_utils.RestErr)
	UpdateExpirationTime(AccessToken) *errors_utils.RestErr
}

type Repository interface {
	GetById(string) (*AccessToken, *errors_utils.RestErr)
	Create(AccessToken) *errors_utils.RestErr
	UpdateExpirationTime(AccessToken) *errors_utils.RestErr
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

func (s *service) GetById(accessTokenId string) (*AccessToken, *errors_utils.RestErr) {
	accessTokenId = strings.TrimSpace(accessTokenId)
	if len(accessTokenId) == 0 {
		return nil, errors_utils.NewBadRequestError("invalid access token id")
	}
	accessToken, err := s.dbRepo.GetById(accessTokenId)
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

func (s *service) Create(request AccessTokenRequest) (*AccessToken, *errors_utils.RestErr) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	user, err := s.restUsersRepo.LoginUser(request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	at :=  GetNewAccessToken(user.Id)
	at.Generate()

	// Save the new access token in Cassandra:
	if err := s.dbRepo.Create(at); err != nil {
		return nil, err
	}
	return &at, nil
}

func (s *service) UpdateExpirationTime(at AccessToken) *errors_utils.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}

	return s.dbRepo.UpdateExpirationTime(at)
}
