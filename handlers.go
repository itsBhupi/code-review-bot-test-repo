package main

import (
	"fmt"
	"net/http"
)

// No middleware, no proper error handling
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// No method validation
	// No proper JSON response
	fmt.Fprintf(w, "OK")
}

// Poor error handling and security
func AdminPanel(w http.ResponseWriter, r *http.Request) {
	// No authentication check
	// No authorization
	// No CSRF protection

	password := r.FormValue("password")

	// Hardcoded admin password
	if password == "admin123" {
		fmt.Fprintf(w, "Welcome Admin! Database password: %s", database_url)
	} else {
		// No proper status code
		fmt.Fprintf(w, "Access denied")
	}
}
