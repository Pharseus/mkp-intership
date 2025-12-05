package handlers

import (
	"database/sql"
	"encoding/json"
	"mkp/config"
	"mkp/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GetSchedules handler untuk mendapatkan semua jadwal tayang
func GetSchedules(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := `
		SELECT 
			s.id, s.movie_id, s.studio_id, s.start_time, s.end_time, 
			s.price, s.status, s.created_at,
			m.title as movie_title,
			st.name as studio_name,
			c.name as cinema_name
		FROM schedules s
		LEFT JOIN movies m ON s.movie_id = m.id
		LEFT JOIN studios st ON s.studio_id = st.id
		LEFT JOIN cinemas c ON st.cinema_id = c.id
		ORDER BY s.start_time DESC
	`

	rows, err := config.DB.Query(query)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer rows.Close()

	schedules := []models.Schedule{}
	for rows.Next() {
		var schedule models.Schedule
		err := rows.Scan(
			&schedule.ID,
			&schedule.MovieID,
			&schedule.StudioID,
			&schedule.StartTime,
			&schedule.EndTime,
			&schedule.Price,
			&schedule.Status,
			&schedule.CreatedAt,
			&schedule.MovieTitle,
			&schedule.StudioName,
			&schedule.CinemaName,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error scanning data")
			return
		}
		schedules = append(schedules, schedule)
	}

	respondWithJSON(w, http.StatusOK, schedules)
}

// GetScheduleByID handler untuk mendapatkan jadwal tayang berdasarkan ID
func GetScheduleByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Ambil ID dari URL path
	id := extractIDFromPath(r.URL.Path, "/api/schedules/")
	if id == 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	query := `
		SELECT 
			s.id, s.movie_id, s.studio_id, s.start_time, s.end_time, 
			s.price, s.status, s.created_at,
			m.title as movie_title,
			st.name as studio_name,
			c.name as cinema_name
		FROM schedules s
		LEFT JOIN movies m ON s.movie_id = m.id
		LEFT JOIN studios st ON s.studio_id = st.id
		LEFT JOIN cinemas c ON st.cinema_id = c.id
		WHERE s.id = $1
	`

	var schedule models.Schedule
	err := config.DB.QueryRow(query, id).Scan(
		&schedule.ID,
		&schedule.MovieID,
		&schedule.StudioID,
		&schedule.StartTime,
		&schedule.EndTime,
		&schedule.Price,
		&schedule.Status,
		&schedule.CreatedAt,
		&schedule.MovieTitle,
		&schedule.StudioName,
		&schedule.CinemaName,
	)

	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "Schedule not found")
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	respondWithJSON(w, http.StatusOK, schedule)
}

// CreateSchedule handler untuk membuat jadwal tayang baru
func CreateSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.ScheduleCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validasi input
	if req.MovieID == 0 || req.StudioID == 0 || req.StartTime == "" || req.EndTime == "" || req.Price <= 0 {
		respondWithError(w, http.StatusBadRequest, "All fields are required and price must be positive")
		return
	}

	// Parse waktu
	startTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid start_time format. Use: YYYY-MM-DD HH:MM:SS")
		return
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid end_time format. Use: YYYY-MM-DD HH:MM:SS")
		return
	}

	// Set default status jika tidak ada
	if req.Status == "" {
		req.Status = "SHOWING"
	}

	// Insert ke database
	query := `
		INSERT INTO schedules (movie_id, studio_id, start_time, end_time, price, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var scheduleID int
	err = config.DB.QueryRow(
		query,
		req.MovieID,
		req.StudioID,
		startTime,
		endTime,
		req.Price,
		req.Status,
		time.Now(),
	).Scan(&scheduleID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create schedule")
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Schedule created successfully",
		"id":      scheduleID,
	})
}

// UpdateSchedule handler untuk mengupdate jadwal tayang
func UpdateSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Ambil ID dari URL path
	id := extractIDFromPath(r.URL.Path, "/api/schedules/")
	if id == 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	var req models.ScheduleUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Build dynamic update query
	updates := []string{}
	args := []interface{}{}
	argID := 1

	if req.MovieID != nil {
		updates = append(updates, "movie_id = $"+strconv.Itoa(argID))
		args = append(args, *req.MovieID)
		argID++
	}
	if req.StudioID != nil {
		updates = append(updates, "studio_id = $"+strconv.Itoa(argID))
		args = append(args, *req.StudioID)
		argID++
	}
	if req.StartTime != nil {
		startTime, err := time.Parse("2006-01-02 15:04:05", *req.StartTime)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid start_time format")
			return
		}
		updates = append(updates, "start_time = $"+strconv.Itoa(argID))
		args = append(args, startTime)
		argID++
	}
	if req.EndTime != nil {
		endTime, err := time.Parse("2006-01-02 15:04:05", *req.EndTime)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid end_time format")
			return
		}
		updates = append(updates, "end_time = $"+strconv.Itoa(argID))
		args = append(args, endTime)
		argID++
	}
	if req.Price != nil {
		updates = append(updates, "price = $"+strconv.Itoa(argID))
		args = append(args, *req.Price)
		argID++
	}
	if req.Status != nil {
		updates = append(updates, "status = $"+strconv.Itoa(argID))
		args = append(args, *req.Status)
		argID++
	}

	if len(updates) == 0 {
		respondWithError(w, http.StatusBadRequest, "No fields to update")
		return
	}

	// Add ID to args
	args = append(args, id)

	query := "UPDATE schedules SET " + strings.Join(updates, ", ") + " WHERE id = $" + strconv.Itoa(argID)

	result, err := config.DB.Exec(query, args...)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update schedule")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "Schedule not found")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Schedule updated successfully",
	})
}

// DeleteSchedule handler untuk menghapus jadwal tayang
func DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Ambil ID dari URL path
	id := extractIDFromPath(r.URL.Path, "/api/schedules/")
	if id == 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	query := "DELETE FROM schedules WHERE id = $1"
	result, err := config.DB.Exec(query, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete schedule")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "Schedule not found")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Schedule deleted successfully",
	})
}

// Helper function untuk extract ID dari URL path
func extractIDFromPath(path string, prefix string) int {
	idStr := strings.TrimPrefix(path, prefix)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0
	}
	return id
}
