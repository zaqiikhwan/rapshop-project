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
func (spdl *servicePembelianDL) CreateDataPembelianMidtrans(input entities.PembelianDL) error {
	return nil
}

func (spdl *servicePembelianDL) CreateDataPembelian(world string, nama string, grow_id string, jenis_item bool, jumlah_dl int, wa string, metode_transfer int, gambar string, id string) error {
	location := time.FixedZone("UTC+7", 7*60*60)
	GMT_7 := time.Now().In(location)

	hargaBeli, err := spdl.ServiceStockDL.GetLatestDataStock()

	if err != nil {
		return err
	}

	var jumlahTransaksi int

	if jumlah_dl > 0 && jumlah_dl < 100 {
		jumlahTransaksi = jumlah_dl * hargaBeli.HargaBeliDL
	} else if jumlah_dl % 100 == 0 && jumlah_dl > 0 {
		jumlahTransaksi = (jumlah_dl / 100) * hargaBeli.HargaBeliBGL
	} else if jumlah_dl > 100 {
		jumlahTransaksi = (jumlah_dl / 100) * hargaBeli.HargaBeliBGL + jumlah_dl % 100 * hargaBeli.HargaBeliDL
	}

	newPembelian := entities.PembelianDL{
		ID: id,
		World: world,
		Nama: nama,
		GrowID: grow_id,
		JenisItem: jenis_item,
		JumlahDL: jumlah_dl,
		WA: wa,
		HargaBeli: hargaBeli.HargaBeliDL,
		MetodeTransfer: metode_transfer,
		JumlahTransaksi: int64(jumlahTransaksi),
		BuktiPembayaran: gambar,
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

	dataPembelian, err := spdl.RepoPembelianDL.GetByID(id)
	if err != nil {
		return err
	}

	kurangiStock := model.InputStockDL {
		StockDL: dataPembelian.JumlahDL,
	}

	if *statusKirim.StatusPengiriman {
		if _,err := spdl.ServiceStockDL.UpdateKurangiStock(&kurangiStock); err != nil {
			return err
		}
	}
	return nil
}

func(spdl *servicePembelianDL) UpdateStatusButtonBayar(id string, input entities.PembelianDL) error {
	statusBayar := entities.PembelianDL {
		ButtonBayar: input.ButtonBayar,
	}
	if err := spdl.RepoPembelianDL.UpdateStatus(statusBayar, id); err != nil {
		return err
	}
	return nil
}

func(spdl *servicePembelianDL) UpdateStatusPembayaranAdmin(id string, input entities.PembelianDL) error {
	statusBayar := entities.PembelianDL {
		StatusPembayaran: input.StatusPembayaran,
	}
	if err := spdl.RepoPembelianDL.UpdateStatus(statusBayar, id); err != nil {
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