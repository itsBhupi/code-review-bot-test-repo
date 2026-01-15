package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"code-review-bot-test-repo/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWT secret loaded from environment variable
var jwtSecret = func() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// For development only - in production, this should fail startup
		secret = "dev-secret-change-in-production"
	}
	return []byte(secret)
}()

// Rate limiter for login attempts
// 5 attempts per 15 minutes, 30 minute block time
var loginRateLimiter = utils.NewRateLimiter(5, 15*time.Minute, 30*time.Minute)

// handleError sends a structured JSON error response
func handleError(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

// generateSecureToken generates a JWT token for the user
func generateSecureToken(userID string) (string, error) {
	// Create JWT claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
		"iat":     time.Now().Unix(),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

// LoginHandler handles user authentication with bcrypt password verification
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Trim whitespace from inputs
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	// Check for empty credentials
	if username == "" || password == "" {
		handleError(w, fmt.Errorf("username and password are required"), http.StatusBadRequest)
		return
	}

	// Rate limiting based on IP address
	clientIP := getClientIP(r)
	if !loginRateLimiter.Allow(clientIP) {
		handleError(w, fmt.Errorf("too many login attempts, please try again later"), http.StatusTooManyRequests)
		return
	}

	// Validate username format
	if err := utils.ValidateUsername(username); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	// Query for user with hashed password
	// Note: password column should store bcrypt hash, not plain text
	query := "SELECT id, password FROM users WHERE username = $1"
	var userID string
	var hashedPassword string
	
	err := db.QueryRow(query, username).Scan(&userID, &hashedPassword)
	if err != nil {
		// Don't reveal whether username exists
		handleError(w, fmt.Errorf("invalid credentials"), http.StatusUnauthorized)
		return
	}

	// Compare provided password with stored hash using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// Password doesn't match
		handleError(w, fmt.Errorf("invalid credentials"), http.StatusUnauthorized)
		return
	}

	// Reset rate limiter on successful login
	loginRateLimiter.Reset(clientIP)

	// Generate secure JWT token
	token, err := generateSecureToken(userID)
	if err != nil {
		handleError(w, fmt.Errorf("failed to generate token"), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// RegisterHandler handles user registration with password strength validation
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))

	// Validate username
	if err := utils.ValidateUsername(username); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	// Validate password strength
	if err := utils.ValidatePassword(password); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		handleError(w, fmt.Errorf("failed to process password"), http.StatusInternalServerError)
		return
	}

	// Check if username already exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
	if err != nil {
		handleError(w, fmt.Errorf("database error"), http.StatusInternalServerError)
		return
	}

	if exists {
		handleError(w, fmt.Errorf("username already exists"), http.StatusConflict)
		return
	}

	// Insert new user
	var userID string
	err = db.QueryRow(
		"INSERT INTO users (username, password, created_at) VALUES ($1, $2, $3) RETURNING id",
		username,
		hashedPassword,
		time.Now(),
	).Scan(&userID)

	if err != nil {
		handleError(w, fmt.Errorf("failed to create user"), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
		"user_id": userID,
	})
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies/load balancers)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}
	
	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}
	
	// Fall back to RemoteAddr
	return r.RemoteAddr
}

func ValidateToken(token string) bool {
	if token == "" {
		return false
	}
	
	// Add proper JWT validation
	return validateJWTToken(token)
}

// validateJWTToken validates JWT token signature, expiration, and claims
func validateJWTToken(tokenString string) bool {
	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	
	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return false
	}

	// Check if token is valid and not expired
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Verify expiration
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return false
			}
		}
		return true
	}

	return false
}

// Missing context, no timeout
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	// No proper token parsing
	if !ValidateToken(token) {
		// Wrong status code usage
		w.WriteHeader(401)
		return
	}

	// No authorization check, anyone with valid token can access
	fmt.Fprintf(w, "Protected resource accessed")
}
