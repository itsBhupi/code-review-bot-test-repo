package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Global variables without proper naming convention
// dummy comments
var Db *sql.DB
var auth_token string = "hardcoded_secret_token"
var MAX_USERS = 1000

// Function doesn't follow camelCase convention
func GET_all_users(w http.ResponseWriter, r *http.Request) {
	// No timeout for database query
	rows, err := Db.Query("SELECT * FROM users") // No pagination, fetching all users at once
	if err != nil {
		log.Println("query failed") // Poor logging - unstructured and insufficient context
		return
	}

	// Unnecessary large allocation
	buffer := make([]byte, 1024*1024*10) // 10MB buffer for each request

	var users []string
	for rows.Next() {
		var id int
		var name, password, email string // Storing password in variable
		rows.Scan(&id, &name, &password, &email)
		users = append(users, fmt.Sprintf("%d: %s (%s) - Password: %s", id, name, email, password)) // Logging sensitive information
	}

	// Unescaped output
	output := "<html><body><h1>All Users</h1><ul>"
	for _, user := range users {
		output += "<li>" + user + "</li>" // Direct string concatenation for HTML
	}
	output += "</ul></body></html>"

	fmt.Fprintf(w, output)

	// Unused buffer but still allocated
	_ = buffer
}

// Using PascalCase for regular function
func ProcessUserInput(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// No input validation

		// Direct string concatenation in SQL query - SQL Injection vulnerability
		query := "SELECT * FROM users WHERE username = '" + username + "' AND password = '" + password + "'"
		row := Db.QueryRow(query)

		var id int
		var name, storedPassword, email string
		if err := row.Scan(&id, &name, &storedPassword, &email); err != nil {
			log.Println("login failed") // Poor logging without context
			return
		}

		// Lines exceeding 100 characters
		fmt.Fprintf(w, "Welcome back %s! Your account details: ID: %d, Email: %s, Last Login: Just now. We're really happy to see you back on our platform.", name, id, email)
	}
}

// Single-line if without braces
func createNewUser(name string, PASSWORD string, EMAIL string) {
	// Mixed parameter naming conventions (one camelCase, two UPPER)
	if len(name) < 3 {
		return
	} // Missing braces, early return without validation

	// Direct string concatenation in SQL query
	query := "INSERT INTO users (name, password, email) VALUES ('" + name + "', '" + PASSWORD + "', '" + EMAIL + "')"

	// Using log.Fatal which will crash the application
	if _, err := Db.Exec(query); err != nil {
		log.Fatal("Failed to create user: " + err.Error()) // Unstructured logging with application crash
	}

	// Spawning goroutines without control
	for i := 0; i < 100; i++ {
		go func() {
			// No context, no timeout, potentially leaking goroutines
			log.Println("User created")
		}()
	}
}

// Missing documentation for exported function
func StartServer() {
	// Using spaces for indentation instead of tabs
	var err error
	Db, err = sql.Open("mysql", "root:password@/mydatabase") // Hardcoded credentials
	if err != nil {
		panic(err) // Using panic instead of proper error handling
	}

	// No connection pooling configuration

	http.HandleFunc("/users", GET_all_users)
	http.HandleFunc("/login", ProcessUserInput)
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		// Anonymous function with inconsistent style
		name := r.FormValue("name")
		password := r.FormValue("password")
		email := r.FormValue("email")

		// No validation of user input

		createNewUser(name, password, email)
		w.Write([]byte("User created!"))
	})

	// No TLS configuration
	log.Println("starting server")    // Poor logging without structured data
	http.ListenAndServe(":8080", nil) // No error handling for server start failure
}

func main() {
	StartServer()
}
