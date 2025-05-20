package main

import (
	"net/http"
	"time"

	"gorm.io/gorm"
)

// VIOLATION: No timeout configuration
func getAllUsers(db *gorm.DB) []map[string]interface{} {
	// VIOLATION: No pagination, fetching all users at once
	var results []map[string]interface{}
	db.Table("users").Find(&results)
	return results
}

// VIOLATION: Creates a large buffer for each request
func processLargeData(w http.ResponseWriter, r *http.Request) {
	// VIOLATION: Unnecessary large allocation
	buffer := make([]byte, 1024*1024*10) // 10MB buffer

	// VIOLATION: No defer to close resources
	f, _ := r.MultipartForm.File["upload"][0].Open()

	// Pretend to use f to avoid unused variable warning
	_ = f

	// Use buffer and file...

	// VIOLATION: Resource leak - never closes the file

	// VIOLATION: Unused buffer but still allocated
	_ = buffer
}

// VIOLATION: Spawning goroutines without control or limits
func notifyAllUsers(message string) {
	// VIOLATION: No limiting mechanism, creating 1000s of goroutines
	for i := 0; i < 1000; i++ {
		go func(id int) {
			// VIOLATION: No context, no timeout
			// Simulate sending notification
			time.Sleep(time.Second)
		}(i)
	}
}
