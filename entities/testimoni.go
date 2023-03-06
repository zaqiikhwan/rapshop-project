package entities

import "gorm.io/gorm"

type Testimoni struct {
	gorm.Model
	Gambar string  
	Testimoni string 
	Username string
	Title string 
}