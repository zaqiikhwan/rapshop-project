package model

import "rapsshop-project/entities"

type InputStockDL struct {
	StockDL int
	Profit  int
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
	UpdateStock(input *InputStockDL) (entities.StockDL, error)
	DeleteStock() error
}