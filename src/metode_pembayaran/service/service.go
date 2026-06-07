package service

import (
	"os"
	"rapsshop-project/entities"
	"rapsshop-project/model"
)

type metodePembayaranUsecase struct {
	RepoMetodePembayaran model.MetodePembayaranRepository
}

func NewMetodePembayaranUsecase(repoMetodePembayaran model.MetodePembayaranRepository) model.MetodePembayaranUsecase {
	return &metodePembayaranUsecase{RepoMetodePembayaran: repoMetodePembayaran}
}

func (mpu *metodePembayaranUsecase) CreateNewPembayaran(input *model.InputMetodePembayaran) error {
	newPembayaran := entities.MetodePembayaran{
		IndexPembayaran:      &input.IndexPembayaran,
		JenisPembayaran:      input.JenisPembayaran,
		KredensialPembayaran: input.KredensialPembayaran,
		Pemilik:              input.Pemilik,
	}

	if err := mpu.RepoMetodePembayaran.Create(newPembayaran); err != nil {
		return err
	}

	return nil
}

func (mpu *metodePembayaranUsecase) GetAllPembayaran() ([]entities.MetodePembayaran, error) {
	allPembayaran, err := mpu.RepoMetodePembayaran.GetAll()

	if err != nil {
		return allPembayaran, err
	}

	return allPembayaran, nil
}

func (mpu *metodePembayaranUsecase) GetCheckoutOptions() (model.CheckoutPaymentOptionsResponse, error) {
	manualMethods, err := mpu.RepoMetodePembayaran.GetAll()
	if err != nil {
		return model.CheckoutPaymentOptionsResponse{}, err
	}

	gatewayEnabled := os.Getenv("AUTHORIZATION_VALUE") != "" && os.Getenv("MIDTRANS") != ""
	return model.NewCheckoutPaymentOptions(manualMethods, gatewayEnabled), nil
}

func (mpu *metodePembayaranUsecase) GetDetailPembayaranByIndex(index int) (entities.MetodePembayaran, error) {
	detailPembayaran, err := mpu.RepoMetodePembayaran.GetByIndex(index)

	if err != nil {
		return detailPembayaran, err
	}

	return detailPembayaran, nil
}

func (mpu *metodePembayaranUsecase) GetDetailPembayaranByID(id uint) (entities.MetodePembayaran, error) {
	detailPembayaran, err := mpu.RepoMetodePembayaran.GetByID(id)

	if err != nil {
		return detailPembayaran, err
	}

	return detailPembayaran, nil
}

func (mpu *metodePembayaranUsecase) PatchDetailPembayaranByID(id uint, input *model.InputMetodePembayaran) error {
	patchPayment := entities.MetodePembayaran{
		IndexPembayaran:      &input.IndexPembayaran,
		JenisPembayaran:      input.JenisPembayaran,
		KredensialPembayaran: input.KredensialPembayaran,
		Pemilik:              input.Pemilik,
	}
	if err := mpu.RepoMetodePembayaran.UpdateKredensialByID(id, patchPayment); err != nil {
		return err
	}
	return nil
}

func (mpu *metodePembayaranUsecase) DeletePembayaranByID(id uint) error {
	if err := mpu.RepoMetodePembayaran.DeleteByID(id); err != nil {
		return err
	}
	return nil
}
