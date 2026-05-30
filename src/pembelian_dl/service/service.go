package service

import (
	"errors"
	"fmt"
	"rapsshop-project/entities"
	"rapsshop-project/lib"
	"rapsshop-project/model"
	"time"

	"gorm.io/gorm"
)

type servicePembelianDL struct {
	db *gorm.DB
	RepoPembelianDL model.PembelianDLRepository
	StockRepo model.StockDLRepository
	midtransCoreClient *lib.CoreAPI
}

func NewServicePembelianDL(db *gorm.DB, repoBeliDL model.PembelianDLRepository, ca *lib.CoreAPI, stockRepo model.StockDLRepository) model.PembelianDLUsecase {
	return &servicePembelianDL{
		db: db,
		RepoPembelianDL: repoBeliDL,
		StockRepo: stockRepo,
		midtransCoreClient: ca,
	}
}

// resolvePaymentType maps the numeric MetodeTransfer to a Midtrans payment_type
// and (for bank transfers) the bank code.
func resolvePaymentType(metode int) (paymentType string, bank string) {
	switch metode {
	case 1:
		return "qris", ""
	case 2:
		return "gopay", ""
	case 3:
		return "shopeepay", ""
	case 4:
		return "bank_transfer", "bca"
	case 5:
		return "bank_transfer", "bri"
	case 6:
		return "bank_transfer", "bni"
	}
	return "", ""
}

// ChargeAndCreate builds the Midtrans charge, sends it via the lib driver, and
// persists the order on success. All Midtrans/HTTP concerns live in lib.
func (spdl *servicePembelianDL) ChargeAndCreate(input entities.PembelianDL) (map[string]any, error) {
	if input.JumlahDL <= 0 {
		return nil, errors.New("jumlah_dl must be greater than 0")
	}
	harga, err := spdl.StockRepo.GetLatest()
	if err != nil {
		return nil, err
	}

	paymentType, bank := resolvePaymentType(input.MetodeTransfer)
	md := model.NewMidtransData(paymentType, bank, input, harga)
	payload, total := md.IniDataPembelian()

	resp, status, err := spdl.midtransCoreClient.Charge(payload)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		msg, _ := resp["status_message"].(string)
		return resp, fmt.Errorf("midtrans charge failed: %s", msg)
	}

	input.JumlahTransaksi = total
	input.HargaBeli = harga.HargaBeliDL
	if err := spdl.CreateDataPembelian(input); err != nil {
		return nil, err
	}
	return resp, nil
}

func (spdl *servicePembelianDL) GetLiveStatus(id string) (map[string]any, error) {
	return spdl.midtransCoreClient.CheckStatus(id)
}

func (spdl *servicePembelianDL) CreateDataPembelian(input entities.PembelianDL) error {
	if input.JumlahDL <= 0 {
		return errors.New("jumlah_dl must be greater than 0")
	}
	location := time.FixedZone("UTC+7", 7*60*60)
	GMT_7 := time.Now().In(location)
	newPembelian := entities.PembelianDL{
		ID: input.ID,
		World: input.World,
		Nama: input.Nama,
		GrowID: input.GrowID,
		JenisItem: input.JenisItem,
		JumlahDL: input.JumlahDL,
		WA: input.WA,
		MetodeTransfer: input.MetodeTransfer,
		JumlahTransaksi: input.JumlahTransaksi,
		HargaBeli: input.HargaBeli,
		CreatedAt: GMT_7,
	}

	if err := spdl.RepoPembelianDL.Create(newPembelian); err != nil {
		return err
	}
	return nil
}

func (spdl *servicePembelianDL) CreateDataPembelianManual(input entities.PembelianDL) error {
	if input.JumlahDL <= 0 {
		return errors.New("jumlah_dl must be greater than 0")
	}

	hargaBeli, err := spdl.StockRepo.GetLatest()
	if err != nil {
		return err
	}

	_, total := model.NewMidtransData("", "", input, hargaBeli).IniDataPembelian()

	location := time.FixedZone("UTC+7", 7*60*60)
	newPembelian := entities.PembelianDL{
		ID: input.ID,
		World: input.World,
		Nama: input.Nama,
		GrowID: input.GrowID,
		JenisItem: input.JenisItem,
		JumlahDL: input.JumlahDL,
		WA: input.WA,
		HargaBeli: hargaBeli.HargaBeliDL,
		MetodeTransfer: input.MetodeTransfer,
		JumlahTransaksi: total,
		BuktiPembayaran: input.BuktiPembayaran,
		CreatedAt: time.Now().In(location),
	}

	return spdl.RepoPembelianDL.Create(newPembelian)
}

func (spdl *servicePembelianDL) GetAllPembelian(_startInt int, _endInt int) ([]entities.PembelianDL, int, error) {
	allData, lenData, err := spdl.RepoPembelianDL.GetAll(_startInt, _endInt)
	if err != nil {
		return allData, lenData, err
	}

	return allData, lenData, nil
}

func(spdl *servicePembelianDL) UpdateStatusPembayaran(id string) error {
	midtransReport, err := spdl.midtransCoreClient.HandleNotification(id)
	if err != nil {
		return err
	}
	if midtransReport == nil {
		return nil
	}
	dataPembelian, err := spdl.RepoPembelianDL.GetByID(id)
	if err != nil {
		return err
	}

	newStatus := ""
	switch midtransReport.TransactionStatus {
	case "capture":
		if midtransReport.FraudStatus == "challenge" {
			newStatus = "challange"
		} else if midtransReport.FraudStatus == "accept" {
			newStatus = "success"
		}
	case "settlement":
		newStatus = "success"
	case "deny":
		newStatus = "deny"
	case "cancel", "expire":
		newStatus = "failure"
	case "pending":
		newStatus = "pending"
	}
	if newStatus == "" {
		return nil
	}
	if newStatus != "success" {
		return spdl.RepoPembelianDL.UpdateByID(entities.PembelianDL{StatusPembayaran: newStatus}, id)
	}

	// Flip to paid and decrement stock exactly once: MarkPaid only reports a
	// transition for the call that actually changes the row, so duplicate/concurrent
	// webhooks can't double-deduct.
	return spdl.db.Transaction(func(tx *gorm.DB) error {
		paid, err := spdl.RepoPembelianDL.WithTx(tx).MarkPaid(id)
		if err != nil || !paid {
			return err
		}
		_, err = spdl.StockRepo.WithTx(tx).AdjustLatest(-dataPembelian.JumlahDL, entities.StockDL{})
		return err
	})
}

func(spdl *servicePembelianDL) UpdateStatusPengiriman(id string, input entities.PembelianDL) error {
	// Decrement stock exactly once, only on the not-shipped -> shipped transition.
	if input.StatusPengiriman != nil && *input.StatusPengiriman {
		return spdl.db.Transaction(func(tx *gorm.DB) error {
			shipped, err := spdl.RepoPembelianDL.WithTx(tx).MarkShipped(id, input.EditorStatus)
			if err != nil || !shipped {
				return err
			}
			dataPembelian, err := spdl.RepoPembelianDL.WithTx(tx).GetByID(id)
			if err != nil {
				return err
			}
			_, err = spdl.StockRepo.WithTx(tx).AdjustLatest(-dataPembelian.JumlahDL, entities.StockDL{})
			return err
		})
	}
	return spdl.RepoPembelianDL.UpdateByID(entities.PembelianDL{
		EditorStatus: input.EditorStatus,
		StatusPengiriman: input.StatusPengiriman,
	}, id)
}

func(spdl *servicePembelianDL) UpdateStatusButtonBayar(id string, input entities.PembelianDL) error {
	statusBayar := entities.PembelianDL {
		ButtonBayar: input.ButtonBayar,
	}
	if err := spdl.RepoPembelianDL.UpdateByID(statusBayar, id); err != nil {
		return err
	}
	return nil
}

func(spdl *servicePembelianDL) UpdateStatusPembayaranAdmin(id string, input entities.PembelianDL) error {
	statusBayar := entities.PembelianDL {
		EditorStatus: input.EditorStatus,
		StatusPembayaran: input.StatusPembayaran,
	}
	if err := spdl.RepoPembelianDL.UpdateByID(statusBayar, id); err != nil {
		return err
	}
	return nil
}

func (spdl *servicePembelianDL) UpdateTambahBukti(id string, input entities.PembelianDL) error {
	buktiBayar := entities.PembelianDL {
		BuktiPembayaran: input.BuktiPembayaran,
	}

	if err := spdl.RepoPembelianDL.UpdateByID(buktiBayar, id); err != nil {
		return err
	}
	return nil
}

func(spdl *servicePembelianDL) GetDetailByID(id string) (entities.PembelianDL, error) {
	dataPenjualan, err := spdl.RepoPembelianDL.GetByID(id)
	if err != nil {
		return dataPenjualan, err
	}
	return dataPenjualan, nil
}

func (spdl *servicePembelianDL) GetTotal(date string) ([]model.RekapTotalPembelian, error) {
	dataPembelianByDate, err := spdl.RepoPembelianDL.GetTotalPembelian(date)

	if err != nil {
		return dataPembelianByDate, err
	}

	return dataPembelianByDate, nil
}