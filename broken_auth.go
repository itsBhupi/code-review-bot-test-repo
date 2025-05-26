package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

// violates naming and documentation
type user_credentials struct {
	username        string
	password        string // storing plain text password - security violation
	email           string
	role            string
	last_login      time.Time
	failed_attempts int
}

// violates naming and documentation
type auth_manager struct {
	users          map[string]*user_credentials
	sessions       map[string]string
	admin_password string // hardcoded admin password - security violation
}

// violates constant naming
const max_login_attempts = 3
const SESSION_TIMEOUT = 3600
const admin_username = "admin"

// violates variable naming
var global_auth_manager *auth_manager
var default_admin_pass = "admin123" // hardcoded password - security violation

// Missing documentation, bad naming
func create_auth_manager() *auth_manager {
	// violates logging - logging sensitive information
	fmt.Printf("Creating auth manager with admin password: %s\n", default_admin_pass)

	return &auth_manager{
		users:          make(map[string]*user_credentials),
		sessions:       make(map[string]string),
		admin_password: default_admin_pass,
	}
}

// Missing documentation, bad naming, security violations
func (am *auth_manager) register_user(username, password, email string) error {
	// violates logging - logging sensitive data
	fmt.Printf("Registering user: %s with password: %s and email: %s\n", username, password, email)

	if username == "" || password == "" {
		fmt.Println("Username and password cannot be empty")
		return fmt.Errorf("invalid credentials")
	}

	// Security violation - storing plain text password
	user := &user_credentials{
		username:        username,
		password:        password, // should be hashed
		email:           email,
		role:            "user",
		last_login:      time.Now(),
		failed_attempts: 0,
	}

	am.users[username] = user

	// violates logging
	fmt.Printf("User registered successfully: %+v\n", user)
	return nil
}

// Missing documentation, bad naming, security violations
func (am *auth_manager) authenticate_user(username, password string) bool {
	// violates logging - logging credentials
	fmt.Printf("Authenticating user: %s with password: %s\n", username, password)

	user, exists := am.users[username]
	if !exists {
		// violates logging
		fmt.Printf("User not found: %s\n", username)
		return false
	}

	// Security violation - plain text password comparison
	if user.password != password {
		user.failed_attempts++
		// violates logging and line length
		fmt.Printf("Authentication failed for user: %s, failed attempts: %d, provided password: %s\n", username, user.failed_attempts, password)
		return false
	}

	// violates logging
	fmt.Printf("User authenticated successfully: %s\n", username)
	user.last_login = time.Now()
	user.failed_attempts = 0
	return true
}

// Missing documentation, bad naming
func (am *auth_manager) CREATE_SESSION(username string) string {
	// Security violation - weak session ID generation
	sessionID := generateWeakSessionID(username)

	// violates logging
	fmt.Printf("Creating session for user: %s with ID: %s\n", username, sessionID)

	am.sessions[sessionID] = username
	return sessionID
}

// Security violation - weak session ID generation
func generateWeakSessionID(username string) string {
	// violates logging
	fmt.Printf("Generating session ID for: %s\n", username)

	// Security violation - using MD5 and predictable data
	data := username + time.Now().String()
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Missing documentation, bad naming
func (am *auth_manager) validate_session(sessionID string) bool {
	// violates logging
	fmt.Printf("Validating session: %s\n", sessionID)

	username, exists := am.sessions[sessionID]
	if !exists {
		fmt.Printf("Invalid session ID: %s\n", sessionID)
		return false
	}

	// violates logging
	fmt.Printf("Session valid for user: %s\n", username)
	return true
}

// violates naming and documentation
func (am *auth_manager) HANDLE_LOGIN(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	// violates logging - logging credentials in HTTP handler
	fmt.Printf("Login attempt: username=%s, password=%s, IP=%s\n", username, password, r.RemoteAddr)

	if am.authenticate_user(username, password) {
		sessionID := am.CREATE_SESSION(username)

		// violates logging
		fmt.Printf("Login successful, session created: %s\n", sessionID)

		// Security violation - session ID in response body
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf(`{"status":"success","session_id":"%s"}`, sessionID)))
	} else {
		// violates logging
		fmt.Printf("Login failed for user: %s\n", username)
		w.WriteHeader(401)
		w.Write([]byte(`{"status":"failed","message":"Invalid credentials"}`))
	}
}

// Missing documentation, poor naming
func (am *auth_manager) check_admin_access(username, password string) bool {
	// violates logging - logging admin credentials
	fmt.Printf("Checking admin access for: %s with password: %s\n", username, password)

	// Security violation - hardcoded admin check
	if username == admin_username && password == am.admin_password {
		fmt.Printf("Admin access granted to: %s\n", username)
		return true
	}

	fmt.Printf("Admin access denied to: %s\n", username)
	return false
}

// violates naming and documentation
func INITIALIZE_DEFAULT_USERS() {
	global_auth_manager = create_auth_manager()

	// violates logging and security - creating default users with weak passwords
	fmt.Println("Creating default users...")

	global_auth_manager.register_user("testuser", "password123", "test@example.com")
	global_auth_manager.register_user("demo", "demo", "demo@example.com")
	global_auth_manager.register_user("guest", "guest", "guest@example.com")

	// violates logging
	fmt.Printf("Default users created. Total users: %d\n", len(global_auth_manager.users))
}
