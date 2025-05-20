package main

import (
	"encoding/json"
	"net/http"
	"time"
)

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
