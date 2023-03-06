package entities

import "gorm.io/gorm"

// TODO
// type HargaDL struct {
// 	gorm.Model
// 	NominalHarga int `json:"nominal_harga"`
// 	IsPembelian *int `json:"is_pembelian"`
// 	IsDL *int `json:"is_dl"`
// }


type HargaDL struct {
	gorm.Model
	HargaJualDL int `json:"harga_jual_dl"`
	HargaBeliDL int `json:"harga_beli_dl"`
	HargaJualBGL int `json:"harga_jual_bgl"`
	HargaBeliBGL int `json:"harga_beli_bgl"`
}