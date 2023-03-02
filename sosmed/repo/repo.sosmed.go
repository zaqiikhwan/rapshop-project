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
	return nil
}

func (sr *sosmedRepository) GetAll() ([]model.SosmedDto, error) {
	return []model.SosmedDto{}, nil
}

func (sr *sosmedRepository) GetByID(id uint) (entities.Sosmed, error) {
	return entities.Sosmed{}, nil
}

func (sr *sosmedRepository) UpdateByID(updateSosmed entities.Sosmed, id uint) error {
	return nil
}

func (sr *sosmedRepository) DeleteByID(id uint) error {
	return nil
}