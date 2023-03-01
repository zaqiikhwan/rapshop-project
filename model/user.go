package model

type UserRegister struct {
	Nama string `json:"nama" validate:"required"`
	Email string `json:"email"`
}