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
	if err := pdlr.db.Find(&lenData).Error; err != nil {
		return allPenjualan, 0, err
	}

	if err := pdlr.db.Order("created_at desc").Offset(_startInt - 1).Limit(_endInt - _startInt + 1).Find(&allPenjualan).Error; err != nil {
		return allPenjualan, 0, err
	}
	return allPenjualan, len(lenData), nil
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