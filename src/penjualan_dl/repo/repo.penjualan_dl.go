package repo

import (
	"fmt"
	"rapsshop-project/entities"
	"rapsshop-project/model"
	"strconv"
	"strings"

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

	if err := pdlr.db.Select("harga_beli, status_pembayaran").Where("created_at LIKE ? and (harga_beli != 0 or harga_beli <> NULL) and status_pembayaran = 'success' order by harga_beli asc", query).Find(&allPembelianByDate).Error; err != nil {
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
	var totalPenjualanTunggal model.RekapTotalPenjualan
	var totalPenjualan []model.RekapTotalPenjualan

	arrayDate := strings.Split(date, "-")
	arrayYearInt, _ := strconv.Atoi(arrayDate[0])
	arrayMonthInt, _ := strconv.Atoi(arrayDate[1])
	
	var arrayTanggalStr []string
	if arrayYearInt % 4 == 0 && arrayMonthInt == 2{
		jumlahHari := 29;
		arrayTanggalStr = GeneratorTanggal(arrayDate, jumlahHari)
	} else if arrayMonthInt == 2 {
		jumlahHari := 28;
		arrayTanggalStr = GeneratorTanggal(arrayDate, jumlahHari)
	} else if ((arrayMonthInt % 2 != 0 && (arrayMonthInt != 9 && arrayMonthInt != 11))|| arrayMonthInt == 8 || arrayMonthInt == 10 || arrayMonthInt == 12){
		jumlahHari := 31;
		arrayTanggalStr = GeneratorTanggal(arrayDate, jumlahHari)
	} else {
		jumlahHari := 30;
		arrayTanggalStr = GeneratorTanggal(arrayDate, jumlahHari)
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
	var eachProfit model.RekapProfit
	var eachTransaksi model.Total
	var allProfit []model.RekapProfit

	arrayDate := strings.Split(date, "-")
	arrayYearInt, _ := strconv.Atoi(arrayDate[0])
	arrayMonthInt, _ := strconv.Atoi(arrayDate[1])

	var arrayTanggalStr []string
	if arrayYearInt % 4 == 0 && arrayMonthInt == 2{
		jumlahHari := 29;
		arrayTanggalStr = GeneratorTanggal(arrayDate, jumlahHari)
	} else if arrayMonthInt == 2 {
		jumlahHari := 28;
		arrayTanggalStr = GeneratorTanggal(arrayDate, jumlahHari)
	} else if ((arrayMonthInt % 2 != 0 && (arrayMonthInt != 9 && arrayMonthInt != 11))|| arrayMonthInt == 8 || arrayMonthInt == 10 || arrayMonthInt == 12){
		jumlahHari := 31;
		arrayTanggalStr = GeneratorTanggal(arrayDate, jumlahHari)
	} else {
		jumlahHari := 30;
		arrayTanggalStr = GeneratorTanggal(arrayDate, jumlahHari)
	}

	for i := 0; i < len(arrayTanggalStr); i++ {

		queryDate := "%" + arrayTanggalStr[i] + "%"
		if err := pdlr.db.Raw("select sum(jumlah_transaksi) from penjualan_dls where created_at LIKE ? and status = 1", queryDate).Scan(&eachTransaksi.TransaksiJual).Error; err != nil {
			eachTransaksi.TransaksiBeli = 0
			eachTransaksi.TransaksiJual = 0
		}

		if err := pdlr.db.Raw("select sum(jumlah_transaksi) from pembelian_dls where created_at LIKE ? and status_pembayaran = 'success'", queryDate).Scan(&eachTransaksi.TransaksiBeli).Error; err != nil {
			eachTransaksi.TransaksiBeli = 0
			eachTransaksi.TransaksiJual = 0
		}

		eachProfit.Tanggal = arrayTanggalStr[i][8:10]
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

func GeneratorTanggal(arrayDate []string, jumlahHari int) []string {
	var arrayTanggalStrCheck []string

	for i := 1; i <= jumlahHari; i++ {
		arrDate := [1]int{i}; // {1, 2, 3, 4, 5, 6, 7, 8, 9, 10} // ini buat tanggal
		var arrDateStr [1]string
		var tglStr string
		arrDateStr = [1]string{strconv.Itoa(arrDate[0])}
		if arrDate[0] < 10 {
			combine := [3]string{"0",arrDateStr[0]}
			tglStr = strings.Join(combine[0:2], "")
		} else {
			tglStr = arrDateStr[0]
		}

		fmt.Println(arrayDate[1])
		combine := [3]string{arrayDate[0],arrayDate[1],tglStr}
		arrayTanggalStrCheck = append(arrayTanggalStrCheck, strings.Join(combine[0:3], "-"))
	}	
	return arrayTanggalStrCheck
}