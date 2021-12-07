package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100);not null;index"`
	Email    string `gorm:"type:varchar(255);not null;uniqueIndex;"`
	Password string `gorm:"type:varchar(255);"`
}
