package model

import (
	"mime/multipart"
	"rapsshop-project/entities"
)

type InputPenjualanDL struct {
	Nama 			 string `json:"nama"`
	JumlahDL         int    `json:"jumlah_dl"`
	JumlahTransaksi  int    `json:"jumlah_transaksi"`
	WA               string `json:"wa"`
	Transfer         string `json:"transfer"`
	EditorStatus     string `json:"editor"`
	NomorTransfer    string `json:"nomor_transfer"`
	StatusPembayaran int    `json:"status"`
	BuktiDL          string `json:"bukti_dl"`
	HargaJual 		 int 	`json:"harga_jual"`
}

// fitur laporan
type RekapTransaksiPenjualan struct {
	JumlahTransaksi int `json:"jumlah_transaksi"`
	JumlahDL int `json:"jumlah_dl"`
	Rate int `json:"rate"`
}

type RekapTransaksiPembelian struct {
	JumlahTransaksi int `json:"jumlah_transaksi"`
	JumlahDL int `json:"jumlah_dl"`
	Rate int `json:"rate"`
}

type RekapTotalPenjualan struct {
	Tanggal string `json:"tanggal"`
	JumlahDL int `json:"jumlah_dl"`
}

type RekapProfit struct {
	Tanggal string `json:"tanggal"`
	Profit int `json:"profit"`
}

type Total struct {
	TransaksiJual int
	TransaksiBeli int
}

type PenjualanDLRepository interface {
	Create(input entities.PenjualanDL) error
	GetAll(_startInt int, _endInt int) ([]entities.PenjualanDL, int, error)
	GetByID(id uint) (entities.PenjualanDL, error)
	GetByDate(date string) ([]RekapTransaksiPenjualan, []RekapTransaksiPembelian, error)
	GetProfit(date string) ([]RekapProfit, error)
	GetTotalPenjualan(date string) ([]RekapTotalPenjualan, error)
	UpdateByID(id uint, input entities.PenjualanDL) error
	DeleteByID(id uint) error
}

type PenjualanDLUsecase interface {
	Create(image *multipart.FileHeader, jumlahDL int, jumlahTransaksi int, wa string, transfer string, nomorTransfer string, nama string, hargaJualDL int) error
	GetAll(_startInt int , _endInt int) ([]entities.PenjualanDL, int, error)
	GetByID(id uint) (entities.PenjualanDL, error)
	GetByDate(date string) ([]RekapTransaksiPenjualan, []RekapTransaksiPembelian, error)
	GetTotal(date string) ([]RekapTotalPenjualan, error)
	GetProfit(date string) ([]RekapProfit, error)
	UpdateByID(id uint, input entities.PenjualanDL) (entities.PenjualanDL, error)
	DeleteByID(id uint) error
}