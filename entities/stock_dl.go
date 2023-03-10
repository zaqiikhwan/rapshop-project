package entities

import (
	"time"

	"gorm.io/gorm"
)

type StockDL struct {
	gorm.Model
	StockDL int `json:"stock_dl"`
	HargaJualDL int `json:"harga_jual_dl"`
	HargaBeliDL int `json:"harga_beli_dl"`
	HargaBeliBGL int `json:"harga_beli_bgl"`
	Profit int `json:"profit"`
	Waktu time.Time `json:"waktu_terjadi"`
}