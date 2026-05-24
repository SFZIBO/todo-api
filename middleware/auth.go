package middleware

import (
	"context"
	"net/http"
	"strings"
	"todo-api/config"

	"github.com/golang-jwt/jwt/v5"
)

// contextKey adalah tipe khusus untuk key di context, mencegah konflik dengan package lain
type contextKey string

const UserIDKey contextKey = "userID"

// AuthMiddleware adalah fungsi yang berjalan SEBELUM handler utama
// Tugasnya: cek apakah request punya token JWT yang valid
// Jika valid → lanjut ke handler. Jika tidak → tolak dengan 401 Unauthorized
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ambil header "Authorization" dari request
		// Format yang diharapkan: "Bearer <token>"
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"Token tidak ditemukan"}`, http.StatusUnauthorized)
			return
		}

		// Pisahkan "Bearer" dan token-nya
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error":"Format token salah"}`, http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		cfg := config.GetConfig()

		// Parse dan verifikasi token JWT menggunakan secret key
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan algoritma yang digunakan adalah HMAC (bukan yang lain)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"error":"Token tidak valid atau sudah expired"}`, http.StatusUnauthorized)
			return
		}

		// Ambil claims (data yang disimpan di dalam token)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"error":"Token tidak bisa dibaca"}`, http.StatusUnauthorized)
			return
		}

		// Ambil user_id dari claims dan simpan ke context
		// Context digunakan untuk mengirim data dari middleware ke handler
		userID := uint(claims["user_id"].(float64))
		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		// Lanjutkan ke handler berikutnya dengan context yang sudah diisi userID
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID adalah helper untuk mengambil userID dari context di dalam handler
func GetUserID(r *http.Request) uint {
	return r.Context().Value(UserIDKey).(uint)
}
