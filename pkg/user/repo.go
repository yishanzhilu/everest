package user

import (
	"github.com/yishanzhilu/everest/pkg/models"
)

// Repository abstract user repository layer logic
type Repository interface {
	Create(user *models.UserModel) error
	ReadByID(id uint64) (*models.UserModel, error)
	ReadByGithubID(id uint64) (*models.UserModel, error)
	UpdateWithStruct(user *models.UserModel, newUser *models.UserModel) error
}

// OauthRepo abstract github oauth repo layer logic
type OauthRepo interface {
	GetUserOauthToken(code string) (*GithubToken, error)
	GetUserOauthInfo(token string) (*GithubUser, error)
}
