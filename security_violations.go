package main

import (
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

// Database connection variable with a different name to avoid conflict
var dbConn *gorm.DB

// SQL Injection vulnerability: direct string concatenation
func getUserByID(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	// No input validation

	// VIOLATION: Direct string concatenation in SQL query
	query := "SELECT name, email FROM users WHERE id = " + userID
	var user struct {
		Name  string
		Email string
	}
	dbConn.Raw(query).Scan(&user)

	// VIOLATION: Storing credentials in plaintext
	const apiKey = "sk_live_51AbCdEfGhIjKlMnOpQrStUvWxYz1234567890AbCdEfGhIjKl"

	// VIOLATION: Improper HTML escaping
	output := "<div>User: " + user.Name + "</div>"
	fmt.Fprintf(w, output)
}
