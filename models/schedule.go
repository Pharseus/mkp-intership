package models

import "time"

type Schedule struct {
	ID         int       `json:"id"`
	MovieID    int       `json:"movie_id"`
	StudioID   int       `json:"studio_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Price      float64   `json:"price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	MovieTitle string    `json:"movie_title,omitempty"`
	StudioName string    `json:"studio_name,omitempty"`
	CinemaName string    `json:"cinema_name,omitempty"`
}

// ScheduleCreateRequest model untuk membuat jadwal baru
type ScheduleCreateRequest struct {
	MovieID   int     `json:"movie_id"`
	StudioID  int     `json:"studio_id"`
	StartTime string  `json:"start_time"`
	EndTime   string  `json:"end_time"`
	Price     float64 `json:"price"`
	Status    string  `json:"status"`
}

// ScheduleUpdateRequest model untuk update jadwal
type ScheduleUpdateRequest struct {
	MovieID   *int     `json:"movie_id,omitempty"`
	StudioID  *int     `json:"studio_id,omitempty"`
	StartTime *string  `json:"start_time,omitempty"`
	EndTime   *string  `json:"end_time,omitempty"`
	Price     *float64 `json:"price,omitempty"`
	Status    *string  `json:"status,omitempty"`
}
