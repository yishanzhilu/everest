package user_test

import (
	"github.com/yishanzhilu/everest/pkg/models"
	. "github.com/yishanzhilu/everest/pkg/user"
)

type stubService struct {
	repo       Repository
	githubRepo OauthRepo
}

func NewStubService() Service {
	repo := NewStubUserRepository()
	githubRepo := NewStubOauthRepository()
	return &stubService{
		repo,
		githubRepo,
	}
}

func (s *stubService) GetByID(id uint64) (*models.UserModel, error) {
	return nil, nil
}

func (s *stubService) UpdateByID(id uint64, newUser *models.UserModel) (*models.UserModel, error) {
	return nil, nil
}

func (s *stubService) GetOrCreateUserWithGithubOauth(code string) (*models.UserModel, error) {
	return nil, nil
}
