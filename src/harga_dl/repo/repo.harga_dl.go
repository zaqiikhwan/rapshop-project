package repo

import (
	"rapsshop-project/entities"
	"rapsshop-project/model"

	"gorm.io/gorm"
)

type hargaDLRepository struct {
	db *gorm.DB
}

func NewHargaDLRepository(db *gorm.DB) model.HargaDLRepository {
	return &hargaDLRepository{db: db}
}

func (hdlr *hargaDLRepository) Create(newHarga entities.HargaDL) error {
	if err := hdlr.db.Create(&newHarga).Error; err != nil {
		return err
	}
	return nil
}

// func (hdlr *hargaDLRepository) GetAll() ([]entities.StockDL, error) {
// 	var allStock []entities.StockDL
// 	if err := hdlr.db.Model(&allStock).Find(&allStock).Error; err != nil {
// 		return allStock, err
// 	}
// 	return allStock, nil
// }

func (hdlr *hargaDLRepository) GetLatest() (entities.HargaDL, error) {
	var detail entities.HargaDL
	if err := hdlr.db.Order("id desc").First(&detail).Error; err != nil {
		return detail, err
	}
	return detail, nil
}

func (hdlr *hargaDLRepository) UpdateByID(updateStock entities.HargaDL, id uint) error {
	var updatePrice entities.HargaDL
	if err := hdlr.db.Model(&updatePrice).Where("id = ?", id).Updates(updateStock).Error; err != nil {
		return err
	}
	return nil
}

func (hdlr *hargaDLRepository) DeleteByID(id uint) error {
	var detailPrice entities.HargaDL
	if err := hdlr.db.Delete(&detailPrice, id).Error; err != nil {
		return err
	}

	return nil
}