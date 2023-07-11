package models

import "time"

type Profile struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      uint `gorm:"unique"`
	DisplayName string
	UserName    string `gorm:"unique"`
	DisplayURL  string
	Bio         string
}
