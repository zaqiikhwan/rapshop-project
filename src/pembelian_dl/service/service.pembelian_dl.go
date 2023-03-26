package service

import (
	"rapsshop-project/entities"
	"rapsshop-project/lib"
	"rapsshop-project/model"
	"time"
)

type servicePembelianDL struct {
	RepoPembelianDL model.PembelianDLRepository
	ServiceStockDL model.StockDLUsecase
	midtransCoreClient *lib.CoreApi
}

func NewServicePembelianDL(repoBeliDL model.PembelianDLRepository, ca *lib.CoreApi, serviceStock model.StockDLUsecase) model.PembelianDLUsecase {
	return &servicePembelianDL{
		RepoPembelianDL: repoBeliDL,
		ServiceStockDL: serviceStock,
		midtransCoreClient: ca,
	}
}

func (spdl *servicePembelianDL) CreateDataPembelian(input entities.PembelianDL) error {
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
		HargaBeli: input.HargaBeli,
		MetodeTransfer: input.MetodeTransfer,
		JumlahTransaksi: input.JumlahTransaksi,
		CreatedAt: GMT_7,
	}

	if err := spdl.RepoPembelianDL.Create(newPembelian); err != nil {
		return err
	}
	return nil
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
	dataPembelian, err := spdl.RepoPembelianDL.GetByID(id)
	if err != nil {
		return err
	}
	if midtransReport != nil {
		if midtransReport.TransactionStatus == "capture" {
			if midtransReport.FraudStatus == "challenge" {
				dataPembelian.StatusPembayaran = "challange"
			} else if midtransReport.FraudStatus == "accept" {
				dataPembelian.StatusPembayaran = "success"
				kurangiStock := model.InputStockDL {
					StockDL: dataPembelian.JumlahDL,
				}
				if _,err := spdl.ServiceStockDL.UpdateKurangiStock(&kurangiStock); err != nil {
					return err
				}
			}
		} else if midtransReport.TransactionStatus == "settlement" {
			dataPembelian.StatusPembayaran = "success"
			kurangiStock := model.InputStockDL {
				StockDL: dataPembelian.JumlahDL,
			}
			if _,err := spdl.ServiceStockDL.UpdateKurangiStock(&kurangiStock); err != nil {
				return err
			}
		} else if midtransReport.TransactionStatus == "deny" {
			dataPembelian.StatusPembayaran = "deny"
		} else if midtransReport.TransactionStatus == "cancel" || midtransReport.TransactionStatus == "expire" {
			dataPembelian.StatusPembayaran = "failure"
		} else if midtransReport.TransactionStatus == "pending" {
			dataPembelian.StatusPembayaran = "pending"
		}
	}
	
	if err := spdl.RepoPembelianDL.UpdateStatus(dataPembelian, id); err != nil {
		return err
	} 
	return nil
}

func(spdl *servicePembelianDL) UpdateStatusPengiriman(id string, input entities.PembelianDL) error {
	statusKirim := entities.PembelianDL {
		EditorStatus: input.EditorStatus,
		StatusPengiriman: input.StatusPengiriman,
	}
	if err := spdl.RepoPembelianDL.UpdateStatus(statusKirim, id); err != nil {
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