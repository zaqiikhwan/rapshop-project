package model

import "rapsshop-project/entities"

type InputHargaDL struct {
	NominalHarga int  `json:"nominal_harga"`
	IsPembelian  int `json:"is_pembelian"`
	IsDL         int `json:"is_dl"`
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