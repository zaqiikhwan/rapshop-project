package model

import "rapsshop-project/entities"

type InputGrowtopiaEnv struct {
	World    string `json:"world"`
	Password string `json:"password"`
	Owner    string `json:"owner"`
}

type GrowtopiaEnvDto struct {
	ID       uint   `json:"id"`
	World    string `json:"world"`
	Password string `json:"password"`
	Owner    string `json:"owner"`
}
type GrowtopiaEnvRepository interface {
	Create(newHarga entities.Growtopia) error
	GetLatest() (GrowtopiaEnvDto, error)
	UpdateByID(updateEnv entities.Growtopia, id uint) error
}

type GrowtopiaEnvUsecase interface {
	CreateNewEnv(input *InputGrowtopiaEnv) error
	GetLatestEnv() (GrowtopiaEnvDto, error)
	UpdateLatestEnv(input *InputGrowtopiaEnv) (GrowtopiaEnvDto, error)
}