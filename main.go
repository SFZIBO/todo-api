package main

import (
	"fmt"
	"log"
	"net/http"
	"todo-api/config"
	"todo-api/database"
	"todo-api/routes"
)

func main() {
	// 1. Load konfigurasi
	cfg := config.GetConfig()

	// 2. Koneksikan ke database
	database.Connect(cfg.DBPath)

	// 3. Setup semua route
	router := routes.SetupRoutes()

	// 4. Jalankan server
	addr := ":" + cfg.Port
	fmt.Printf("🚀 Server berjalan di http://localhost%s\n", addr)
	fmt.Printf("🌐 Buka browser: http://localhost%s\n", addr)

	// ListenAndServe akan memblokir program dan terus mendengarkan request
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("Server gagal start:", err)
	}
}
