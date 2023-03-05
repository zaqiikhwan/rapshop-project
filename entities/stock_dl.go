package entities

import (
	"time"

	"gorm.io/gorm"
)

type StockDL struct {
	gorm.Model
	StockDL int `json:"stock_dl"`
	Profit int `json:"profit"`
	Waktu time.Time `json:"waktu_terjadi"`
}