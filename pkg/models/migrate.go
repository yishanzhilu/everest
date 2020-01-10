package models

import "github.com/jinzhu/gorm"

// Migrate .
func Migrate(db *gorm.DB) {
	db.AutoMigrate(&UserModel{})
}
