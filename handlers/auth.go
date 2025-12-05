package handlers

import (
	"database/sql"
	"encoding/json"
	"mkp/config"
	"mkp/middleware"
	"mkp/models"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Register handler untuk registrasi user baru
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var input struct {
		Fullname string `json:"fullname"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validasi input
	if input.Fullname == "" {
		respondWithError(w, http.StatusBadRequest, "Fullname is required")
		return
	}
	if input.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required")
		return
	}
	if input.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Password is required")
		return
	}
	if len(input.Password) < 6 {
		respondWithError(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	// Cek apakah email sudah terdaftar
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err = config.DB.QueryRow(checkQuery, input.Email).Scan(&exists)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if exists {
		respondWithError(w, http.StatusConflict, "Email already registered")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	// Insert user baru ke database
	var user models.User
	query := `
		INSERT INTO users (fullname, email, password_hash, created_at, updated_at) 
		VALUES ($1, $2, $3, NOW(), NOW()) 
		RETURNING id, fullname, email, created_at, updated_at
	`
	err = config.DB.QueryRow(query, input.Fullname, input.Email, string(hashedPassword)).
		Scan(&user.ID, &user.Fullname, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	// Generate JWT token
	token, err := generateJWT(user.ID, user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := models.LoginResponse{
		Token: token,
		User:  user,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

// Login handler untuk autentikasi user
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	var user models.User
	query := "SELECT id, fullname, email, password_hash, created_at, updated_at FROM users WHERE email = $1"
	err := config.DB.QueryRow(query, req.Email).Scan(
		&user.ID,
		&user.Fullname,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := generateJWT(user.ID, user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := models.LoginResponse{
		Token: token,
		User:  user,
	}

	respondWithJSON(w, http.StatusOK, response)
}

// generateJWT membuat JWT token untuk user
func generateJWT(userID int, email string) (string, error) {
	// Set expiration time (24 jam)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Buat claims
	claims := &middleware.Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Buat token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(middleware.JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Helper function untuk mengirim JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Helper function untuk mengirim error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
