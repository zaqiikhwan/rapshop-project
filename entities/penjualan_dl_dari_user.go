package entities

import "time"

// user menjual DL dan dibeli oleh admin rapshop
type PenjualanDL struct {
	ID uint `gorm:"primaryKey" json:"id"`
	JumlahDL int `json:"jumlah_dl"`
	JumlahTransaksi int `json:"jumlah_transaksi"`
	WA string `json:"wa"`
	Transfer string `json:"transfer"`
	NomorTransfer string `json:"nomor_transfer"`
	StatusPembayaran *int `gorm:"default:0" json:"status"`
	BuktiDL string `json:"bukti_dl"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}