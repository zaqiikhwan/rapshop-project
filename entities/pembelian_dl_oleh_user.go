package entities

import "time"

// dl dibeli user
type PembelianDL struct {
	ID string `gorm:"primaryKey"`
	World string `gorm:"size:255" json:"world"`
	Nama string `gorm:"size:255" json:"nama"`
	GrowID string `gorm:"size:255" json:"grow_id"`
	JenisItem bool `json:"jenis_item"` // if i == 0 -> maka dl, else maka bgl 
	JumlahDL int `json:"jumlah_dl"` 
	WA string `gorm:"size:20" json:"wa"`
	MetodeTransfer int `json:"metode_transfer"`
	JumlahTransaksi int64 `json:"jumlah_transaksi"`
	ButtonBayar *bool `gorm:"default:false" json:"button_bayar"`
	StatusPembayaran string `gorm:"default:belum_dibayar" json:"status_pembayaran"`
	StatusPengiriman *bool `gorm:"default:false" json:"status_pengiriman"`
	EditorStatus     string `json:"editor"`
	HargaBeli int `json:"harga_beli"`
	BuktiPembayaran string `json:"bukti_pembayaran"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

