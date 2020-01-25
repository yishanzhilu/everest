package models

import (
	"encoding/json"

	"github.com/yishanzhilu/everest/pkg/common"
)

// UserModel defines user info
type UserModel struct {
	BaseModel
	Name         string `json:"name" gorm:"type:varchar(80)"`
	AvatarURL    string `json:"avatarUrl" gorm:"type:text"`
	Minutes      uint64 `json:"minutes"`
	GithubToken  string `json:"-" gorm:"type:char(40);"`
	GithubID     uint64 `json:"-" gorm:"unique_index;type:BIGINT"`
	RefreshToken string `json:"-" gorm:"type:char(40);"`
}

type userModelSerializer struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatarUrl" `
	Minutes   uint64 `json:"minutes"`
	CreatedAt string `json:"createdAt" `
	UpdatedAt string `json:"updatedAt"`
}

func newUserModelSerializer(u *UserModel) *userModelSerializer {
	return &userModelSerializer{
		ID:        u.ID,
		Name:      u.Name,
		AvatarURL: u.AvatarURL,
		Minutes:   u.Minutes,
		CreatedAt: u.CreatedAt.Format(common.TIMESTAMP),
		UpdatedAt: u.UpdatedAt.Format(common.TIMESTAMP),
	}
}

// MarshalJSON implementation.
func (u UserModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(newUserModelSerializer(&u))
}
