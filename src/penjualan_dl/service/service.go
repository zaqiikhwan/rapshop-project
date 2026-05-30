package service

import (
	"errors"
	"mime/multipart"
	"os"
	"rapsshop-project/entities"
	"rapsshop-project/model"
	"time"

	storage_go "github.com/supabase-community/storage-go"
	"gorm.io/gorm"
)

type penjualanDLUsecase struct {
	db *gorm.DB
	PenjualanDLRepository model.PenjualanDLRepository
	StockRepo model.StockDLRepository
}

func NewPenjualanDLUsecase(db *gorm.DB, repoJualDL model.PenjualanDLRepository, stockRepo model.StockDLRepository) model.PenjualanDLUsecase {
	return &penjualanDLUsecase{db: db, PenjualanDLRepository: repoJualDL, StockRepo: stockRepo}
}

func (pdlu *penjualanDLUsecase) Create(image *multipart.FileHeader, jumlahDL int, jumlahTransaksi int, wa string, transfer string, nomorTransfer string, nama string, hargaJualDL int) error {
	if jumlahDL <= 0 {
		return errors.New("jumlah_dl must be greater than 0")
	}

	client := storage_go.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SERVICE_TOKEN"), nil)

	if client == nil {
		return errors.New("storage authentication failed")
	}

	imageIo, err := image.Open()
	if err != nil {
		return err
	}
	defer imageIo.Close()

	if resp := client.UploadFile(os.Getenv("STORAGE_NAME"), image.Filename, imageIo); resp.Key == "" {
		return errors.New("upload failed: " + resp.Message)
	}

	location := time.FixedZone("UTC+7", 7*60*60)
	GMT_7 := time.Now().In(location)
	newPenjualan := entities.PenjualanDL{
		Nama: nama,
		BuktiDL: os.Getenv("BASE_URL") + image.Filename,
		HargaJual: hargaJualDL,
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

func (pdlu *penjualanDLUsecase) GetAll(_startInt int, _endInt int) ([]entities.PenjualanDL, int, error) {
	allPenjualan, lenData, err := pdlu.PenjualanDLRepository.GetAll(_startInt, _endInt)

	if err != nil {
		return allPenjualan, lenData, err
	}

	return allPenjualan, lenData, nil
}

func (pdlu *penjualanDLUsecase) GetByDate(date string) ([]model.RekapTransaksiPenjualan, []model.RekapTransaksiPembelian, error) {
	dataPenjualanByDate, rekapPembelian, err := pdlu.PenjualanDLRepository.GetByDate(date)

	if err != nil {
		return dataPenjualanByDate, rekapPembelian, err
	}

	return dataPenjualanByDate, rekapPembelian, nil
}

func (pdlu *penjualanDLUsecase) GetTotal(date string) ([]model.RekapTotalPenjualan, error) {
	dataPenjualanByDate, err := pdlu.PenjualanDLRepository.GetTotalPenjualan(date)

	if err != nil {
		return dataPenjualanByDate, err
	}

	return dataPenjualanByDate, nil
}

func (pdlu *penjualanDLUsecase) GetProfit(date string) ([]model.RekapProfit, error) {
	profit, err := pdlu.PenjualanDLRepository.GetProfit(date)

	if err != nil {
		return profit, err
	}

	return profit, nil
}

func (pdlu *penjualanDLUsecase) GetByID(id uint) (entities.PenjualanDL, error) {
	detail, err := pdlu.PenjualanDLRepository.GetByID(id)

	if err != nil {
		return detail, err
	}

	return detail, nil
}

func (pdlu *penjualanDLUsecase) UpdateByID(id uint, input entities.PenjualanDL) (entities.PenjualanDL, error) {
	detail, err := pdlu.PenjualanDLRepository.GetByID(id)
	if err != nil {
		return detail, err
	}

	if input.Status == nil {
		if err := pdlu.PenjualanDLRepository.UpdateByID(id, entities.PenjualanDL{EditorStatus: input.EditorStatus}); err != nil {
			return detail, err
		}
		return pdlu.PenjualanDLRepository.GetByID(id)
	}

	from, to := *detail.Status, *input.Status

	// Compare-and-set the status, then adjust stock only when this call actually
	// performed the transition, so concurrent double-approve can't double-count stock.
	err = pdlu.db.Transaction(func(tx *gorm.DB) error {
		changed, err := pdlu.PenjualanDLRepository.WithTx(tx).UpdateStatusIfCurrent(id, from, to, input.EditorStatus)
		if err != nil || !changed {
			return err
		}
		switch {
		case (from == 0 || from == -1) && to == 1:
			_, err = pdlu.StockRepo.WithTx(tx).AdjustLatest(detail.JumlahDL, entities.StockDL{})
		case from == 1 && (to == 0 || to == -1):
			_, err = pdlu.StockRepo.WithTx(tx).AdjustLatest(-detail.JumlahDL, entities.StockDL{})
		}
		return err
	})
	if err != nil {
		return detail, err
	}

	return pdlu.PenjualanDLRepository.GetByID(id)
}

func (pdlu *penjualanDLUsecase) DeleteByID(id uint) error {
	err := pdlu.PenjualanDLRepository.DeleteByID(id)

	if err != nil {
		return err
	}

	return nil
}