package model

import (
	"mime/multipart"
	"rapsshop-project/entities"
)

type TestimoniDto struct {
	ID        uint   `json:"id"`
	Gambar    string `json:"gambar"`
	Testimoni string `json:"testimoni"`
	JumlahDL  int    `json:"jumlah_dl"`
}

type InputTestimoni struct {
	Gambar    string `json:"gambar"`
	Testimoni string `json:"testimoni"`
	JumlahDL  int    `json:"jumlah_dl"`
}

type TestimoniRepository interface {
	Create(newTesti entities.Testimoni) error
	GetAll() ([]TestimoniDto, error) // need paginate in here
	GetByID(id uint) (entities.Testimoni, error)
	UpdateByID(updateTesti entities.Testimoni, id uint) error
	DeleteByID(id uint) error
}

type TestimoniUsecase interface {
	CreateTestimoni(image *multipart.FileHeader, testi string, jumlah int) error
	GetAllTestimoni() ([]TestimoniDto, error) // need paginate implement later
	GetTestimoniByID(id uint) (entities.Testimoni, error)
	UpdateTestimoniByID(id uint, image *multipart.FileHeader, testi string, jumlah int) (entities.Testimoni, error)
	DeleteTestimoniByID(id uint) error
}