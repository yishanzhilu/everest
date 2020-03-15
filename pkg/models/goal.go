package models

import (
	"encoding/json"

	"github.com/yishanzhilu/everest/pkg/common"
)

// GoalModel defines user info
type GoalModel struct {
	BaseModel
	Title       string     `gorm:"not null;size:80"`
	Description string     `gorm:"not null;size:255"`
	Status      WorkStatus `gorm:"index:idx_status;"`
	Minutes     uint64
	Missions    []MissionModel `gorm:"foreignkey:GoalID"`
	User        UserModel      `gorm:"foreignkey:UserID;save_associations:false"`
	UserID      uint64         `gorm:"index"`
}

type goalModelSerializer struct {
	ID          uint64         `json:"id"`
	Title       string         `json:"title" `
	Description string         `json:"description"`
	Status      string         `json:"status" `
	Minutes     uint64         `json:"minutes"`
	Missions    []MissionModel `json:"missions,omitempty"`
	UserID      uint64         `json:"ownerID,omitempty"`
	CreatedAt   string         `json:"createdAt" `
	UpdatedAt   string         `json:"updatedAt"`
}

// MarshalJSON .
func (g GoalModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(&goalModelSerializer{
		ID:          g.ID,
		Title:       g.Title,
		Description: g.Description,
		Status:      WorkStatsMapJSON[g.Status],
		Minutes:     g.Minutes,
		Missions:    g.Missions,
		UserID:      g.UserID,
		CreatedAt:   g.CreatedAt.UTC().Format(common.TIMESTAMP),
		UpdatedAt:   g.UpdatedAt.UTC().Format(common.TIMESTAMP),
	})
}
