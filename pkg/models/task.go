package models

import (
	"encoding/json"

	"github.com/yishanzhilu/everest/pkg/common"
)

// TaskModel is user created scalar job.
// If status is todo, it's a todo. Todo is editable. And todo has only content, which means review and minutes
// should be empty value
// IF status is done, it's a record. Record is not editable, since it will send to user's follower's timeline
// for Minutes, one task can only have max minutes 480, which 8 hours
type TaskModel struct {
	BaseModel
	Content   string
	Review    string
	Status    WorkStatus
	Minutes   uint16
	Goal      GoalModel    `gorm:"foreignkey:GoalID;association_autoupdate:false;association_autocreate:false"`
	GoalID    uint64       `gorm:"index"`
	Mission   MissionModel `gorm:"foreignkey:MissionID;association_autoupdate:false;association_autocreate:false"`
	MissionID uint64       `gorm:"index"`
	User      UserModel    `gorm:"foreignkey:UserID;association_autoupdate:false;association_autocreate:false"`
	UserID    uint64       `gorm:"index"`
}

type taskModelSerializer struct {
	ID           uint64 `json:"id"`
	Content      string `json:"content" `
	Review       string `json:"review"`
	Status       string `json:"status" `
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
func (t TaskModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(&taskModelSerializer{
		ID:           t.ID,
		Content:      t.Content,
		Review:       t.Review,
		Status:       WorkStatsMapJSON[t.Status],
		Minutes:      t.Minutes,
		GoalID:       t.GoalID,
		GoalTitle:    t.Goal.Title,
		MissionID:    t.MissionID,
		MissionTitle: t.Mission.Title,
		UserID:       t.UserID,
		CreatedAt:    t.CreatedAt.Format(common.TIMESTAMP),
		UpdatedAt:    t.UpdatedAt.Format(common.TIMESTAMP),
	})
}
