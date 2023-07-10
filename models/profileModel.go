package models

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	UserID      uint
	Email       string `gorm:"unique"`
	DisplayName string
	UserName    string
	DisplayURL  string
	Bio         string
}
