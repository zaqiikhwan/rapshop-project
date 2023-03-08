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

func(rp *repoPembelianDL) GetAll() ([]entities.PembelianDL, error) {
	var allPembelian []entities.PembelianDL
	if err := rp.db.Find(&allPembelian).Error; err != nil {
		return allPembelian, err
	}
	return allPembelian, nil
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