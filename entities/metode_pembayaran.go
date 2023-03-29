package entities

import "gorm.io/gorm"

type MetodePembayaran struct {
	gorm.Model
	IndexPembayaran int `json:"index_pembayaran"`
	JenisPembayaran string `json:"jenis_pembayaran"`
	KredensialPembayaran string `json:"kredensial_pembayaran"`
	Pemilik string `json:"pemilik"`
}

type InputMetodePembayaran struct {
	IndexPembayaran int `json:"index_pembayaran"`
	JenisPembayaran string `json:"jenis_pembayaran"`
	KredensialPembayaran string `json:"kredensial_pembayaran"`
	Pemilik string `json:"pemilik"`
}

type MetodePembayaranRepository interface {
	Create(newMetode MetodePembayaran) error 
	GetByIndex(index int) (MetodePembayaran, error)
	GetByID(id uint) (MetodePembayaran, error)
	GetAll() ([]MetodePembayaran, error)
	UpdateKredensialByID(id uint, patchKredensial MetodePembayaran) error
}

type MetodePembayaranUsecase interface {
	CreateNewPembayaran(input *InputMetodePembayaran) error
	GetAllPembayaran() ([]MetodePembayaran, error) 
	GetDetailPembayaranByIndex(index int) (MetodePembayaran, error)
	GetDetailPembayaranByID(id uint) (MetodePembayaran, error)
	PatchDetailPembayaranByID(id uint, input *InputMetodePembayaran) error
}