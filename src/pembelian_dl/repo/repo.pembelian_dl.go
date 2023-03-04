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