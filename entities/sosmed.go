package entities

import "gorm.io/gorm"

type Sosmed struct {
	gorm.Model
	Username string `gorm:"size:255"`
	Platform string `gorm:"size:255"`
	Link string `gorm:"size:255"`
}