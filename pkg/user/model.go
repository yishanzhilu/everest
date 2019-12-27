package user

import (
	"github.com/yishanzhilu/everest/pkg/common"
)

// Model defines user info
type Model struct {
	common.BaseModel
	GithubToken  string `json:"-"`
	GithubID     int64  `json:"-" gorm:"unique_index"`
	Name         string `json:"name"`
	RefreshToken string `json:"-" gorm:"type:char(40);"`
}

// TableName used for gorm to associate with db name
func (Model) TableName() string {
	return "user_models"
}
