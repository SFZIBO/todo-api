package models

import "time"

// User adalah representasi data pengguna di database
// Tag `json:"-"` berarti field Password TIDAK akan dikirim ke client saat response JSON
// Tag `gorm:"unique"` berarti Email harus unik di database
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

// RegisterRequest adalah struct untuk menerima data saat registrasi
// `binding:"required"` artinya field wajib diisi
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginRequest adalah struct untuk menerima data saat login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
