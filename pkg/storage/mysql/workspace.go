package mysql

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/yishanzhilu/everest/pkg/common"
	"github.com/yishanzhilu/everest/pkg/workspace"
)

type mysqlWorkprofileRepository struct {
	db *gorm.DB
}

// NewMysqlWorkprofileRepository ..
func NewMysqlWorkprofileRepository(db *gorm.DB) workspace.WorkprofileRepository {
	return &mysqlWorkprofileRepository{
		db,
	}
}
func (r *mysqlWorkprofileRepository) Create(workprofile *workspace.WorkprofileModel) error {
	common.Logger.Debug("Create", workprofile)
	if exist := r.db.NewRecord(workprofile); !exist {
		return r.db.Create(&workprofile).Error
	} // => returns `true` as primary key is blank
	return fmt.Errorf("workprofile with id \"%s\" is already exist", workprofile.ID)
}

func (r *mysqlWorkprofileRepository) FindByID(id string) (*workspace.WorkprofileModel, error) {
	common.Logger.Debug(id)
	workprofile := new(workspace.WorkprofileModel)
	workprofile.ID = id
	err := r.db.First(&workprofile).Error
	return workprofile, err
}

func (r *mysqlWorkprofileRepository) FindList() ([]workspace.WorkprofileModel, error) {
	var workprofiles []workspace.WorkprofileModel
	err := r.db.Find(&workprofiles).Error
	return workprofiles, err
}
