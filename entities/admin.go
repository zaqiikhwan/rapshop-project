package entities

import "time"

type Admin struct {
	ID        string `gorm:"primaryKey"`
	Username  string
	Password  string
	Nama      string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}