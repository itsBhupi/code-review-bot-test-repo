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
	username := r.FormValue("username")
	password := r.FormValue("password")

	// SQL injection vulnerability
	query := fmt.Sprintf("SELECT id FROM users WHERE username='%s' AND password='%s'", username, password)

	// Using global db without proper connection handling
	rows, _ := db.Query(query) // Ignoring error
	defer rows.Close()

	if rows.Next() {
		// Logging sensitive information
		fmt.Printf("User %s logged in with password %s", username, password)

		// No proper session management
		w.Header().Set("Authorization", "Bearer fake-token-123")
		fmt.Fprintf(w, "Login successful")
	} else {
		// Poor error response, no status code
		fmt.Fprintf(w, "Invalid credentials")
	}
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
