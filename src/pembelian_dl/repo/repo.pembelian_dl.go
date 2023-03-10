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
	if err := rp.db.Order("created_at desc").Where("status_pembayaran = ?", "success").Offset(_startInt - 1).Limit(_endInt - _startInt + 1).Find(&allPembelian).Error; err != nil {
		return allPembelian, 0, err
	}
	if err := rp.db.Where("status_pembayaran = ?", "success").Find(&lenData).Error; err != nil {
		return allPembelian, 0, err
	}
	return allPembelian, len(lenData), nil
}

func (rp *repoPembelianDL) GetByID(id string) (entities.PembelianDL, error) {
	var model entities.PembelianDL
	if err := rp.db.First(&model).Where("id = ?", id).Take(&model).Error; err != nil {
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