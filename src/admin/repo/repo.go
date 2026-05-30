package repo

import (
	"rapsshop-project/entities"
	"rapsshop-project/model"

	"gorm.io/gorm"
)

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) model.AdminRepository {
	return &adminRepository{db: db}
}

func (ar *adminRepository) Create(newAdmin entities.Admin) error {
	if err := ar.db.Create(&newAdmin).Error; err != nil {
		return err
	}
	return nil
}

func (ar *adminRepository) GetByUsername(username string) (entities.Admin, error) {
	var admin entities.Admin
	err := ar.db.Where("username = ?", username).First(&admin).Error; if err != nil {
		return admin, err
	}
	return admin, nil	
}

func (ar *adminRepository) GetByID(id string) (model.AdminDto, error) {
	var admin entities.Admin
	var adminDto model.AdminDto
	err := ar.db.Model(&admin).Where("id = ?", id).First(&adminDto).Error; if err != nil {
		return adminDto, err
	}
	return adminDto, nil
}