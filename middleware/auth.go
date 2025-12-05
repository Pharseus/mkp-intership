package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Secret key untuk JWT - dalam production sebaiknya disimpan di environment variable
var JWTSecret = []byte("your-secret-key-change-this-in-production")

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// AuthMiddleware middleware untuk validasi JWT token
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondWithError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		tokenString := parts[1]

		// Parse dan validasi token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Validasi signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return JWTSecret, nil
		})

		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Ambil claims
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			// Simpan user info ke context untuk digunakan di handler
			ctx := context.WithValue(r.Context(), "userID", claims.UserID)
			ctx = context.WithValue(ctx, "email", claims.Email)

			// Lanjutkan ke handler berikutnya dengan context yang sudah berisi user info
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			respondWithError(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}
	}
}

// Helper function untuk mengirim error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"error": "%s"}`, message)
}
