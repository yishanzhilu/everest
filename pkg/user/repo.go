package user

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/yishanzhilu/everest/pkg/common"
)

// Repository abstract user repository layer logic
type Repository interface {
	Create(user *Model) error
	ReadByID(id string) (*Model, error)
	ReadByGithubID(id int64) (*Model, error)
	UpdateStruct(user *Model, new Model) error
}

type mysqlUserRepository struct {
	db *gorm.DB
}

// NewMysqlUserRepository ..
func NewMysqlUserRepository(db *gorm.DB) Repository {
	return &mysqlUserRepository{
		db,
	}
}
func (r *mysqlUserRepository) Create(user *Model) error {
	common.Logger.Debug("Create", user)
	if exist := r.db.NewRecord(user); !exist {
		common.Logger.WithField("exist", exist).Debug("Create NewRecord")
		return r.db.Create(&user).Error
	} // => returns `true` as primary key is blank
	return fmt.Errorf("workprofile with id \"%s\" is already exist", user.ID)
}

func (r *mysqlUserRepository) ReadByID(id string) (*Model, error) {
	common.Logger.Debug("ReadByID", id)
	var user Model
	user.ID = id
	err := r.db.First(&user).Error
	return &user, err
}

func (r *mysqlUserRepository) ReadByGithubID(id int64) (*Model, error) {
	common.Logger.Debug("ReadByGithubID", id)
	var user Model
	user.GithubID = id
	err := r.db.First(&user).Error
	return &user, err
}

// UpdateStruct implementation.
func (r *mysqlUserRepository) UpdateStruct(user *Model, new Model) error {
	common.Logger.Debug("UpdateGithubToken", user.ID)
	err := r.db.Model(user).Updates(new).Error
	return err
}
