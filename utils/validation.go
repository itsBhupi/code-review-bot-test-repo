package utils

import (
	"fmt"
	"strings"
	"unicode"
)

// ValidateUsername checks if the username meets requirements
func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	
	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}
	
	if len(username) > 50 {
		return fmt.Errorf("username must not exceed 50 characters")
	}
	
	// Check for valid characters (alphanumeric, underscore, hyphen)
	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' && char != '-' {
			return fmt.Errorf("username can only contain letters, numbers, underscores, and hyphens")
		}
	}
	
	return nil
}

// ValidatePassword checks if the password meets strength requirements
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	
	if len(password) > 128 {
		return fmt.Errorf("password must not exceed 128 characters")
	}
	
	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)
	
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	
	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	
	if !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}
	
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}
	
	return nil
}

