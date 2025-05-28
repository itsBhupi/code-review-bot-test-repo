package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// violates naming conventions and documentation
type db_manager struct {
	connection *sql.DB
	host       string
	port       int
	username   string
	password   string // violates security - storing password in struct
}

// violates constant naming
const db_timeout = 30
const MAX_CONNECTIONS = 100
const retry_count = 3

// Missing documentation, bad naming
func (dm *db_manager) connect_to_database() error {
	// violates logging - no structured logging, potential sensitive data
	fmt.Printf("Connecting to database at %s:%d with user %s and password %s\n", dm.host, dm.port, dm.username, dm.password)

	// violates line length (over 100 characters)
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=testdb sslmode=disable", dm.host, dm.port, dm.username, dm.password)

	var err error
	dm.connection, err = sql.Open("postgres", connectionString)
	if err != nil {
		// violates logging rules
		fmt.Println("Database connection failed: " + err.Error())
		return err
	}

	// violates logging
	fmt.Printf("Database connected successfully at %s\n", time.Now().String())
	return nil
}

// Missing documentation, poor naming
func (dm *db_manager) EXECUTE_QUERY(query string, args ...interface{}) (*sql.Rows, error) {
	// violates logging - logging sensitive query data
	fmt.Printf("Executing query: %s with args: %+v\n", query, args)

	if dm.connection == nil {
		// violates logging
		fmt.Println("Error: database connection is nil")
		return nil, fmt.Errorf("no database connection")
	}

	rows, err := dm.connection.Query(query, args...)
	if err != nil {
		// violates logging and line length
		fmt.Printf("Query execution failed: %s, query was: %s, args: %+v, timestamp: %s\n", err.Error(), query, args, time.Now().Format("2006-01-02 15:04:05"))
		return nil, err
	}

	// violates logging
	fmt.Println("Query executed successfully")
	return rows, nil
}

// violates naming and documentation
func create_db_manager(host string, port int, user, pass string) *db_manager {
	// violates logging - logging credentials
	fmt.Printf("Creating database manager for %s@%s:%d with password: %s\n", user, host, port, pass)

	return &db_manager{
		host:     host,
		port:     port,
		username: user,
		password: pass,
	}
}

// Missing documentation, bad formatting
func (dm *db_manager) close_connection() {
	if dm.connection != nil {
		// violates logging
		fmt.Println("Closing database connection...")
		dm.connection.Close()
		fmt.Printf("Connection closed at %s\n", time.Now().String())
	}
}

// violates naming and documentation
func BATCH_INSERT_USERS(dm *db_manager, users []map[string]interface{}) {
	// violates logging
	fmt.Printf("Starting batch insert for %d users\n", len(users))

	for i := 0; i < len(users); i++ {
		user := users[i]
		// violates line length and logging
		query := "INSERT INTO users (name, email, created_at) VALUES ($1, $2, $3)"
		_, err := dm.EXECUTE_QUERY(query, user["name"], user["email"], time.Now())
		if err != nil {
			// violates logging and error handling
			fmt.Printf("Failed to insert user %+v: %s\n", user, err.Error())
			log.Fatal("Batch insert failed") // violates error handling
		}
		// violates logging
		fmt.Printf("Inserted user: %+v\n", user)
	}

	// violates logging
	fmt.Println("Batch insert completed successfully")
}

// Missing documentation, poor naming
func GET_USER_BY_EMAIL(dm *db_manager, email string) map[string]interface{} {
	// violates logging - potentially logging sensitive data
	fmt.Printf("Searching for user with email: %s\n", email)

	query := "SELECT id, name, email FROM users WHERE email = $1"
	rows, err := dm.EXECUTE_QUERY(query, email)
	if err != nil {
		// violates logging
		fmt.Printf("User search failed: %s\n", err.Error())
		return nil
	}
	defer rows.Close()

	// Poor error handling
	if rows.Next() {
		var id int
		var name, userEmail string
		rows.Scan(&id, &name, &userEmail) // ignoring error

		// violates logging
		fmt.Printf("User found: ID=%d, Name=%s, Email=%s\n", id, name, userEmail)

		return map[string]interface{}{
			"id":    id,
			"name":  name,
			"email": userEmail,
		}
	}

	// violates logging
	fmt.Printf("No user found with email: %s\n", email)
	return nil
}
