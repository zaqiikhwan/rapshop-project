package usecase

import (
	"rapsshop-project/entities"
	"rapsshop-project/model"
)

type sosmedUsecase struct {
	TestimoniRepository model.SosmedRepository
}

func NewSosmedUsecase(repoSosmed model.SosmedRepository) model.SosmedUsecase {
	return &sosmedUsecase{TestimoniRepository: repoSosmed}
}

func (su *sosmedUsecase) CreateSosmed(input *model.InputSosmed) error {
	newSosmed := entities.Sosmed{
		Username: input.Username,
		Platform: input.Platform,
		Link: input.Link,
	}
	if err := su.TestimoniRepository.Create(newSosmed); err != nil {
		return err
	}

	return nil
}

func (su *sosmedUsecase) GetAllSosmed() ([]model.SosmedDto, error) {
	var allSosmed []model.SosmedDto
	allSosmed, err := su.TestimoniRepository.GetAll()
	if err != nil {
		return allSosmed, err
	}
	return allSosmed, nil
}

func (su *sosmedUsecase) GetSosmedByID(id uint) (entities.Sosmed, error) {
	var detailSosmed entities.Sosmed

	detailSosmed, err := su.TestimoniRepository.GetByID(id)

	if err != nil {
		return detailSosmed, err
	}
	return detailSosmed, nil
}

func (su *sosmedUsecase) UpdateSosmedByID(id uint, input *model.InputSosmed) (entities.Sosmed, error) {
	return entities.Sosmed{}, nil
} 

func (su *sosmedUsecase) DeleteSosmedByID(id uint) error {
	return nil
}


