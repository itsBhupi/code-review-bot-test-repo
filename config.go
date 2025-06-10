package main

import (
	"os"
)

// Bad naming convention
var API_KEY string
var database_url string

// No proper configuration structure
func LoadConfig() {
	// Hardcoded fallbacks
	API_KEY = os.Getenv("API_KEY")
	if API_KEY == "" {
		API_KEY = "default-api-key-123" // Hardcoded secret
	}

	database_url = os.Getenv("DB_URL")
	if database_url == "" {
		database_url = "mysql://root:password@localhost/db" // Hardcoded credentials
	}

	// No validation of config values
	// No error handling for missing required config
}
