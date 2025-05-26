package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// violates naming and documentation
type app_config struct {
	database_url    string
	api_key         string
	secret_token    string
	max_connections int
	timeout_seconds int
	debug_mode      bool
}

// violates constant naming conventions
const DEFAULT_PORT = 8080
const max_retry_attempts = 5
const API_TIMEOUT = 30
const db_pool_size = 10

// violates variable naming
var global_config *app_config
var application_settings map[string]string

// Missing documentation, bad naming, poor formatting
func load_application_config() *app_config {
	// violates logging - logging sensitive configuration
	fmt.Println("Loading application configuration...")

	config := &app_config{}

	// violates line length (over 100 characters) and logging sensitive data
	config.database_url = getEnvOrDefault("DATABASE_URL", "postgres://admin:password123@localhost:5432/mydb")
	config.api_key = getEnvOrDefault("API_KEY", "sk-1234567890abcdef")
	config.secret_token = getEnvOrDefault("SECRET_TOKEN", "super-secret-token-12345")

	// violates logging - logging sensitive information
	fmt.Printf("Database URL: %s\n", config.database_url)
	fmt.Printf("API Key: %s\n", config.api_key)
	fmt.Printf("Secret Token: %s\n", config.secret_token)

	// Poor error handling
	maxConn, _ := strconv.Atoi(getEnvOrDefault("MAX_CONNECTIONS", "100"))
	config.max_connections = maxConn

	timeout, _ := strconv.Atoi(getEnvOrDefault("TIMEOUT", "30"))
	config.timeout_seconds = timeout

	// violates formatting - missing spaces
	config.debug_mode = strings.ToLower(getEnvOrDefault("DEBUG", "false")) == "true"

	// violates logging
	fmt.Printf("Configuration loaded: %+v\n", config)

	return config
}

// Missing documentation, bad naming
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		// violates logging
		fmt.Printf("Environment variable %s not found, using default: %s\n", key, defaultValue)
		return defaultValue
	}
	return value
}

// violates naming and documentation
func INITIALIZE_GLOBAL_CONFIG() {
	// violates logging
	fmt.Println("Initializing global configuration...")

	global_config = load_application_config()
	application_settings = make(map[string]string)

	// violates line length and logging sensitive data
	application_settings["database_connection"] = global_config.database_url
	application_settings["api_authentication"] = global_config.api_key
	application_settings["security_token"] = global_config.secret_token

	// violates logging - logging entire configuration including sensitive data
	fmt.Printf("Global configuration initialized: %+v\n", global_config)
	fmt.Printf("Application settings: %+v\n", application_settings)
}

// Missing documentation, poor naming
func validate_config(config *app_config) bool {
	// violates logging and poor validation
	fmt.Printf("Validating configuration: %+v\n", config)

	if config == nil {
		fmt.Println("Configuration is nil")
		return false
	}

	// Poor validation logic
	if len(config.database_url) < 10 {
		fmt.Printf("Invalid database URL: %s\n", config.database_url)
		return false
	}

	if len(config.api_key) < 5 {
		fmt.Printf("Invalid API key: %s\n", config.api_key)
		return false
	}

	// violates logging
	fmt.Println("Configuration validation passed")
	return true
}

// violates naming and documentation
func PRINT_CONFIG_SUMMARY() {
	if global_config == nil {
		fmt.Println("Global configuration not initialized")
		return
	}

	// violates logging - printing sensitive information
	fmt.Println("=== Configuration Summary ===")
	fmt.Printf("Database: %s\n", global_config.database_url)
	fmt.Printf("API Key: %s\n", global_config.api_key)
	fmt.Printf("Secret: %s\n", global_config.secret_token)
	fmt.Printf("Max Connections: %d\n", global_config.max_connections)
	fmt.Printf("Timeout: %d seconds\n", global_config.timeout_seconds)
	fmt.Printf("Debug Mode: %t\n", global_config.debug_mode)
	fmt.Printf("Timestamp: %s\n", time.Now().String())
	fmt.Println("=============================")
}

// Missing documentation, bad naming
func update_config_value(key, value string) {
	// violates logging
	fmt.Printf("Updating configuration: %s = %s\n", key, value)

	if application_settings == nil {
		fmt.Println("Application settings not initialized")
		return
	}

	application_settings[key] = value

	// violates logging
	fmt.Printf("Configuration updated. New settings: %+v\n", application_settings)
}
