package model

import "rapsshop-project/entities"

type InputSosmed struct {
	Username string `json:"username"`
	Platform string `json:"platform"`
	Link     string `json:"link"`
}

type SosmedDto struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Platform string `json:"platform"`
	Link     string `json:"link"`
}

type SosmedRepository interface {
	Create(newSosmed entities.Sosmed) error
	GetAll() ([]SosmedDto, error) 
	GetByID(id uint) (entities.Sosmed, error)
	UpdateByID(updateSosmed entities.Sosmed, id uint) error
	DeleteByID(id uint) error
}

type SosmedUsecase interface {
	CreateSosmed(input *InputSosmed) error
	GetAllSosmed() ([]SosmedDto, error) // need paginate implement later
	GetSosmedByID(id uint) (entities.Sosmed, error)
	UpdateSosmedByID(id uint, input *InputSosmed) (entities.Sosmed, error)
	DeleteSosmedByID(id uint) error
}