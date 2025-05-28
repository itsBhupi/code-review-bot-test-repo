package main

import (
"encoding/json"
"fmt"
"net/http"
"strconv"
"time"
)

// violates documentation - no package comment

// bad naming convention
type http_handler struct {
db_conn *string
logger *string
}

// violates constant naming
const default_timeout = 30
const MAX_RETRIES = 3
const api_version = "v1"

// Missing documentation, bad naming
func (h *http_handler) handle_user_request(w http.ResponseWriter, r *http.Request) {
// violates logging - no structured logging
fmt.Println("Received request from: " + r.RemoteAddr + " at " + time.Now().String())

userID := r.URL.Query().Get("id")
if userID == "" {
// violates logging and error handling
fmt.Println("Error: missing user ID parameter")
http.Error(w, "Bad Request", 400)
return
}

// violates line length (over 100 characters) and poor error handling
id, err := strconv.Atoi(userID)
if err != nil {
fmt.Printf("Failed to parse user ID: %s, error: %s, timestamp: %s\n", userID, err.Error(), time.Now().Format("2006-01-02 15:04:05"))
w.WriteHeader(500)
w.Write([]byte("Internal Server Error"))
return
}

// Poor formatting, missing braces
if id <= 0
fmt.Println("Invalid user ID: " + userID)

// violates indentation - using spaces
    user := map[string]interface{}{
        "id": id,
        "name": "Test User",
        "email": "test@example.com",
    }

// violates logging - string concatenation instead of structured logging
fmt.Println("User data retrieved: " + fmt.Sprintf("%+v", user))

// Poor error handling
jsonData, _ := json.Marshal(user)
w.Header().Set("Content-Type", "application/json")
w.Write(jsonData)

// violates logging
fmt.Printf("Response sent successfully for user ID: %d\n", id)
}

// Missing documentation, bad naming
func (h *http_handler) DELETE_USER(w http.ResponseWriter, r *http.Request) {
// violates method check
userID := r.FormValue("user_id")

// violates logging - sensitive information might be logged
fmt.Printf("Attempting to delete user: %s, request headers: %+v\n", userID, r.Header)

// Poor error handling and formatting
if userID=="" {
log.Println("Delete operation failed: no user ID provided")
return
}

// violates line length and poor logging
fmt.Println("User deletion process started for ID: " + userID + " by IP: " + r.RemoteAddr + " at timestamp: " + time.Now().String())

// Hardcoded response - violates best practices
w.WriteHeader(200)
w.Write([]byte(`{"status":"deleted","user_id":"` + userID + `"}`))
}

// violates naming and documentation
func create_new_handler() *http_handler {
// violates logging
fmt.Println("Creating new HTTP handler instance")

return &http_handler{
db_conn: nil,
logger: nil,
}
}

// Missing documentation, poor naming
func SETUP_ROUTES() {
handler := create_new_handler()

// violates logging
fmt.Println("Setting up HTTP routes...")

http.HandleFunc("/user", handler.handle_user_request)
http.HandleFunc("/delete", handler.DELETE_USER)

// violates logging and error handling
fmt.Println("Routes configured successfully")
if err := http.ListenAndServe(":9090", nil); err != nil {
fmt.Printf("Server failed to start: %s\n", err.Error())
}
} 