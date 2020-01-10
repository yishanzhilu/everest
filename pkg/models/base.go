package models

import "time"

// BaseModel extends gorm.Model, but change ID to sortable uuid
type BaseModel struct {
	ID        uint64     `json:"id" gorm:"type:bigint(20) unsigned auto_increment;not null;primary_key;comment:'主键id'"`
	CreatedAt time.Time  `json:"createdAt" gorm:"comment:'创建时间'"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"comment:'更新时间'"`
	DeletedAt *time.Time `json:"-" gorm:"index;comment:'删除时间'"`
}

// WorkStatus is the enum of status for Goal, Mission and Task
type WorkStatus uint8

const (
	// Doing .
	Doing WorkStatus = iota + 1
	// Todo .
	Todo
	// Done .
	Done
	// Drop .
	Drop
)
