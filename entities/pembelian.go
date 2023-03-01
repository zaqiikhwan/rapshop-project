package entities

type OrderDL struct {
	ID string `gorm:"primaryKey"`
	Nama string `gorm:"size:255"`
	World string `gorm:"size:255"`
	GrowID string `gorm:"size:255"`
	JumlahDL int
	WA string `gorm:"size:20"`
	Transfer string `gorm:"size:255"`
	JumlahTransaksi int
	StatusPengiriman bool
}