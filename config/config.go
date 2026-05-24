package config

// Config menyimpan semua pengaturan aplikasi
// Dalam proyek nyata, nilai ini sebaiknya dibaca dari environment variable
type Config struct {
	Port      string // Port server berjalan
	JWTSecret string // Kunci rahasia untuk membuat dan memverifikasi token JWT
	DBPath    string // Lokasi file database SQLite
}

// GetConfig mengembalikan konfigurasi default
func GetConfig() Config {
	return Config{
		Port:      "8080",
		JWTSecret: "rahasia-jwt-super-aman-ganti-ini-di-produksi",
		DBPath:    "todo.db",
	}
}
