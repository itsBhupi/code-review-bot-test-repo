package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// TestDuplicateDetection - This file is designed to test the duplicate detection service
// It contains intentional code issues that should trigger AI review comments

// UserService handles user operations
type UserService struct {
	users map[string]User
}

// User represents a user in the system
type User struct {
	ID       string
	Name     string
	Email    string
	Password string // This should trigger a security comment about storing passwords in plain text
	Age      int
}

// CreateUser creates a new user - has multiple issues
func (s *UserService) CreateUser(name, email, password string, ageStr string) (*User, error) {
	// Issue 1: No input validation
	// Issue 2: Converting string to int without proper error handling
	age, _ := strconv.Atoi(ageStr) // Should trigger error handling comment

	// Issue 3: No duplicate email check
	// Issue 4: Storing password in plain text
	user := &User{
		ID:       generateID(), // This function doesn't exist - should trigger error
		Name:     name,
		Email:    email,
		Password: password, // Plain text password storage
		Age:      age,
	}

	s.users[user.ID] = *user
	return user, nil
}

// GetUserByEmail retrieves user by email - inefficient implementation
func (s *UserService) GetUserByEmail(email string) *User {
	// Issue: Inefficient O(n) search instead of using a proper index
	for _, user := range s.users {
		if user.Email == email {
			return &user
		}
	}
	return nil
}

// UpdateUserPassword updates user password - security issues
func (s *UserService) UpdateUserPassword(userID, newPassword string) error {
	// Issue 1: No authentication check
	// Issue 2: No password validation
	// Issue 3: Still storing in plain text
	user, exists := s.users[userID]
	if !exists {
		return fmt.Errorf("user not found")
	}

	user.Password = newPassword // Direct assignment without hashing
	s.users[userID] = user
	return nil
}

// HTTPHandler handles HTTP requests - multiple issues
func HTTPHandler(w http.ResponseWriter, r *http.Request) {
	// Issue 1: No error handling for request parsing
	email := r.URL.Query().Get("email")

	// Issue 2: Direct string concatenation (potential injection)
	query := "SELECT * FROM users WHERE email = '" + email + "'"

	// Issue 3: Hardcoded response without proper content type
	fmt.Fprintf(w, "Query: %s", query)
}

// ProcessFile processes a file - resource leak issues
func ProcessFile(filename string) error {
	// Issue 1: No input validation
	// Issue 2: Not handling file close properly (resource leak)
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening file: %v", err) // Basic logging
		return err
	}

	// Issue 3: Missing defer file.Close()
	// This will cause resource leaks

	content := make([]byte, 1024)
	_, err = file.Read(content)
	if err != nil {
		return err // File handle leaked here
	}

	// Issue 4: Not checking file close error
	file.Close()

	return nil
}

// DatabaseOperation simulates database operations - error handling issues
func DatabaseOperation(userID string) (string, error) {
	// Issue 1: No input sanitization
	// Issue 2: Potential SQL injection if this were real SQL
	query := fmt.Sprintf("SELECT name FROM users WHERE id = %s", userID)

	// Issue 3: Simulating database error without proper handling
	if userID == "" {
		return "", fmt.Errorf("invalid user ID")
	}

	// Issue 4: No context timeout for database operation
	// Issue 5: No connection pooling considerations
	result := executeQuery(query) // This function doesn't exist

	return result, nil
}

// ConcurrentOperation demonstrates concurrency issues
func ConcurrentOperation(data map[string]int, key string, value int) {
	// Issue 1: Race condition - map is not safe for concurrent access
	// Issue 2: No synchronization mechanism
	data[key] = value // Potential race condition

	// Issue 3: No error handling for map operations
	total := 0
	for _, v := range data {
		total += v // Another potential race condition
	}

	fmt.Printf("Total: %d\n", total)
}

// JSONProcessing handles JSON data - type safety issues
func JSONProcessing(jsonStr string) interface{} {
	// Issue 1: Using interface{} instead of proper types
	// Issue 2: No JSON validation
	var result interface{}

	// Issue 3: No error handling for JSON unmarshaling
	json.Unmarshal([]byte(jsonStr), &result) // Missing error check

	// Issue 4: Type assertion without safety check
	data := result.(map[string]interface{}) // Potential panic

	// Issue 5: Not handling nil cases
	return data["field"] // Could panic if field doesn't exist
}

// StringManipulation shows string handling issues
func StringManipulation(input string) string {
	// Issue 1: No input validation for nil/empty
	// Issue 2: Inefficient string concatenation in loop
	result := ""
	words := strings.Split(input, " ")

	for i, word := range words {
		// Issue 3: String concatenation in loop (should use strings.Builder)
		result += word
		if i < len(words)-1 {
			result += "-" // Inefficient concatenation
		}
	}

	// Issue 4: No bounds checking
	return result[0:10] // Potential panic if result is shorter than 10 chars
}

// NetworkOperation demonstrates network handling issues
func NetworkOperation(url string) ([]byte, error) {
	// Issue 1: No URL validation
	// Issue 2: No timeout configuration
	resp, err := http.Get(url) // No timeout, could hang indefinitely
	if err != nil {
		return nil, err
	}

	// Issue 3: Not deferring body close
	// Issue 4: Not checking resp.StatusCode
	body := make([]byte, 1000)
	resp.Body.Read(body) // No error handling

	resp.Body.Close() // Should be deferred

	return body, nil
}

// MemoryLeak demonstrates potential memory leak
func MemoryLeak() {
	// Issue 1: Growing slice without bounds
	var data [][]byte

	for i := 0; i < 1000000; i++ {
		// Issue 2: Allocating large chunks without cleanup
		chunk := make([]byte, 1024*1024) // 1MB per iteration
		data = append(data, chunk)

		// Issue 3: No cleanup or size limits
		// This will consume significant memory
	}

	// Issue 4: Data never gets cleaned up or used
	fmt.Printf("Allocated %d chunks\n", len(data))
}

// main function with basic issues
func main() {
	// Issue 1: No proper configuration management
	// Issue 2: Hardcoded values
	service := &UserService{
		users: make(map[string]User),
	}

	// Issue 3: No error handling for user creation
	user, _ := service.CreateUser("John Doe", "john@example.com", "password123", "25")

	// Issue 4: Printing sensitive information
	fmt.Printf("Created user: %+v\n", user) // This will print the password

	// Issue 5: No graceful shutdown handling
	// Issue 6: No proper logging setup
	log.Println("Application started")
}
