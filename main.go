package main

import (
	"log"
	"mkp/config"
	"mkp/handlers"
	"mkp/middleware"
	"net/http"
)

func main() {
	// Inisialisasi database
	config.InitDB()
	defer config.CloseDB()

	// Setup routes
	setupRoutes()

	// Start server
	port := ":8080"
	log.Printf("Server running on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func setupRoutes() {
	// Public routes (tidak perlu authentication)
	http.HandleFunc("/api/register", handlers.Register)
	http.HandleFunc("/api/login", handlers.Login)

	// Protected routes (perlu authentication)
	// GET semua jadwal
	http.HandleFunc("/api/schedules", middleware.AuthMiddleware(handlers.GetSchedules))

	// POST create jadwal baru
	http.HandleFunc("/api/schedules/create", middleware.AuthMiddleware(handlers.CreateSchedule))

	// GET, PUT, DELETE jadwal by ID - menggunakan pattern yang sama
	http.HandleFunc("/api/schedules/", func(w http.ResponseWriter, r *http.Request) {
		// Pastikan ada ID di path
		if r.URL.Path == "/api/schedules/" || r.URL.Path == "/api/schedules" {
			middleware.AuthMiddleware(handlers.GetSchedules)(w, r)
			return
		}

		// Route berdasarkan HTTP method
		switch r.Method {
		case http.MethodGet:
			middleware.AuthMiddleware(handlers.GetScheduleByID)(w, r)
		case http.MethodPut:
			middleware.AuthMiddleware(handlers.UpdateSchedule)(w, r)
		case http.MethodDelete:
			middleware.AuthMiddleware(handlers.DeleteSchedule)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Root endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Welcome to MKP Cinema API", "version": "1.0"}`))
	})
}
