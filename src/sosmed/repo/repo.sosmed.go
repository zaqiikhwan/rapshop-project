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
	if err := sr.db.Model(&models).Find(&model.SosmedDto{}).Error; err != nil {
		return []model.SosmedDto{}, err
	}
	return []model.SosmedDto{}, nil
}

func (sr *sosmedRepository) GetByID(id uint) (entities.Sosmed, error) {
	if err := sr.db.First(&entities.Sosmed{}, id).Error; err != nil {
		return entities.Sosmed{}, err
	}
	return entities.Sosmed{}, nil
}

func (sr *sosmedRepository) UpdateByID(updateSosmed entities.Sosmed, id uint) error {
	return nil
}

func (sr *sosmedRepository) DeleteByID(id uint) error {
	return nil
}