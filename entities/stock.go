package entities

import (
	"time"

	"gorm.io/gorm"
)

type StockDL struct {
	gorm.Model
	StockDL int
	Profit int
	Waktu time.Time
}