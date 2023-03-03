package entities

import "gorm.io/gorm"

type Sosmed struct {
	gorm.Model
	Username string `gorm:"size:255" json:"username"` 
	Platform string `gorm:"size:255" json:"platform"`
	Link string `gorm:"size:255" json:"link"`
}