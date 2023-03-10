package model

import "rapsshop-project/entities"

// must initiate new value with json for data transfer object
type InputStockDL struct {
	StockDL int `json:"stock_dl"`
	Profit  int `json:"profit"`
	HargaJualDL int `json:"harga_jual_dl"`
	HargaBeliDL int `json:"harga_beli_dl"`
	HargaJualBGL int `json:"harga_jual_bgl"`
	HargaBeliBGL int `json:"harga_beli_bgl"`
}

type StockDLRepository interface {
	Create(newStock entities.StockDL) error
	GetAll() ([]entities.StockDL, error)
	GetLatest() (entities.StockDL, error)
	UpdateByID(updateStock entities.StockDL, id uint) error
	DeleteByID(id uint) error
}

type StockDLUsecase interface {
	CreateNewStock(input *InputStockDL) error
	GetAllStock() ([]entities.StockDL, error) // need paginate implement later
	GetLatestDataStock() (entities.StockDL, error)
	UpdateTambahStock(input *InputStockDL) (entities.StockDL, error)
	UpdateKurangiStock(input *InputStockDL) (entities.StockDL, error)
	DeleteStock() error
}