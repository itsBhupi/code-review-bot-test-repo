package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Bad struct naming and no validation tags
type user struct {
	id       int    // unexported fields in API struct
	Name     string // no validation
	email    string
	Password string // password in response struct
}

// No proper HTTP method handling
func UserAPI(w http.ResponseWriter, r *http.Request) {
	// No method validation - handles all methods the same way
	// No content-type validation
	// No rate limiting

	switch r.Method {
	case "GET":
		// No pagination, returns all users
		users := GetAllUsers() // Function from database.go with issues

		// No proper JSON encoding error handling
		json.NewEncoder(w).Encode(users)

	case "POST":
		// Reading entire body without size limit
		body, _ := ioutil.ReadAll(r.Body) // Ignoring error, deprecated function

		var newUser user
		// No JSON unmarshaling error handling
		json.Unmarshal(body, &newUser)

		// No input validation whatsoever
		// Returning 200 instead of 201 for creation
		w.WriteHeader(200)
		fmt.Fprintf(w, "User created: %+v", newUser)

	default:
		// Wrong status code for method not allowed
		w.WriteHeader(400)
		fmt.Fprintf(w, "Bad method")
	}
}

// No proper RESTful design
func DeleteEverything(w http.ResponseWriter, r *http.Request) {
	// Dangerous operation with no authentication
	// No confirmation required
	// Using GET for destructive operation

	if r.Method == "GET" { // Should be DELETE
		// No error handling
		connection.Exec("DELETE FROM users")
		connection.Exec("DELETE FROM accounts")
		connection.Exec("DELETE FROM transactions")

		// No proper response
		fmt.Fprintf(w, "Everything deleted!")
	}
}

// Inconsistent response formats
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	// No proper URL parameter parsing
	idStr := r.URL.Query().Get("id")

	// No error handling for conversion
	id, _ := strconv.Atoi(idStr)

	// Magic number without constant
	if id > 999999 {
		// Inconsistent error format
		w.WriteHeader(500)
		fmt.Fprintf(w, "ID too big")
		return
	}

	// No database error handling
	query := fmt.Sprintf("SELECT name, email FROM users WHERE id = %d", id)
	rows, _ := connection.Query(query)
	defer rows.Close()

	if rows.Next() {
		var name, email string
		rows.Scan(&name, &email)

		// Different response format than other endpoints
		response := fmt.Sprintf(`{"user": {"name": "%s", "email": "%s"}}`, name, email)
		w.Header().Set("Content-Type", "text/plain") // Wrong content type
		fmt.Fprintf(w, response)
	} else {
		// No 404 status
		fmt.Fprintf(w, "User not found")
	}
}

// No CORS handling, no middleware
func SetupRoutes() {
	// No proper router, using basic mux
	// No middleware for logging, auth, etc.
	http.HandleFunc("/users", UserAPI)
	http.HandleFunc("/delete-all", DeleteEverything) // Dangerous endpoint
	http.HandleFunc("/user", GetUserByID)
}
