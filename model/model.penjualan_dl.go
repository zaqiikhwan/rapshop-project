package model

import (
	"mime/multipart"
	"rapsshop-project/entities"
)

type InputPenjualanDL struct {
	JumlahDL         int    `json:"jumlah_dl"`
	JumlahTransaksi  int    `json:"jumlah_transaksi"`
	WA               string `json:"wa"`
	Transfer         string `json:"transfer"`
	NomorTransfer    string `json:"nomor_transfer"`
	StatusPembayaran int    `json:"status"`
	BuktiDL          string `json:"bukti_dl"`
}

type PenjualanDLRepository interface {
	Create(input entities.PenjualanDL) error
	GetAll() ([]entities.PenjualanDL, error)
	GetByID(id uint) (entities.PenjualanDL, error)
	UpdateByID(id uint, input entities.PenjualanDL) error
	DeleteByID(id uint) error
}

type PenjualanDLUsecase interface {
	Create(image *multipart.FileHeader, jumlahDL int, jumlahTransaksi int, wa string, transfer string, nomorTransfer string) error
	GetAll() ([]entities.PenjualanDL, error)
	GetByID(id uint) (entities.PenjualanDL, error)
	UpdateByID(id uint, input entities.PenjualanDL) (entities.PenjualanDL, error)
	DeleteByID(id uint) error
}

// type TestimoniUsecase interface {
// 	CreateTestimoni(image *multipart.FileHeader, testi string, uname string, title string) error
// 	GetAllTestimoni() ([]TestimoniDto, error) // need paginate implement later
// 	GetTestimoniByID(id uint) (entities.Testimoni, error)
// 	UpdateTestimoniByID(id uint, image *multipart.FileHeader, testi string, uname string, title string) (entities.Testimoni, error)
// 	DeleteTestimoniByID(id uint) error
// }