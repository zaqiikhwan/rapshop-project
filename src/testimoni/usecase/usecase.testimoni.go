package usecase

import (
	"errors"
	"mime/multipart"
	"os"
	"rapsshop-project/entities"
	"rapsshop-project/model"

	storage_go "github.com/supabase-community/storage-go"
)

type testimoniUsecase struct {
	TestimoniRepository model.TestimoniRepository
}

func NewTestimoniUsecase(repoTesti model.TestimoniRepository) model.TestimoniUsecase {
	return &testimoniUsecase{TestimoniRepository: repoTesti}
}

func (tu *testimoniUsecase) CreateTestimoni(image *multipart.FileHeader, testi string, jumlah int) error {

	client := storage_go.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SERVICE_TOKEN"), nil)

	if client == nil {
		return errors.New("storage authentication failed")
	}

	imageIo, err := image.Open()

	client.UploadFile(os.Getenv("STORAGE_NAME"), image.Filename, imageIo)

	if err != nil {
		return err
	}
	newTesti := entities.Testimoni{
		Gambar: os.Getenv("BASE_URL") + image.Filename,
		Testimoni: testi,
		JumlahDL: jumlah,
	}

	if err := tu.TestimoniRepository.Create(newTesti); err != nil {
		return err
	}
	return nil
}

func (tu *testimoniUsecase) GetAllTestimoni() ([]model.TestimoniDto, error) {
	allTesti, err := tu.TestimoniRepository.GetAll()
	if err != nil {
		return allTesti, err
	}
	return allTesti, nil
}

func (tu *testimoniUsecase) GetTestimoniByID(id uint) (entities.Testimoni, error) {
	detailTestimoni, err := tu.TestimoniRepository.GetByID(id)

	if err != nil {
		return detailTestimoni, err
	}

	return detailTestimoni, nil
}

func (tu *testimoniUsecase) UpdateTestimoniByID(id uint, image *multipart.FileHeader, testi string, jumlah int) (entities.Testimoni, error) {
	client := storage_go.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SERVICE_TOKEN"), nil)

	if client == nil {
		return entities.Testimoni{}, errors.New("storage authentication failed")
	}

	detailTestimoni, err := tu.TestimoniRepository.GetByID(id)

	if err != nil {
		return detailTestimoni, err
	}

	var updateTesti entities.Testimoni

	if image != nil {
		paths := make([]string, 1) 
		paths = append(paths, detailTestimoni.Gambar)
	
		client.RemoveFile(os.Getenv("STORAGE_NAME"), paths)
		imageIo, err := image.Open()
	
		if err != nil {
			return entities.Testimoni{}, err
		}
	
		client.UploadFile(os.Getenv("STORAGE_NAME"), image.Filename, imageIo)

		updateTesti = entities.Testimoni{
			Testimoni: testi,
			JumlahDL: jumlah,
			Gambar: os.Getenv("BASE_URL") + image.Filename,
		}

		err = tu.TestimoniRepository.UpdateByID(updateTesti, id)

		if err != nil {
			return entities.Testimoni{}, err
		}
	} else {
		updateTesti = entities.Testimoni{
			Testimoni: testi,
			JumlahDL: jumlah,
		}
	
		err = tu.TestimoniRepository.UpdateByID(updateTesti, id)
	
		if err != nil {
			return entities.Testimoni{}, err
		}
	}

	detailTestimoni, err = tu.TestimoniRepository.GetByID(id)

	if err != nil {
		return entities.Testimoni{}, err
	}

	return detailTestimoni, nil
}

func (tu *testimoniUsecase) DeleteTestimoniByID(id uint) error {
	if err := tu.TestimoniRepository.DeleteByID(id); err != nil {
		return err
	}
	return nil
}