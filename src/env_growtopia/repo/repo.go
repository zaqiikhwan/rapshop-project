package repo

import (
	"rapsshop-project/entities"
	"rapsshop-project/model"

	"gorm.io/gorm"
)

type envGrowtopiaRepository struct {
	db *gorm.DB
}

func NewEnvGrowtopiaRepo(db *gorm.DB) model.GrowtopiaEnvRepository {
	return &envGrowtopiaRepository{db: db}
}

func (egr *envGrowtopiaRepository) Create(newEnv entities.Growtopia) error {
	if err := egr.db.Create(&newEnv).Error; err != nil {
		return err
	}
	return nil
}

func (egr *envGrowtopiaRepository) GetLatest() (model.GrowtopiaEnvDto, error) {
	var detail model.GrowtopiaEnvDto
	if err := egr.db.Model(entities.Growtopia{}).Order("id desc").First(&detail).Error; err != nil {
		return detail, err
	}
	return detail, nil
}

func (egr *envGrowtopiaRepository) UpdateByID(updateEnvGrow entities.Growtopia, id uint) error {
	var updateEnv entities.Growtopia
	if err := egr.db.Model(&updateEnv).Where("id = ?", id).Updates(updateEnvGrow).Error; err != nil {
		return err
	}
	return nil
}