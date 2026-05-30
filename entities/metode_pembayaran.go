package entities

import "gorm.io/gorm"

type MetodePembayaran struct {
	gorm.Model
	IndexPembayaran      *int   `json:"index_pembayaran"`
	JenisPembayaran      string `json:"jenis_pembayaran"`
	KredensialPembayaran string `json:"kredensial_pembayaran"`
	Pemilik              string `json:"pemilik"`
}
