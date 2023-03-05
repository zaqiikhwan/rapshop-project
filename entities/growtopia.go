package entities

import "gorm.io/gorm"

type Growtopia struct {
	gorm.Model
	World string `gorm:"size:255" json:"world"`
	Password string `json:"password"`
	Owner string `json:"owner"`
}