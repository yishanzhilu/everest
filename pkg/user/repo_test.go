package user_test

import (
	"github.com/yishanzhilu/everest/pkg/models"
	. "github.com/yishanzhilu/everest/pkg/user"
)

type stubUserRepository struct {
}

// NewMysqlUserRepository ..
func NewStubUserRepository() Repository {
	return &stubUserRepository{}
}

func (s *stubUserRepository) Create(user *models.UserModel) error {
	return nil
}

func (s *stubUserRepository) ReadByID(id uint64) (*models.UserModel, error) {
	return nil, nil
}

func (s *stubUserRepository) ReadByGithubID(id uint64) (*models.UserModel, error) {

	return nil, nil
}

func (s *stubUserRepository) UpdateWithStruct(user *models.UserModel, newUser *models.UserModel) error {

	return nil
}

/* -------------------------------------------------------------------------- */
/*                               Stub Oauth Repo                              */
/* -------------------------------------------------------------------------- */

type stubOauthRepository struct {
}

// NewMysqlUserRepository ..
func NewStubOauthRepository() OauthRepo {
	return &stubOauthRepository{}
}

func (s *stubOauthRepository) GetUserOauthToken(code string) (*GithubToken, error) {
	return nil, nil
}

func (s stubOauthRepository) GetUserOauthInfo(token string) (*GithubUser, error) {
	return nil, nil
}
