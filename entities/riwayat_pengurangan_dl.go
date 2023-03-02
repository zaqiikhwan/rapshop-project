package entities

import "gorm.io/gorm"

type RiwayatPenguranganDL struct {
	gorm.Model
	JumlahDL int
	Harga    int
}