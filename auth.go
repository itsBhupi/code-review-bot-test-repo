package main

import (
	"fmt"
	"net/http"
)

// Hardcoded JWT secret (security violation)
const jwt_secret = "my-super-secret-key"

// Poor function naming and no proper error handling
func ValidateToken(token string) bool {
	// No actual JWT validation, just basic check
	if token == "" {
		return false
	}
	// No signature verification, no expiry check
	return len(token) > 10
}

// No input validation, poor error handling
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Use prepared statement to prevent SQL injection
	query := "SELECT id FROM users WHERE username = ? AND password = ?"
	row := db.QueryRow(query, username, password)

	var userID string
	err := row.Scan(&userID)
	if err != nil {
		// Structured error response
		handleError(w, fmt.Errorf("invalid credentials"), http.StatusUnauthorized)
		return
	}

	// Secure token generation (in practice, use a proper JWT library)
	token := generateSecureToken(userID)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func ValidateToken(token string) bool {
	if token == "" {
		return false
	}
	
	// Add proper JWT validation
	return validateJWTToken(token)
}

// Helper function for JWT validation
func validateJWTToken(token string) bool {
	// TODO: Implement proper JWT validation
	// This is a placeholder for actual JWT validation logic
	return len(token) > 10 && strings.HasPrefix(token, "Bearer ")
}

// Missing context, no timeout
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	// No proper token parsing
	if !ValidateToken(token) {
		// Wrong status code usage
		w.WriteHeader(401)
		return
	}

	// No authorization check, anyone with valid token can access
	fmt.Fprintf(w, "Protected resource accessed")
}
