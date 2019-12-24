package common

import "time"

// BaseModel base model definition, including fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`, which could be embedded in your models
//    type User struct {
//      common.Model
//    }
// It forks form gorm model.go, but change ID from int to bigint(20)
type BaseModel struct {
	ID        int64 `gorm:"type:bigint(20) unsigned auto_increment;not null;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
