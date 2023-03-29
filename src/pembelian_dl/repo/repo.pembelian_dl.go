package repo

import (
	"rapsshop-project/entities"
	"rapsshop-project/model"

	"gorm.io/gorm"
)

type repoPembelianDL struct {
	db *gorm.DB
}

func NewRepoPembelianDL(db *gorm.DB) model.PembelianDLRepository {
	return &repoPembelianDL{db: db}
}

func(rp *repoPembelianDL) Create(input entities.PembelianDL) error {
	if err := rp.db.Create(&input).Error; err != nil {
		return err
	}
	return nil
}

func(rp *repoPembelianDL) GetAll(_startInt int, _endInt int) ([]entities.PembelianDL, int, error) {
	var allPembelian []entities.PembelianDL
	var lenData []entities.PembelianDL
	if err := rp.db.Order("created_at desc").Where("status_pembayaran = ? or button_bayar = ?", "success", true).Offset(_startInt - 1).Limit(_endInt - _startInt + 1).Find(&allPembelian).Error; err != nil {
		return allPembelian, 0, err
	}
	if err := rp.db.Select("id").Where("status_pembayaran = ? or button_bayar = ?", "success", true).Find(&lenData).Error; err != nil {
		return allPembelian, 0, err
	}
	return allPembelian, len(lenData), nil
}

func (rp *repoPembelianDL) GetByID(id string) (entities.PembelianDL, error) {
	var model entities.PembelianDL
	if err := rp.db.Where("id = ?", id).Take(&model).Error; err != nil {
		return model, err
	}
	return model, nil
}

func (rp *repoPembelianDL) UpdateStatus(input entities.PembelianDL, id string) error {
	var statusPayment entities.PembelianDL
	if err := rp.db.Where("id = ?", id).Model(&statusPayment).Updates(input).Error; err != nil {
		return err
	}
	return nil
}

func (rp *repoPembelianDL) GetTotalPembelian(date string) ([]model.RekapTotalPembelian, error) {
	var allPembelianByDate []entities.PembelianDL

	var totalPembelianTunggal model.RekapTotalPembelian
	var totalPembelian []model.RekapTotalPembelian

	query := "%" + date + "%"

	if err := rp.db.Select("created_at").Where("created_at LIKE ? order by created_at asc", query).Find(&allPembelianByDate).Error; err != nil {
		return totalPembelian, err
	}

	var arrayTanggalStr []string
	
	for _, v := range allPembelianByDate {
		if len(arrayTanggalStr) == 0 {
			arrayTanggalStr = append(arrayTanggalStr, v.CreatedAt.String()[0:10])
		} else if len(arrayTanggalStr) > 0 && arrayTanggalStr[len(arrayTanggalStr) - 1] != v.CreatedAt.String()[0:10] {
			arrayTanggalStr = append(arrayTanggalStr, v.CreatedAt.String()[0:10])
		}
	}

	for i := 0; i < len(arrayTanggalStr); i++ {

		queryDate := "%" + arrayTanggalStr[i] + "%"
	
		if err := rp.db.Raw("select sum(jumlah_dl) from pembelian_dls where created_at LIKE ? and status_pembayaran = ?", queryDate, "success").Scan(&totalPembelianTunggal.JumlahDL).Error; err != nil {
			totalPembelianTunggal.JumlahDL = 0
		}

		totalPembelianTunggal.Tanggal = arrayTanggalStr[i][8:10]
		totalPembelian = append(totalPembelian, totalPembelianTunggal)
	}
	
	return totalPembelian, nil
}