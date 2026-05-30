package repo

import (
	"errors"
	"rapsshop-project/entities"
	"rapsshop-project/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type stockDLRepository struct {
	db *gorm.DB
}

func NewStockDLRepository(db *gorm.DB) model.StockDLRepository {
	return &stockDLRepository{db: db}
}

func (sdlr *stockDLRepository) WithTx(tx *gorm.DB) model.StockDLRepository {
	return &stockDLRepository{db: tx}
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

// AdjustLatest atomically applies `delta` to the latest stock row (locking it for
// the duration of the transaction) and applies non-zero price fields from `prices`.
// It refuses to drive stock below zero.
func (sdlr *stockDLRepository) AdjustLatest(delta int, prices entities.StockDL) (entities.StockDL, error) {
	var latest entities.StockDL
	err := sdlr.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Order("id desc").First(&latest).Error; err != nil {
			return err
		}
		newStock := latest.StockDL + delta
		if newStock < 0 {
			return errors.New("insufficient stock")
		}
		if err := tx.Model(&entities.StockDL{}).Where("id = ?", latest.ID).Updates(prices).Error; err != nil {
			return err
		}
		return tx.Model(&entities.StockDL{}).Where("id = ?", latest.ID).UpdateColumn("stock_dl", newStock).Error
	})
	if err != nil {
		return latest, err
	}
	return sdlr.GetLatest()
}

func (sdlr *stockDLRepository) DeleteByID(id uint) error {
	var detailStock entities.StockDL
	if err := sdlr.db.Delete(&detailStock, id).Error; err != nil {
		return err
	}

	return nil
}
