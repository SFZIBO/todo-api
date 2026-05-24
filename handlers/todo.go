package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"todo-api/database"
	"todo-api/middleware"
	"todo-api/models"

	"github.com/gorilla/mux"
)

// GetTodos mengambil semua todo milik user yang sedang login
// Method: GET /api/todos
// Header: Authorization: Bearer <token>
func GetTodos(w http.ResponseWriter, r *http.Request) {
	// Ambil userID dari context (sudah diset oleh AuthMiddleware)
	userID := middleware.GetUserID(r)

	var todos []models.Todo

	// Query: SELECT * FROM todos WHERE user_id = ? ORDER BY created_at DESC
	database.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&todos)

	// Jika tidak ada todo, kembalikan array kosong (bukan null)
	if todos == nil {
		todos = []models.Todo{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"todos": todos,
		"total": len(todos),
	})
}

// GetTodo mengambil satu todo berdasarkan ID
// Method: GET /api/todos/{id}
func GetTodo(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	// Ambil parameter {id} dari URL
	vars := mux.Vars(r)
	todoID, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID tidak valid"})
		return
	}

	var todo models.Todo

	// Cari todo dengan id yang diminta DAN milik user yang login
	// Ini penting agar user tidak bisa akses todo milik orang lain!
	if err := database.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Todo tidak ditemukan"})
		return
	}

	writeJSON(w, http.StatusOK, todo)
}

// CreateTodo membuat todo baru
// Method: POST /api/todos
// Body: {"title":"...", "description":"..."}
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req models.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON salah"})
		return
	}

	if req.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Title wajib diisi"})
		return
	}

	todo := models.Todo{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Completed:   false, // Default: belum selesai
	}

	if err := database.DB.Create(&todo).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal membuat todo"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Todo berhasil dibuat",
		"todo":    todo,
	})
}

// UpdateTodo mengupdate todo yang sudah ada
// Method: PUT /api/todos/{id}
// Body: {"title":"...", "description":"...", "completed": true/false}
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	vars := mux.Vars(r)
	todoID, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID tidak valid"})
		return
	}

	// Pastikan todo ada dan milik user yang login
	var todo models.Todo
	if err := database.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Todo tidak ditemukan"})
		return
	}

	var req models.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON salah"})
		return
	}

	// Update hanya field yang dikirim (pakai pointer untuk cek nil)
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}

	// Simpan perubahan ke database
	database.DB.Save(&todo)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Todo berhasil diupdate",
		"todo":    todo,
	})
}

// DeleteTodo menghapus todo berdasarkan ID
// Method: DELETE /api/todos/{id}
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	vars := mux.Vars(r)
	todoID, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID tidak valid"})
		return
	}

	// Pastikan todo ada dan milik user yang login sebelum dihapus
	var todo models.Todo
	if err := database.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Todo tidak ditemukan"})
		return
	}

	database.DB.Delete(&todo)

	writeJSON(w, http.StatusOK, map[string]string{"message": "Todo berhasil dihapus"})
}
