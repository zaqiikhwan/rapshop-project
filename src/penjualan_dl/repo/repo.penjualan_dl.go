package repo

import (
	"rapsshop-project/entities"
	"rapsshop-project/model"

	"gorm.io/gorm"
)

type penjualanDLRepository struct {
	db *gorm.DB
}

func NewPenjualanDLRepository(db *gorm.DB) model.PenjualanDLRepository {
	return &penjualanDLRepository{db: db}
}

func (pdlr *penjualanDLRepository) Create(input entities.PenjualanDL) error {
	if err := pdlr.db.Create(&input).Error; err != nil {
		return err
	}
	return nil
}

func (pdlr *penjualanDLRepository) GetAll(_startInt int, _endInt int) ([]entities.PenjualanDL, int,error) {
	var allPenjualan []entities.PenjualanDL
	var lenData []entities.PenjualanDL
	if err := pdlr.db.Select("id").Find(&lenData).Error; err != nil {
		return allPenjualan, 0, err
	}

	if err := pdlr.db.Order("created_at desc").Offset(_startInt - 1).Limit(_endInt - _startInt + 1).Find(&allPenjualan).Error; err != nil {
		return allPenjualan, 0, err
	}
	return allPenjualan, len(lenData), nil
}

func (pdlr *penjualanDLRepository) GetByDate(date string) ([]model.RekapTransaksiPenjualan, []model.RekapTransaksiPembelian, error) {
	var allPenjualanByDate []entities.PenjualanDL
	var allPembelianByDate []entities.PembelianDL

	var rekapBeliTunggal model.RekapTransaksiPembelian
	var rekapBeli []model.RekapTransaksiPembelian
	
	var rekapJualTunggal model.RekapTransaksiPenjualan
	var rekapJual []model.RekapTransaksiPenjualan
	
	var arrayHarga []int
	var arrayHargaBeli []int

	query := "%" + date + "%"

	if err := pdlr.db.Select("harga_jual, status").Where("created_at LIKE ? and (harga_jual != 0 or harga_jual <> NULL) and status = 1 order by harga_jual asc", query).Find(&allPenjualanByDate).Error; err != nil {
		return rekapJual, rekapBeli, err
	}

	if err := pdlr.db.Select("harga_beli").Where("created_at LIKE ? and (harga_beli != 0 or harga_beli <> NULL) and status_pengiriman = 1 order by harga_beli asc", query).Find(&allPembelianByDate).Error; err != nil {
		return rekapJual, rekapBeli, err
	}

	for _, v := range allPenjualanByDate {
		if len(arrayHarga) == 0 {
			arrayHarga  = append(arrayHarga, v.HargaJual)
		} else if len(arrayHarga) > 0 && arrayHarga[len(arrayHarga) - 1] != v.HargaJual {
			arrayHarga = append(arrayHarga, v.HargaJual)
		}
	}

	for _, v := range allPembelianByDate {
		if len(arrayHargaBeli) == 0 {
			arrayHargaBeli = append(arrayHargaBeli, v.HargaBeli)
		} else if len(arrayHargaBeli) > 0 && arrayHargaBeli[len(arrayHargaBeli) - 1] != v.HargaBeli {
			arrayHargaBeli = append(arrayHargaBeli, v.HargaBeli)
		}
	}
	
	for i := 0; i < len(arrayHarga); i++ {
	
		if err := pdlr.db.Raw("select sum(jumlah_transaksi) from penjualan_dls where created_at LIKE ? and harga_jual = ? and status = ?", query, arrayHarga[i], 1).Scan(&rekapJualTunggal.JumlahTransaksi).Error; err != nil {
			rekapJualTunggal.JumlahTransaksi = 0
		}
	
		if err := pdlr.db.Raw("select sum(jumlah_dl) from penjualan_dls where created_at LIKE ? and harga_jual = ? and status = ?", query, arrayHarga[i], 1).Scan(&rekapJualTunggal.JumlahDL).Error; err != nil {
			rekapJualTunggal.JumlahDL = 0
		}

		rekapJualTunggal.Rate = arrayHarga[i]
		rekapJual = append(rekapJual, rekapJualTunggal)
	}

	for i := 0; i < len(arrayHargaBeli); i++ {
		if err := pdlr.db.Raw("select sum(jumlah_transaksi) from pembelian_dls where created_at LIKE ? and harga_beli = ? and status_pembayaran = 'success'", query, arrayHargaBeli[i]).Scan(&rekapBeliTunggal.JumlahTransaksi).Error; err != nil {
			rekapBeliTunggal.JumlahTransaksi = 0
		}
	
		if err := pdlr.db.Raw("select sum(jumlah_dl) from pembelian_dls where created_at LIKE ? and harga_beli = ? and status_pembayaran = 'success'", query, arrayHargaBeli[i]).Scan(&rekapBeliTunggal.JumlahDL).Error; err != nil {
			rekapBeliTunggal.JumlahDL = 0
		}

		rekapBeliTunggal.Rate = arrayHargaBeli[i]
		rekapBeli = append(rekapBeli, rekapBeliTunggal)
	}
	
	return rekapJual, rekapBeli, nil
}

func (pdlr *penjualanDLRepository) GetTotalPenjualan(date string) ([]model.RekapTotalPenjualan, error) {
	var allPenjualanByDate []entities.PenjualanDL

	var totalPenjualanTunggal model.RekapTotalPenjualan
	var totalPenjualan []model.RekapTotalPenjualan

	query := "%" + date + "%"

	if err := pdlr.db.Select("created_at").Where("created_at LIKE ? order by created_at asc", query).Find(&allPenjualanByDate).Error; err != nil {
		return totalPenjualan, err
	}

	var arrayTanggalStr []string

	for _, v := range allPenjualanByDate {
		if len(arrayTanggalStr) == 0 {
			arrayTanggalStr = append(arrayTanggalStr, v.CreatedAt.String()[0:10])
		} else if len(arrayTanggalStr) > 0 && arrayTanggalStr[len(arrayTanggalStr) - 1] != v.CreatedAt.String()[0:10] {
			arrayTanggalStr = append(arrayTanggalStr, v.CreatedAt.String()[0:10])
		}
	}

	for i := 0; i < len(arrayTanggalStr); i++ {

		query := "%" + arrayTanggalStr[i] + "%"
	
		if err := pdlr.db.Raw("select sum(jumlah_dl) from penjualan_dls where created_at LIKE ? and status = 1", query).Scan(&totalPenjualanTunggal.JumlahDL).Error; err != nil {
			totalPenjualanTunggal.JumlahDL = 0
		}

		totalPenjualanTunggal.Tanggal = arrayTanggalStr[i][8:10]
		totalPenjualan = append(totalPenjualan, totalPenjualanTunggal)
	}
	
	return totalPenjualan, nil
}

func (pdlr *penjualanDLRepository) GetProfit(date string) ([]model.RekapProfit, error){
	var allPenjualanByDate []entities.PenjualanDL
	var allPembelianByDate []entities.PembelianDL

	var eachProfit model.RekapProfit
	var eachTransaksi model.Total
	var allProfit []model.RekapProfit

	query := "%" + date + "%"

	if err := pdlr.db.Select("created_at").Where("created_at LIKE ? order by created_at asc", query).Find(&allPenjualanByDate).Error; err != nil {
		return allProfit, err
	}

	if err := pdlr.db.Select("created_at").Where("created_at LIKE ? order by created_at asc", query).Find(&allPembelianByDate).Error; err != nil {
		return allProfit, err
	}

	var arrayTanggalJual []string

	for _, v := range allPenjualanByDate {
		if len(arrayTanggalJual) == 0 {
			arrayTanggalJual = append(arrayTanggalJual, v.CreatedAt.String()[0:10])
		} else if len(arrayTanggalJual) > 0 && arrayTanggalJual[len(arrayTanggalJual) - 1] != v.CreatedAt.String()[0:10] {
			arrayTanggalJual = append(arrayTanggalJual, v.CreatedAt.String()[0:10])
		}
	}

	var arrayTanggalBeli []string

	for _, v := range allPembelianByDate {
		if len(arrayTanggalBeli) == 0 {
			arrayTanggalBeli = append(arrayTanggalBeli, v.CreatedAt.String()[0:10])
		} else if len(arrayTanggalBeli) > 0 && arrayTanggalBeli[len(arrayTanggalBeli) - 1] != v.CreatedAt.String()[0:10] {
			arrayTanggalBeli = append(arrayTanggalBeli, v.CreatedAt.String()[0:10])
		}
	}


	for i := 0; i < len(arrayTanggalBeli); i++ {

		queryDate := "%" + arrayTanggalBeli[i] + "%"
		if err := pdlr.db.Raw("select sum(jumlah_transaksi) from penjualan_dls where created_at LIKE ? and status = 1", queryDate).Scan(&eachTransaksi.TransaksiJual).Error; err != nil {
			eachTransaksi.TransaksiBeli = 0
			eachTransaksi.TransaksiJual = 0
		}

		if err := pdlr.db.Raw("select sum(jumlah_transaksi) from pembelian_dls where created_at LIKE ? and status_pembayaran = 'success'", queryDate).Scan(&eachTransaksi.TransaksiBeli).Error; err != nil {
			eachTransaksi.TransaksiBeli = 0
			eachTransaksi.TransaksiJual = 0
		}

		eachProfit.Tanggal = arrayTanggalBeli[i][8:10]
		eachProfit.Profit = eachTransaksi.TransaksiBeli - eachTransaksi.TransaksiJual

		allProfit = append(allProfit, eachProfit)
	}

	return allProfit, nil
}

func (pdlr *penjualanDLRepository) GetByID(id uint) (entities.PenjualanDL, error) {
	var penjualan entities.PenjualanDL
	if err := pdlr.db.First(&penjualan, id).Error; err != nil {
		return penjualan, err
	}
	return penjualan, nil
}

func (pdlr *penjualanDLRepository) UpdateByID(id uint, input entities.PenjualanDL) error {
	var penjualan entities.PenjualanDL
	if err := pdlr.db.Where("id = ?", id).Model(&penjualan).Updates(&input).Error; err != nil {
		return err
	}
	return nil
}

func (pdlr *penjualanDLRepository) DeleteByID(id uint) error {
	var penjualan entities.PenjualanDL
	if err := pdlr.db.Delete(&penjualan, id).Error; err != nil {
		return  err
	}
	return nil
}