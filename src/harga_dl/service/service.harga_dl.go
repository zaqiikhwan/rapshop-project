package service

import (
	"rapsshop-project/entities"
	"rapsshop-project/model"
)

type hargaDLUsecase struct {
	HargaDLRepository model.HargaDLRepository
}

func NewHargaDLUsecase(repoHargaDL model.HargaDLRepository) model.HargaDLUsecase {
	return &hargaDLUsecase{HargaDLRepository: repoHargaDL}
}

func (hdlu *hargaDLUsecase) CreateNewPrice(input *model.InputHargaDL) error {
	newPrice := entities.HargaDL{
		NominalHarga: input.NominalHarga,
		IsPembelian: input.IsPembelian,
		IsDL: input.IsDL,
	}

	if err := hdlu.HargaDLRepository.Create(newPrice); err != nil {
		return err
	}
	return nil
}

func (hdlu *hargaDLUsecase) GetLatestPrice() (entities.HargaDL, error) {
	hargaDL, err := hdlu.HargaDLRepository.GetLatest()
	if err != nil {
		return hargaDL,err
	}
	return hargaDL, nil
}

func (hdlu *hargaDLUsecase) UpdateLatestPrice(input *model.InputHargaDL) (entities.HargaDL, error) {
	var hargaDL entities.HargaDL
	updatePrice := entities.HargaDL{
		NominalHarga: input.NominalHarga,
		IsPembelian: input.IsPembelian,
		IsDL: input.IsDL,
	}

	hargaDL, err := hdlu.HargaDLRepository.GetLatest()
	if err != nil {
		return hargaDL, err
	}

	err = hdlu.HargaDLRepository.UpdateByID(updatePrice, hargaDL.ID)
	if err != nil {
		return hargaDL, err
	}

	updatedPrice, err := hdlu.HargaDLRepository.GetLatest()
	if err != nil {
		return updatedPrice, err
	}
	
	return updatedPrice, nil
}

func (hdlu *hargaDLUsecase) DeleteLatestPrice() error {
	hargaDL, err := hdlu.HargaDLRepository.GetLatest()
	if err != nil {
		return err
	}
	if err := hdlu.HargaDLRepository.DeleteByID(hargaDL.ID); err != nil {
		return err
	}
	return nil
}