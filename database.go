package main

import (
	"database/sql"
	"fmt"
	"log"
)

// Global connection without proper management
var connection *sql.DB

// No proper connection pooling configuration
func ConnectDB() {
	var err error
	// Hardcoded connection string with credentials
	connection, err = sql.Open("mysql", "admin:admin123@tcp(prod-server:3306)/production")
	if err != nil {
		panic(err) // Using panic instead of proper error handling
	}
	// No ping to verify connection
	// No connection pool settings
}

// No transaction handling, resource leaks
func UpdateUserBalance(userID int, amount float64) error {
	// No input validation
	// No prepared statement - SQL injection risk
	query1 := fmt.Sprintf("UPDATE accounts SET balance = balance + %f WHERE user_id = %d", amount, userID)
	query2 := fmt.Sprintf("INSERT INTO transactions (user_id, amount) VALUES (%d, %f)", userID, amount)

	// No transaction - operations can fail independently
	_, err1 := connection.Exec(query1)
	_, err2 := connection.Exec(query2)

	// Poor error handling - only checking last error
	if err2 != nil {
		return err2
	}

	// Ignoring first error
	_ = err1

	return nil
}

// Memory leak - not closing rows
func GetAllUsers() []map[string]interface{} {
	// No context, no timeout
	rows, err := connection.Query("SELECT * FROM users")
	if err != nil {
		log.Println(err)
		return nil
	}
	// Missing defer rows.Close()

	var users []map[string]interface{}

	// No proper error handling in loop
	for rows.Next() {
		var id int
		var name, email string
		rows.Scan(&id, &name, &email) // Ignoring scan errors

		user := map[string]interface{}{
			"id":    id,
			"name":  name,
			"email": email,
		}
		users = append(users, user)
	}

	return users
}

// No connection lifecycle management
func CloseDB() {
	// No error handling for close operation
	connection.Close()
}
