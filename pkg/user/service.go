package user

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
)

// Service abstract a yishan user's business logic
type Service interface {
	GetOrCreateUserWithGithubOauth(code string) (*Model, error)
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

func (s userService) GetOrCreateUserWithGithubOauth(code string) (*Model, error) {
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
		err = s.repo.UpdateStruct(u, Model{
			GithubToken:  t.AccessToken,
			RefreshToken: genereateRefreshToken(),
		})
		if err != nil {
			u = nil
		}
		return u, err
	}
	// If Not, create ser
	if gorm.IsRecordNotFoundError(err) {
		return s.CreateUserWithGithubOauth(gu, t)
	}
	return nil, err
}

func (s userService) CreateUserWithGithubOauth(gu *GithubUser, t *GithubToken) (*Model, error) {
	var u Model
	u.ID = xid.New().String()
	u.RefreshToken = genereateRefreshToken()
	u.GithubToken = t.AccessToken
	u.GithubID = gu.ID
	u.Name = gu.Name
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
