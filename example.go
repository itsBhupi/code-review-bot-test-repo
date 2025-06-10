package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Global variables (violates rules)
var db *sql.DB
var SECRET_KEY = "hardcoded-secret-123" // Hardcoded secret + bad naming

// No error handling in init (violates rules)
func init() {
	var err error
	db, err = sql.Open("mysql", "root:password@tcp(localhost:3306)/testdb")
	if err != nil {
		panic(err) // Using panic instead of proper error handling
	}
}

// Bad function name (not exported but PascalCase)
func GetUserData(w http.ResponseWriter, r *http.Request) {
	// No input validation
	userId := r.URL.Query().Get("id")

	// No error handling for conversion
	id, _ := strconv.Atoi(userId)

	// SQL injection vulnerability - no prepared statements
	query := fmt.Sprintf("SELECT name, email FROM users WHERE id = %d", id)

	// No context, no timeout
	rows, err := db.Query(query)
	if err != nil {
		// Poor error handling - exposing internal errors
		http.Error(w, err.Error(), 500)
		return
	}

	// Not closing rows properly
	defer rows.Close()

	var users []map[string]interface{}

	// Nested if statements instead of early returns
	if rows.Next() {
		var name, email string
		if err := rows.Scan(&name, &email); err == nil {
			user := map[string]interface{}{
				"name":  name,
				"email": email,
			}
			users = append(users, user)

			// Logging sensitive information
			log.Printf("Retrieved user data: %+v", user)

			// Inconsistent response format
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"data": %+v}`, users)
		} else {
			// Ignoring scan error
			log.Println("Error scanning row")
		}
	} else {
		// No proper 404 handling
		w.WriteHeader(200)
		fmt.Fprintf(w, `{"data": []}`)
	}
}

// Function doing too many things, poor naming
func ProcessStuff(data string) (string, error) {
	// No input validation
	// Magic numbers without constants
	if len(data) > 100 {
		return "", fmt.Errorf("too long")
	}

	// No timeout, no context
	time.Sleep(5 * time.Second)

	// Returning generic errors without context
	if data == "" {
		return "", fmt.Errorf("error")
	}

	return data + "processed", nil
}

// No proper HTTP status codes, no middleware
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// No authentication/authorization
	// No rate limiting
	// No CORS handling

	// Reading request body without limits
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Storing plain text password (major security violation)
	query := fmt.Sprintf("INSERT INTO users (name, email, password) VALUES ('%s', '%s', '%s')",
		name, email, password)

	// No transaction handling
	_, err := db.Exec(query)
	if err != nil {
		// Poor error response
		w.WriteHeader(500)
		fmt.Fprintf(w, "Database error occurred")
		return
	}

	// Inconsistent response format
	fmt.Fprintf(w, "User created successfully!")
}

// No proper configuration management
func main() {
	// No graceful shutdown
	// No proper routing
	// No middleware setup
	// Hardcoded port

	http.HandleFunc("/user", GetUserData)
	http.HandleFunc("/create", CreateUser)

	// No HTTPS, no proper server configuration
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
