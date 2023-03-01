package repository

import (
	"rapsshop-project/entities"
	"rapsshop-project/model"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (ur *userRepository) CreateUser(user *model.UserRegister) error {
	var userReg entities.User

	err := ur.db.Model(&userReg).Create(user).Error

	return err
}