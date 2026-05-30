package model

import "rapsshop-project/entities"

// type InputHargaDL struct {
// 	NominalHarga int  `json:"nominal_harga"`
// 	IsPembelian  int `json:"is_pembelian"`
// 	IsDL         int `json:"is_dl"`
// }

type InputHargaDL struct {
	HargaJualDL int `json:"harga_jual_dl"`
	HargaBeliDL int `json:"harga_beli_dl"`
	HargaJualBGL int `json:"harga_jual_bgl"`
	HargaBeliBGL int `json:"harga_beli_bgl"`
}

type HargaDLRepository interface {
	Create(newHarga entities.HargaDL) error
	// GetAll() ([]entities.HargaDL, error)
	GetLatest() (entities.HargaDL, error)
	UpdateByID(updateHarga entities.HargaDL, id uint) error
	DeleteByID(id uint) error
}

type HargaDLUsecase interface {
	CreateNewPrice(input *InputHargaDL) error
	// GetAllSosmed() ([]SosmedDto, error) // need paginate implement later
	GetLatestPrice() (entities.HargaDL, error)
	UpdateLatestPrice(input *InputHargaDL) (entities.HargaDL, error)
	DeleteLatestPrice() error
}