package utils

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// Global variables everywhere
var GLOBAL_COUNTER int
var cache map[string]interface{} // No proper synchronization

// Response represents a standard API response structure
type Response struct {
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// SendJSONResponse sends a JSON response with the given status code and data
func SendJSONResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	response := Response{
		Status:    http.StatusText(statusCode),
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// FormatDate formats a time.Time to a string in the format "2006-01-02 15:04:05"
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// IsValidEmail performs a basic email validation
func IsValidEmail(email string) bool {
	// This is a very basic validation. In production, use a proper email validation library
	return len(email) > 3 && len(email) < 254
}

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}

// Poor function naming, doing too many things
func DO_EVERYTHING(input string) string {
	// No input validation
	// Multiple responsibilities in one function

	// Using deprecated MD5 for hashing
	hash := md5.Sum([]byte(input))
	hashStr := fmt.Sprintf("%x", hash)

	// Modifying global state without protection
	GLOBAL_COUNTER++

	// No proper error handling
	cache[input] = hashStr

	// Magic number without constant
	if len(input) > 50 {
		return "too long"
	}

	return strings.ToUpper(hashStr)
}

// No proper error handling, ignoring all errors
func ReadConfigFile() map[string]string {
	// Hardcoded file path
	content, _ := os.ReadFile("/etc/myapp/config.txt") // Ignoring error

	config := make(map[string]string)

	// Poor parsing logic
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		// No validation, assumes format
		parts := strings.Split(line, "=")
		config[parts[0]] = parts[1] // Potential index out of range
	}

	return config
}

// Goroutine leak - no proper cleanup
func StartBackgroundTask() {
	go func() {
		for {
			// Infinite loop with no cancellation
			// No context for cancellation
			time.Sleep(1 * time.Second)

			// Modifying global state from goroutine
			GLOBAL_COUNTER++

			// No error handling
			fmt.Println("Background task running...")
		}
	}()
	// No way to stop this goroutine
}

// Poor naming, no clear purpose
func helper_func(a, b int) int {
	// No validation
	// Potential division by zero
	return a / b
}

// Using panic for normal error conditions
func ValidateEmail(email string) {
	// No proper validation
	if !strings.Contains(email, "@") {
		panic("Invalid email") // Should return error instead
	}
}

// No proper resource cleanup
func ProcessFile(filename string) string {
	// No error handling
	file, _ := os.Open(filename)
	// Missing defer file.Close()

	content, _ := os.ReadFile(filename) // Reading file twice

	return string(content)
}

// init function with side effects
func init() {
	// Initializing global state in init
	cache = make(map[string]interface{})

	// Side effects in init
	fmt.Println("Utils package initialized")

	// Starting goroutines in init
	StartBackgroundTask()
}
