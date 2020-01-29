package models

import (
	"encoding/json"

	"github.com/yishanzhilu/everest/pkg/common"
)

// MissionModel defines mission
type MissionModel struct {
	BaseModel
	User        UserModel  `gorm:"foreignkey:UserID;save_associations:false"`
	UserID      uint64     `gorm:"index"`
	Goal        GoalModel  `gorm:"foreignkey:GoalID;save_associations:false"`
	GoalID      uint64     `gorm:"index"`
	Title       string     `gorm:"not null;size:80"`
	Description string     `gorm:"not null;size:225"`
	Status      WorkStatus `gorm:"index"`
	Minutes     uint64
}

type missionModelSerializer struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title" `
	Description string `json:"description"`
	Status      string `json:"status" `
	Minutes     uint64 `json:"minutes"`
	GoalID      uint64 `json:"goalID,omitempty"`
	GoalTitle   string `json:"goalTitle,omitempty"`
	UserID      uint64 `json:"ownerID,omitempty"`
	CreatedAt   string `json:"createdAt" `
	UpdatedAt   string `json:"updatedAt"`
}

// MarshalJSON .
func (m MissionModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(&missionModelSerializer{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Status:      WorkStatsMapJSON[m.Status],
		Minutes:     m.Minutes,
		GoalID:      m.Goal.ID,
		GoalTitle:   m.Goal.Title,
		UserID:      m.UserID,
		CreatedAt:   m.CreatedAt.UTC().Format(common.TIMESTAMP),
		UpdatedAt:   m.UpdatedAt.UTC().Format(common.TIMESTAMP),
	})
}
