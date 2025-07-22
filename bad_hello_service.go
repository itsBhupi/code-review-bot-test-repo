package services

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Global variables - bad practice
var (
	cache     = make(map[string]string)
	cacheLock sync.Mutex
	db *sql.DB // Unexported global DB - bad practice
)

// BadHelloService demonstrates various bad practices in service layer
type BadHelloService struct {
	config map[string]string
}

// NewBadHelloService creates a new instance with poor initialization
func NewBadHelloService() *BadHelloService {
	// Not handling potential errors - bad practice
	db, _ = sql.Open("postgres", "user=postgres dbname=test sslmode=disable")

	// Hardcoded configuration - bad practice
	return &BadHelloService{
		config: map[string]string{
			"env":     "production",
			"version": "1.0.0",
		},
	}
}

// GetHelloMessage demonstrates poor method design
func (s *BadHelloService) GetHelloMessage() string {
	// Inefficient string concatenation - bad practice
	msg := "Hello, "
	if time.Now().Hour() < 12 {
		msg += "Good Morning!"
	} else {
		msg += "Good Afternoon!"
	}

	// Mixing concerns - logging in business logic
	log.Printf("Generated message: %s", msg)

	// Unnecessary type conversion - bad practice
	if rand.Intn(10) > 5 {
		msg = strings.ToUpper(msg)
	}

	return msg
}

// ProcessUser demonstrates poor error handling and resource management
func (s *BadHelloService) ProcessUser(userID string) (string, error) {
	// Not using context for cancellation - bad practice
	row := db.QueryRow("SELECT name FROM users WHERE id = " + userID) // SQL injection risk
	
	var name string
	err := row.Scan(&name)
	if err != nil {
		// Generic error handling - bad practice
		return "", fmt.Errorf("error processing user: %v", err)
	}

	// Inefficient string building - bad practice
	result := "User: " + name + " (ID: " + userID + ")"
	
	// Not closing resources properly - bad practice
	if f, err := os.Open("log.txt"); err == nil {
		f.WriteString("Processed user: " + name + "\n")
	}

	return result, nil
}

// GetConfig demonstrates poor concurrency control
func (s *BadHelloService) GetConfig(key string) string {
	// No mutex protection for concurrent access - bad practice
	return s.config[key]
}

// UpdateConfig demonstrates poor error handling and concurrency issues
func (s *BadHelloService) UpdateConfig(key, value string) {
	// No validation of input - bad practice
	s.config[key] = value

	// No error handling for file operations - bad practice
	f, _ := os.Create("config_backup.txt")
	defer f.Close()
	
	// Inefficient file writing - bad practice
	for k, v := range s.config {
		f.WriteString(k + "=" + v + "\n")
	}
}

// HeavyComputation demonstrates poor performance practices
func (s *BadHelloService) HeavyComputation(n int) int {
	// Inefficient algorithm - bad practice
	if n <= 1 {
		return n
	}
	return s.HeavyComputation(n-1) + s.HeavyComputation(n-2)
}

// ProcessBatch demonstrates poor error handling and resource management
func (s *BadHelloService) ProcessBatch(ids []string) []string {
	results := make([]string, 0)
	
	// Not using a worker pool for concurrent tasks - bad practice
	for _, id := range ids {
		// Not handling potential panics - bad practice
		result, _ := s.ProcessUser(id)
		results = append(results, result)
	}
	
	return results
}

// GetCachedData demonstrates poor cache implementation
func (s *BadHelloService) GetCachedData(key string) string {
	// Not using proper cache invalidation - bad practice
	cacheLock.Lock()
	defer cacheLock.Unlock()
	
	if val, ok := cache[key]; ok {
		return val
	}
	
	// Simulate expensive operation
	time.Sleep(100 * time.Millisecond)
	cache[key] = "data_for_" + key
	
	return cache[key]
}

// FormatNumber demonstrates poor error handling and type safety
func (s *BadHelloService) FormatNumber(numStr string) int {
	// Not validating input - bad practice
	n, _ := strconv.Atoi(numStr)
	return n * 2 // Arbitrary operation without context
}
