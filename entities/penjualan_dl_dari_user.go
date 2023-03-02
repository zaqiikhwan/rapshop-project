package entities

// user menjual DL dan dibeli oleh admin rapshop
type PenjualanDL struct {
	ID string `gorm:"primaryKey"`
	JumlahDL int
	JumlahTransaksi int
	WA string
	Transfer string
	NomorTransfer string
	NamaEWallet string
	StatusPembayaran int
	BuktiDL string
}