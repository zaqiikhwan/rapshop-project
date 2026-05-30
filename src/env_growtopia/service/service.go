package service

import (
	"rapsshop-project/entities"
	"rapsshop-project/model"
)

type envGrowtopiaUsecase struct {
	EnvGrowtopiaRepo model.GrowtopiaEnvRepository
}

func NewEnvGrowtopiaUsecase(envGrowRepo model.GrowtopiaEnvRepository) model.GrowtopiaEnvUsecase {
	return &envGrowtopiaUsecase{EnvGrowtopiaRepo: envGrowRepo}
}

func (egu *envGrowtopiaUsecase) CreateNewEnv(input *model.InputGrowtopiaEnv) error {
	newEnv := entities.Growtopia {
		World: input.World,
		Password: input.Password,
		Owner: input.Owner,
	}
	if err := egu.EnvGrowtopiaRepo.Create(newEnv); err != nil {
		return err
	}
	return nil
}

func (egu *envGrowtopiaUsecase) GetLatestEnv() (model.GrowtopiaEnvDto, error) {
	detail, err := egu.EnvGrowtopiaRepo.GetLatest()
	if err != nil {
		return detail, err
	}
	return detail, nil
}

func (egu *envGrowtopiaUsecase) UpdateLatestEnv(input *model.InputGrowtopiaEnv) (model.GrowtopiaEnvDto, error) {
	detail, err := egu.EnvGrowtopiaRepo.GetLatest()
	if err != nil {
		return detail, err
	}
	updateEnv := entities.Growtopia {
		World: input.World,
		Password: input.Password,
		Owner: input.Owner,
	}

	err = egu.EnvGrowtopiaRepo.UpdateByID(updateEnv, detail.ID)
	if err != nil {
		return detail, err
	}

	updatedEnv, err := egu.EnvGrowtopiaRepo.GetLatest()
	if err != nil {
		return updatedEnv, err
	}

	return updatedEnv, nil
}
