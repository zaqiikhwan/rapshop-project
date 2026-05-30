package repo

import (
	"rapsshop-project/entities"
	"rapsshop-project/model"

	"gorm.io/gorm"
)

type testimoniRepository struct {
	db *gorm.DB
}

func NewTestimoniRepository(db *gorm.DB) model.TestimoniRepository {
	return &testimoniRepository{db: db}
}

func (tr *testimoniRepository) Create(newTesti entities.Testimoni) error {
	if err := tr.db.Create(&newTesti).Error; err != nil {
		return err
	}
	return nil
}

func(tr *testimoniRepository) GetAll() ([]model.TestimoniDto, error) {
	var testimoni entities.Testimoni
	var allTestimoni []model.TestimoniDto
	if err := tr.db.Model(&testimoni).Find(&allTestimoni).Error; err != nil {
		return allTestimoni, err
	}
	return allTestimoni, nil
}

func (tr *testimoniRepository) GetByID(id uint) (entities.Testimoni, error) {
	var detailTestimoni entities.Testimoni

	if err := tr.db.First(&detailTestimoni, id).Error; err != nil {
		return detailTestimoni, err
	}
	return detailTestimoni, nil
}

func (tr *testimoniRepository) UpdateByID(updateTesti entities.Testimoni, id uint) error {
	var detailTestimoni entities.Testimoni

	if err := tr.db.Where("id = ?", id).Model(&detailTestimoni).Updates(updateTesti).Error; err != nil {
		return err
	}
	return nil
}

func (tr *testimoniRepository) DeleteByID(id uint) error {
	var testi entities.Testimoni
	if err := tr.db.Delete(&testi, id).Error; err != nil {
		return err
	}
	return nil
}
