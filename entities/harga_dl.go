package entities

import "gorm.io/gorm"

type HargaDL struct {
	gorm.Model
	NominalHarga int `json:"nominal_harga"`
	IsPembelian bool `json:"is_pembelian"`
	IsDL bool `json:"is_dl"`
}