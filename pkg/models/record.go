package models

import (
	"encoding/json"

	"github.com/yishanzhilu/everest/pkg/common"
)

// RecordModel is for user finished job.
// Record is not editable, since it will send to user's follower's timeline in next vision.
// For Minutes, one record can only have max minutes 480, which 8 hours
type RecordModel struct {
	BaseModel
	Content   string
	Review    string
	Minutes   uint16
	Goal      GoalModel    `gorm:"foreignkey:GoalID;association_autoupdate:false;association_autocreate:false"`
	GoalID    uint64       `gorm:"index"`
	Mission   MissionModel `gorm:"foreignkey:MissionID;association_autoupdate:false;association_autocreate:false"`
	MissionID uint64       `gorm:"index"`
	User      UserModel    `gorm:"foreignkey:UserID;association_autoupdate:false;association_autocreate:false"`
	UserID    uint64       `gorm:"index"`
}

type recordModelSerializer struct {
	ID           uint64 `json:"id"`
	Content      string `json:"content" `
	Review       string `json:"review"`
	Minutes      uint16 `json:"minutes"`
	GoalID       uint64 `json:"goalID,omitempty"`
	GoalTitle    string `json:"goalTitle,omitempty"`
	MissionID    uint64 `json:"missionID,omitempty"`
	MissionTitle string `json:"missionTitle,omitempty"`
	UserID       uint64 `json:"ownerID,omitempty"`
	CreatedAt    string `json:"createdAt" `
	UpdatedAt    string `json:"updatedAt"`
}

// MarshalJSON .
func (t RecordModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(&recordModelSerializer{
		ID:           t.ID,
		Content:      t.Content,
		Review:       t.Review,
		Minutes:      t.Minutes,
		GoalID:       t.GoalID,
		GoalTitle:    t.Goal.Title,
		MissionID:    t.MissionID,
		MissionTitle: t.Mission.Title,
		UserID:       t.UserID,
		CreatedAt:    t.CreatedAt.UTC().Format(common.TIMESTAMP),
		UpdatedAt:    t.UpdatedAt.UTC().Format(common.TIMESTAMP),
	})
}