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
	GetByJenis(jenis string) (MetodePembayaran, error)
	GetByIndex(index int) (MetodePembayaran, error)
	GetAll() ([]MetodePembayaran, error)
	UpdateKredensial(jenis string, patchKredensial MetodePembayaran) error
}

type MetodePembayaranUsecase interface {
	CreateNewPembayaran(input *InputMetodePembayaran) error
	GetAllPembayaran() ([]MetodePembayaran, error) 
	GetDetailPembayaran(jenis string) (MetodePembayaran, error)
	GetDetailPembayaranByIndex(index int) (MetodePembayaran, error)
	PatchDetailPembayaran(jenis string, input *InputMetodePembayaran) error
}