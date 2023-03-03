package service

import (
	"errors"
	"rapsshop-project/entities"
	"rapsshop-project/middleware"
	"rapsshop-project/model"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AdminUsecase struct {
	AdminRepository model.AdminRepository
}

func NewAdminUsecase(repoAdmin model.AdminRepository) model.AdminUsecase {
	return &AdminUsecase{AdminRepository: repoAdmin}
}

func (a *AdminUsecase) Register(input *model.NewAdmin) error {
	admin, _ := a.AdminRepository.GetByUsername(input.Username)

	if len(admin.Username) > 0 {
		return errors.New("username is exist")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	newAdmin := entities.Admin{
		ID: uuid.NewString(),
		Username: input.Username,
		Password: string(hashedPassword),
		Nama: input.Nama,
	}

	if err := a.AdminRepository.Create(newAdmin); err != nil {
		return err
	}

	return nil
}

func(a *AdminUsecase) Login(input *model.AdminLogin) (string, error) {
	admin, err := a.AdminRepository.GetByUsername(input.Username)
	if err != nil {
		return "", err
	}

	if len(admin.Username) == 0 {
		return "", errors.New("username is not exist")
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password))
	if err != nil {
		return "", errors.New("password not match")
	}

	token, err := middleware.GenerateToken(admin.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *AdminUsecase) Profile(id string) (model.AdminDto, error) {
	admin, err := a.AdminRepository.GetByID(id)
	if err != nil {
		return admin, err
	}
	return admin, nil
}