package models

import "time"

// Todo adalah representasi satu item tugas di database
// UserID menghubungkan todo ke pemiliknya (relasi foreign key)
type Todo struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTodoRequest adalah data yang dikirim client saat membuat todo baru
type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateTodoRequest adalah data untuk mengupdate todo yang sudah ada
// Pointer (*string, *bool) digunakan agar bisa membedakan "tidak dikirim" vs "dikirim sebagai kosong"
type UpdateTodoRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}
