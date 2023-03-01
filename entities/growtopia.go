package entities

import "gorm.io/gorm"

type Growtopia struct {
	gorm.Model
	World string `gorm:"size:255"`
	Password string
	Owner string
}