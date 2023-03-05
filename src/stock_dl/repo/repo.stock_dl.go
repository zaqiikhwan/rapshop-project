package repo

import (
	"rapsshop-project/entities"
	"rapsshop-project/model"

	"gorm.io/gorm"
)

type stockDLRepository struct {
	db *gorm.DB
}

func NewStockDLRepository(db *gorm.DB) model.StockDLRepository {
	return &stockDLRepository{db: db}
}

func (sdlr *stockDLRepository) Create(newStock entities.StockDL) error {
	if err := sdlr.db.Create(&newStock).Error; err != nil {
		return err
	}
	return nil
}

func (sdlr *stockDLRepository) GetAll() ([]entities.StockDL, error) {
	var allStock []entities.StockDL
	if err := sdlr.db.Model(&allStock).Find(&allStock).Error; err != nil {
		return allStock, err
	}
	return allStock, nil
}

func (sdlr *stockDLRepository) GetLatest() (entities.StockDL, error) {
	var detail entities.StockDL
	if err := sdlr.db.Order("id desc").First(&detail).Error; err != nil {
		return detail, err
	}
	return detail, nil
}

func (sdlr *stockDLRepository) UpdateByID(updateStock entities.StockDL, id uint) error {
	var stock entities.StockDL
	if err := sdlr.db.Model(&stock).Where("id = ?", id).Updates(updateStock).Error; err != nil {
		return err
	}
	return nil
}

func (sdlr *stockDLRepository) DeleteByID(id uint) error {
	var detailStock entities.StockDL
	if err := sdlr.db.Delete(&detailStock, id).Error; err != nil {
		return err
	}

	return nil
}