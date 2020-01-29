package models

import (
	"encoding/json"

	"github.com/yishanzhilu/everest/pkg/common"
)

// TodoModel .
type TodoModel struct {
	BaseModel
	WorkStats
	Content   string
	Goal      GoalModel    `gorm:"foreignkey:GoalID;save_associations:false"`
	GoalID    uint64       `gorm:"index"`
	Mission   MissionModel `gorm:"foreignkey:MissionID;save_associations:false"`
	MissionID uint64       `gorm:"index"`
	User      UserModel    `gorm:"foreignkey:UserID;save_associations:false"`
	UserID    uint64       `gorm:"index"`
}

type todoModelSerializer struct {
	ID           uint64 `json:"id"`
	Content      string `json:"content"`
	Status       string `json:"status"`
	GoalID       uint64 `json:"goalID,omitempty"`
	GoalTitle    string `json:"goalTitle,omitempty"`
	MissionID    uint64 `json:"missionID,omitempty"`
	MissionTitle string `json:"missionTitle,omitempty"`
	UserID       uint64 `json:"ownerID,omitempty"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

// MarshalJSON .
func (t TodoModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(&todoModelSerializer{
		ID:           t.ID,
		Content:      t.Content,
		Status:       WorkStatsMapJSON[t.Status],
		GoalID:       t.GoalID,
		GoalTitle:    t.Goal.Title,
		MissionID:    t.MissionID,
		MissionTitle: t.Mission.Title,
		UserID:       t.UserID,
		CreatedAt:    t.CreatedAt.UTC().Format(common.TIMESTAMP),
		UpdatedAt:    t.UpdatedAt.UTC().Format(common.TIMESTAMP),
	})
}
