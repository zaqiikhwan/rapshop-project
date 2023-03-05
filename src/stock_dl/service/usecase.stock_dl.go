package service

import (
	"rapsshop-project/entities"
	"rapsshop-project/model"
	"time"
)

type stockDLUsecase struct {
	StockDLRepository model.StockDLRepository
}

func NewSosmedUsecase(repoStockDL model.StockDLRepository) model.StockDLUsecase {
	return &stockDLUsecase{StockDLRepository: repoStockDL}
}

func (sdlu *stockDLUsecase) CreateNewStock(input *model.InputStockDL) error {
	newStock := entities.StockDL{
		StockDL: input.StockDL,
		Waktu: time.Now(),
	}
	if err := sdlu.StockDLRepository.Create(newStock); err != nil {
		return err
	}

	return nil
}

func (sdlu *stockDLUsecase) GetAllStock() ([]entities.StockDL, error) {
	var allStock []entities.StockDL
	allStock, err := sdlu.StockDLRepository.GetAll()
	if err != nil {
		return allStock, err
	}
	return allStock, nil
}

func (sdlu *stockDLUsecase) GetLatestDataStock() (entities.StockDL, error) {
	detailStock, err := sdlu.StockDLRepository.GetLatest()

	if err != nil {
		return detailStock, err
	}
	return detailStock, nil
}

func (sdlu *stockDLUsecase) UpdateStock(input *model.InputStockDL) (entities.StockDL, error) {
	stock, err := sdlu.StockDLRepository.GetLatest()

	if err != nil {
		return stock, err
	}

	updateStockDL := entities.StockDL {
		StockDL: input.StockDL,
	}

	if err := sdlu.StockDLRepository.UpdateByID(updateStockDL, stock.ID); err != nil {
		return stock, err
	}

	stockLatest, err := sdlu.StockDLRepository.GetLatest()

	if err != nil {
		return stock, err
	}
	return stockLatest, nil
} 

func (sdlu *stockDLUsecase) DeleteStock() error {
	stock, err := sdlu.StockDLRepository.GetLatest()

	if err != nil {
		return err
	}

	if err := sdlu.StockDLRepository.DeleteByID(stock.ID); err != nil {
		return err
	}
	return nil
}


