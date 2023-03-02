package entities

// dl dibeli user
type PembelianDL struct {
	ID string `gorm:"primaryKey"`
	World string `gorm:"size:255"`
	Nama string `gorm:"size:255"`
	GrowID string `gorm:"size:255"`
	JumlahDL int
	WA string `gorm:"size:20"`
	Transfer string `gorm:"size:255"`
	JumlahTransaksi int
	StatusPengiriman bool
}