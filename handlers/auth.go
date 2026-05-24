package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"todo-api/config"
	"todo-api/database"
	"todo-api/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// writeJSON adalah helper untuk mengirim response JSON dengan mudah
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Register menangani pendaftaran pengguna baru
// Method: POST /api/register
// Body: {"name":"...", "email":"...", "password":"..."}
func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	// Decode body JSON dari request ke struct RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON salah"})
		return
	}

	// Validasi field wajib diisi
	if req.Name == "" || req.Email == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Nama, email, dan password wajib diisi"})
		return
	}

	// Hash password menggunakan bcrypt sebelum disimpan ke database
	// JANGAN PERNAH simpan password plain text!
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal hash password"})
		return
	}

	// Buat objek User baru
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	// Simpan ke database
	// GORM akan mengembalikan error jika email sudah terdaftar (karena constraint unique)
	if err := database.DB.Create(&user).Error; err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "Email sudah terdaftar"})
		return
	}

	// Kirim response sukses (201 Created)
	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Registrasi berhasil",
		"user":    user,
	})
}

// Login menangani proses masuk pengguna
// Method: POST /api/login
// Body: {"email":"...", "password":"..."}
// Response: {"token":"..."} — token ini dipakai untuk request selanjutnya
func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON salah"})
		return
	}

	// Cari user berdasarkan email
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// Pesan error generik agar hacker tidak tahu apakah email terdaftar atau tidak
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Email atau password salah"})
		return
	}

	// Bandingkan password yang dimasukkan dengan hash di database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Email atau password salah"})
		return
	}

	cfg := config.GetConfig()

	// Buat JWT token dengan data (claims) yang akan disimpan di dalamnya
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Token expired dalam 24 jam
	})

	// Tandatangani token dengan secret key
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal membuat token"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Login berhasil",
		"token":   tokenString,
		"user":    user,
	})
}
