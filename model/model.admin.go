package model

import (
	"rapsshop-project/entities"
)

type NewAdmin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nama     string `json:"nama"`
	Token    string `json:"token"`
}

type AdminLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminDto struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Nama      string `jsom:"nama"`
}

type AdminRepository interface {
	Create(admin entities.Admin) error
	GetByUsername(username string) (entities.Admin, error)
	GetByID(id string) (AdminDto, error)
}

type AdminUsecase interface {
	Register(input *NewAdmin) error
	Login(input *AdminLogin) (string, error)
	Profile(id string) (AdminDto, error)
}