package database

import (
	"log"
	"todo-api/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB adalah variabel global yang menyimpan koneksi ke database
// Dengan cara ini, semua handler bisa mengakses DB yang sama
var DB *gorm.DB

// Connect membuka koneksi ke database SQLite dan membuat tabel otomatis
func Connect(dbPath string) {
	var err error

	// Buka file database SQLite (dibuat otomatis jika belum ada)
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		// Jika koneksi gagal, hentikan program
		log.Fatal("Gagal koneksi ke database:", err)
	}

	// AutoMigrate: buat/update tabel secara otomatis berdasarkan struct model
	// Tidak perlu tulis SQL CREATE TABLE secara manual
	err = DB.AutoMigrate(&models.User{}, &models.Todo{})
	if err != nil {
		log.Fatal("Gagal migrasi database:", err)
	}

	log.Println("✅ Database terhubung:", dbPath)
}
