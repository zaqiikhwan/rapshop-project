package service

import (
	"rapsshop-project/entities"
	"rapsshop-project/model"
	"time"
)

type stockDLUsecase struct {
	StockDLRepository model.StockDLRepository
}

func NewStockDLUsecase(repoStockDL model.StockDLRepository) model.StockDLUsecase {
	return &stockDLUsecase{StockDLRepository: repoStockDL}
}

func (sdlu *stockDLUsecase) CreateNewStock(input *model.InputStockDL) error {
	newStock := entities.StockDL{
		StockDL: input.StockDL,
		HargaJualDL: input.HargaJualDL,
		HargaBeliDL: input.HargaBeliDL,
		HargaBeliBGL: input.HargaBeliBGL,
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

func (sdlu *stockDLUsecase) UpdateTambahStock(input *model.InputStockDL) (entities.StockDL, error) {
	prices := entities.StockDL{
		Profit:       input.Profit,
		HargaJualDL:  input.HargaJualDL,
		HargaBeliDL:  input.HargaBeliDL,
		HargaBeliBGL: input.HargaBeliBGL,
	}
	return sdlu.StockDLRepository.AdjustLatest(input.StockDL, prices)
}

func (sdlu *stockDLUsecase) UpdateKurangiStock(input *model.InputStockDL) (entities.StockDL, error) {
	prices := entities.StockDL{
		Profit:      input.Profit,
		HargaJualDL: input.HargaJualDL,
		HargaBeliDL: input.HargaBeliDL,
	}
	return sdlu.StockDLRepository.AdjustLatest(-input.StockDL, prices)
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


