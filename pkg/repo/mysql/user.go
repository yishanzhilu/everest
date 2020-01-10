package mysql

import (
	"fmt"

	"github.com/yishanzhilu/everest/pkg/models"

	"github.com/jinzhu/gorm"
	"github.com/yishanzhilu/everest/pkg/common"
	"github.com/yishanzhilu/everest/pkg/user"
)

type mysqlUserRepository struct {
	db *gorm.DB
}

// NewMysqlUserRepository ..
func NewMysqlUserRepository(db *gorm.DB) user.Repository {
	return &mysqlUserRepository{
		db,
	}
}

func (r *mysqlUserRepository) Create(user *models.UserModel) error {
	common.Logger.WithField("user", user).Debug("Create")
	if ok := r.db.NewRecord(user); ok {
		return r.db.Create(&user).Error
	} // => returns `true` as primary key is blank
	return fmt.Errorf("Can't create user with primary key: %d", user.ID)
}

func (r *mysqlUserRepository) ReadByID(id uint64) (*models.UserModel, error) {
	common.Logger.Debug("ReadByID", id)
	var user models.UserModel
	user.ID = id
	err := r.db.First(&user).Error
	return &user, err
}

func (r *mysqlUserRepository) ReadByGithubID(id int64) (*models.UserModel, error) {
	common.Logger.Debug("ReadByGithubID", id)
	var user models.UserModel
	err := r.db.Where("github_id = ?", id).First(&user).Error
	return &user, err
}

// UpdateStruct implementation.
func (r *mysqlUserRepository) UpdateStruct(user *models.UserModel, new models.UserModel) error {
	common.Logger.Debug("UpdateGithubToken", user.ID)
	err := r.db.Model(user).Updates(new).Error
	return err
}
