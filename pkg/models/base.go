package models

import "time"

// BaseModel extends gorm.Model, but change ID to sortable uuid
type BaseModel struct {
	ID        uint64     `json:"id" gorm:"type:bigint(20) unsigned auto_increment;not null;primary_key;comment:'主键id'"`
	CreatedAt time.Time  `json:"createdAt" gorm:"comment:'创建时间'"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"comment:'更新时间'"`
	DeletedAt *time.Time `json:"-" gorm:"index;comment:'删除时间'"`
}

// WorkStatus is the enum of status for Goal and Mission
type WorkStatus uint8

const (
	// StatusAny .
	StatusAny WorkStatus = iota
	// StatusDoing .
	StatusDoing
	// StatusTodo .
	StatusTodo
	// StatusDone .
	StatusDone
	// StatusDrop .
	StatusDrop
)

// WorkStatsMap .
var WorkStatsMap = map[string]WorkStatus{"any": 0, "doing": 1, "todo": 2, "done": 3, "drop": 4}

// WorkStatsMapJSON .
var WorkStatsMapJSON = map[WorkStatus]string{1: "doing", 2: "todo", 3: "done", 4: "drop"}

// WorkStats .
type WorkStats struct {
	Status WorkStatus `gorm:"comment:'0:any, 1:doing, 2:todo, 3:done, 4:drop'"`
}
