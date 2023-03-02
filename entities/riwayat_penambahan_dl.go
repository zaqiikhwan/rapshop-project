package entities

import "gorm.io/gorm"

type RiwayatPenambahanDL struct {
	gorm.Model
	JumlahDL int 
	Harga int
}