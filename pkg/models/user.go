package models

// UserModel defines user info
type UserModel struct {
	BaseModel
	Name         string `json:"name" gorm:"type:varchar(80)"`
	GithubToken  string `json:"-" gorm:"type:char(40);"`
	GithubID     int64  `json:"-" gorm:"unique_index;type:BIGINT"`
	RefreshToken string `json:"-" gorm:"type:char(40);"`
	Minutes      int64  `json:"minutes" gorm:"comment:'用户历程时长'"`
}
