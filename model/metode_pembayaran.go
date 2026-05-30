package model

import "rapsshop-project/entities"

type InputMetodePembayaran struct {
	IndexPembayaran      int    `json:"index_pembayaran"`
	JenisPembayaran      string `json:"jenis_pembayaran"`
	KredensialPembayaran string `json:"kredensial_pembayaran"`
	Pemilik              string `json:"pemilik"`
}

type MetodePembayaranRepository interface {
	Create(newMetode entities.MetodePembayaran) error
	GetAll() ([]entities.MetodePembayaran, error)
	GetByID(id uint) (entities.MetodePembayaran, error)
	GetByIndex(index int) (entities.MetodePembayaran, error)
	UpdateKredensialByID(id uint, patchKredensial entities.MetodePembayaran) error
	DeleteByID(id uint) error
}

type MetodePembayaranUsecase interface {
	CreateNewPembayaran(input *InputMetodePembayaran) error
	GetAllPembayaran() ([]entities.MetodePembayaran, error)
	GetDetailPembayaranByIndex(index int) (entities.MetodePembayaran, error)
	GetDetailPembayaranByID(id uint) (entities.MetodePembayaran, error)
	PatchDetailPembayaranByID(id uint, input *InputMetodePembayaran) error
	DeletePembayaranByID(id uint) error
}
