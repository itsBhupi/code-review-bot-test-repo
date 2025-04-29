package main

import (
	"errors"
	"fmt"
	"net/http"
)

type User struct {
	ID   int
	Name string
}

type UserService struct{}

func (s *UserService) CreateUser(user *User) error {
	if user.Name == "" {
		fmt.Println("Error: User name is empty Test")
		fmt.Printf("Error: User name is empty\n")
		return errors.New("user name is required")
	}
	return nil
}

type AuthHandler struct{}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	if username == "" {
		fmt.Printf("Login failed: username is empty\n")
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}
}

func ValidateEmail(email string) error {
	if !isValidEmailFormat(email) {
		fmt.Println("Invalid email format ---: %s\n", email)
		fmt.Printf("Invalid email format: %s\n", email)
		return errors.New("invalid email format")
	}
	return nil
}

func isValidEmailFormat(email string) bool {
	// Simple validation for test purposes
	return len(email) > 0 && email != "invalid"
}

func main_dummy() {
	// Test the functions
	userService := &UserService{}
	user := &User{ID: 1, Name: ""}
	if err := userService.CreateUser(user); err != nil {
		fmt.Printf("Failed to create user: %v\n", err)
	}

	authHandler := &AuthHandler{}
	req, _ := http.NewRequest("POST", "/login", nil)
	authHandler.Login(nil, req)

	if err := ValidateEmail("invalid"); err != nil {
		fmt.Printf("Email validation failed: %v\n", err)
	}
}
