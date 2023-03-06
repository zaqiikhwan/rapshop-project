package entities

import "gorm.io/gorm"

type HargaDL struct {
	gorm.Model
	NominalHarga int `json:"nominal_harga"`
	IsPembelian *int `json:"is_pembelian"`
	IsDL *int `json:"is_dl"`
}