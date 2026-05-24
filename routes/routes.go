package routes

import (
	"net/http"
	"todo-api/handlers"
	"todo-api/middleware"

	"github.com/gorilla/mux"
)

// SetupRoutes mendaftarkan semua endpoint API ke router
// Dibagi menjadi dua kelompok:
// 1. Public routes — bisa diakses tanpa login
// 2. Protected routes — harus login (pakai JWT token)
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Middleware global: izinkan semua origin untuk CORS
	// Agar website testing di browser bisa memanggil API ini
	router.Use(corsMiddleware)

	// === PUBLIC ROUTES (tidak perlu token) ===
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/register", handlers.Register).Methods("POST")
	api.HandleFunc("/login", handlers.Login).Methods("POST")

	// === PROTECTED ROUTES (harus ada token JWT) ===
	// Semua route di sini akan melewati AuthMiddleware terlebih dahulu
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/todos", handlers.GetTodos).Methods("GET")           // Ambil semua todo
	protected.HandleFunc("/todos/{id}", handlers.GetTodo).Methods("GET")       // Ambil satu todo
	protected.HandleFunc("/todos", handlers.CreateTodo).Methods("POST")        // Buat todo baru
	protected.HandleFunc("/todos/{id}", handlers.UpdateTodo).Methods("PUT")    // Update todo
	protected.HandleFunc("/todos/{id}", handlers.DeleteTodo).Methods("DELETE") // Hapus todo

	// Sajikan file HTML untuk testing (website sederhana)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))

	return router
}

// corsMiddleware mengizinkan browser mengakses API dari domain mana pun
// Ini diperlukan agar file HTML bisa memanggil API di localhost
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight request dari browser
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
