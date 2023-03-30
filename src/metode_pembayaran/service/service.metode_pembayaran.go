package service

import "rapsshop-project/entities"

type metodePembayaranUsecase struct {
	RepoMetodePembayaran entities.MetodePembayaranRepository
}

func NewMetodePembayaranUsecase(repoMetodePembayaran entities.MetodePembayaranRepository) entities.MetodePembayaranUsecase {
	return &metodePembayaranUsecase{RepoMetodePembayaran: repoMetodePembayaran}
}

func (mpu *metodePembayaranUsecase) CreateNewPembayaran(input *entities.InputMetodePembayaran) error {
	newPembayaran := entities.MetodePembayaran{
		IndexPembayaran: &input.IndexPembayaran,
		JenisPembayaran: input.JenisPembayaran,
		KredensialPembayaran: input.KredensialPembayaran,
		Pemilik: input.Pemilik,
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

func (mpu *metodePembayaranUsecase) PatchDetailPembayaranByID(id uint, input *entities.InputMetodePembayaran) error {
	patchPayment := entities.MetodePembayaran {
		IndexPembayaran: &input.IndexPembayaran,
		JenisPembayaran: input.JenisPembayaran,
		KredensialPembayaran: input.KredensialPembayaran,
		Pemilik: input.Pemilik,
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