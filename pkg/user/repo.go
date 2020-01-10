package user

import (
	"github.com/yishanzhilu/everest/pkg/models"
)

// Repository abstract user repository layer logic
type Repository interface {
	Create(user *models.UserModel) error
	ReadByID(id uint64) (*models.UserModel, error)
	ReadByGithubID(id int64) (*models.UserModel, error)
	UpdateStruct(user *models.UserModel, new models.UserModel) error
}
