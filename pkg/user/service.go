package user

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/jinzhu/gorm"
	"github.com/yishanzhilu/everest/pkg/models"
)

// Service abstract a yishan user's business logic
type Service interface {
	GetByID(id uint64) (*models.UserModel, error)
	UpdateByID(id uint64, newUser *models.UserModel) (*models.UserModel, error)
	// VerifyUserRefreshToken(id string, token string) (*models.UserModel, error)
	GetOrCreateUserWithGithubOauth(code string) (*models.UserModel, error)
}

type userService struct {
	repo       Repository
	githubRepo OauthRepo
}

// NewUserService create new user service
func NewUserService(repo Repository, githubRepo OauthRepo) Service {
	return &userService{
		repo,
		githubRepo,
	}
}

func (s *userService) GetByID(id uint64) (*models.UserModel, error) {
	return s.repo.ReadByID(id)
}

func (s *userService) UpdateByID(id uint64, newUser *models.UserModel) (*models.UserModel, error) {
	u, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	err = s.repo.UpdateWithStruct(u, newUser)
	if err != nil {
		u = nil
	}
	return u, err
}

func (s *userService) GetOrCreateUserWithGithubOauth(code string) (*models.UserModel, error) {
	t, err := s.githubRepo.GetUserOauthToken(code)
	if err != nil {
		return nil, err
	}
	if t.ErrorMessage != "" {
		return nil, &githubTokenRespError{
			t.ErrorMessage,
			t.ErrorDesc,
			t.ErrorURI,
		}
	}
	gu, err := s.githubRepo.GetUserOauthInfo(t.AccessToken)
	if err != nil {
		return nil, err
	}

	// Does user exist?
	u, err := s.repo.ReadByGithubID(gu.ID)
	// If yes, update token and return
	if err == nil {
		err = s.repo.UpdateWithStruct(u, &models.UserModel{
			GithubToken: t.AccessToken,
		})
		if err != nil {
			u = nil
		}
		return u, err
	}
	// If Not, create user and return
	if gorm.IsRecordNotFoundError(err) {
		return s.CreateUserWithGithubOauth(gu, t)
	}
	// shit happens
	return nil, err
}

func (s *userService) CreateUserWithGithubOauth(gu *GithubUser, t *GithubToken) (*models.UserModel, error) {
	var u models.UserModel
	u.RefreshToken = genereateRefreshToken()
	u.GithubToken = t.AccessToken
	u.GithubID = gu.ID
	u.Name = gu.Name
	if u.Name == "" {
		u.Name = gu.Login
	}
	u.AvatarURL = gu.AvatarURL
	err := s.repo.Create(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func genereateRefreshToken() string {
	refreshTokenByte := make([]byte, 20)
	_, err := rand.Read(refreshTokenByte)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(refreshTokenByte)
}
