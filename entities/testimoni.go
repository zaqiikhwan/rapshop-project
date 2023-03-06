package entities

import "gorm.io/gorm"

type Testimoni struct {
	gorm.Model
	Gambar string `json:"gambar"`
	Testimoni string `json:"testi"`
	Username string `json:"uname"`
	Title string `json:"title"`
}