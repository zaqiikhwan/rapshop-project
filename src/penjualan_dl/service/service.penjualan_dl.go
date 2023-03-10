package service

import (
	"errors"
	"mime/multipart"
	"os"
	"rapsshop-project/entities"
	"rapsshop-project/model"
	"time"

	storage_go "github.com/supabase-community/storage-go"
)

type penjualanDLUsecase struct {
	PenjualanDLRepository model.PenjualanDLRepository
}

func NewTestimoniUsecase(repoJualDL model.PenjualanDLRepository) model.PenjualanDLUsecase {
	return &penjualanDLUsecase{PenjualanDLRepository: repoJualDL}
}

func (pdlu *penjualanDLUsecase) Create(image *multipart.FileHeader, jumlahDL int, jumlahTransaksi int, wa string, transfer string, nomorTransfer string) error {

	client := storage_go.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SERVICE_TOKEN"), nil)

	if client == nil {
		return errors.New("storage authentication failed")
	}

	imageIo, err := image.Open()

	client.UploadFile(os.Getenv("STORAGE_NAME"), image.Filename, imageIo)

	if err != nil {
		return err
	}
	location := time.FixedZone("UTC+7", 7*60*60)
	GMT_7 := time.Now().In(location)
	newPenjualan := entities.PenjualanDL{
		BuktiDL: os.Getenv("BASE_URL") + image.Filename,
		JumlahDL: jumlahDL,
		JumlahTransaksi: jumlahTransaksi,
		WA: wa,
		Transfer: transfer,
		NomorTransfer: nomorTransfer,
		CreatedAt: GMT_7,
	}

	if err := pdlu.PenjualanDLRepository.Create(newPenjualan); err != nil {
		return err
	}
	return nil
}

func (pdlu *penjualanDLUsecase) GetAll(_startInt int, _endInt int) ([]entities.PenjualanDL, int,error) {
	allPenjualan, lenData, err := pdlu.PenjualanDLRepository.GetAll(_startInt, _endInt)

	if err != nil {
		return allPenjualan, lenData, err
	}

	return allPenjualan, lenData, nil
}

func (pdlu *penjualanDLUsecase) GetByID(id uint) (entities.PenjualanDL, error) {
	detail, err := pdlu.PenjualanDLRepository.GetByID(id)

	if err != nil {
		return detail, err
	}

	return detail, nil
}

func (pdlu *penjualanDLUsecase) UpdateByID(id uint, input entities.PenjualanDL) (entities.PenjualanDL, error) {
	updateStatus := entities.PenjualanDL{
		StatusPembayaran: input.StatusPembayaran,
	}
	err := pdlu.PenjualanDLRepository.UpdateByID(id, updateStatus)
	if err != nil {
		return updateStatus, err
	}

	updatedData, err := pdlu.PenjualanDLRepository.GetByID(id)

	if err != nil {
		return updatedData, err
	}

	return updatedData, nil
}

func (pdlu *penjualanDLUsecase) DeleteByID(id uint) error {
	err := pdlu.PenjualanDLRepository.DeleteByID(id)

	if err != nil {
		return err
	}

	return nil
}