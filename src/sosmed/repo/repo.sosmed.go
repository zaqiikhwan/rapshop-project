package repo

import (
	"rapsshop-project/entities"
	"rapsshop-project/model"

	"gorm.io/gorm"
)

type sosmedRepository struct {
	db *gorm.DB
}

func NewSosmedRepository(db *gorm.DB) model.SosmedRepository {
	return &sosmedRepository{db: db}
}

func (sr *sosmedRepository) Create(newSosmed entities.Sosmed) error {
	if err := sr.db.Create(&newSosmed).Error; err != nil {
		return err
	}
	return nil
}

func (sr *sosmedRepository) GetAll() ([]model.SosmedDto, error) {
	var models entities.Sosmed
	var allPlatform []model.SosmedDto
	if err := sr.db.Model(&models).Find(&allPlatform).Error; err != nil {
		return []model.SosmedDto{}, err
	}
	return allPlatform, nil
}

func (sr *sosmedRepository) GetByID(id uint) (entities.Sosmed, error) {
	var detail entities.Sosmed
	if err := sr.db.First(&detail, id).Error; err != nil {
		return entities.Sosmed{}, err
	}
	return detail, nil
}

func (sr *sosmedRepository) UpdateByID(updateSosmed entities.Sosmed, id uint) error {
	var detailSosmed entities.Testimoni

	if err := sr.db.Where("id = ?", id).Model(&detailSosmed).Updates(updateSosmed).Error; err != nil {
		return err
	}
	return nil
}

func (sr *sosmedRepository) DeleteByID(id uint) error {
	var detailSosmed entities.Testimoni
	if err := sr.db.Delete(&detailSosmed, id).Error; err != nil {
		return err
	}

	return nil
}