package services

import (
	"errors"
	"strings"
)

// HelloService handles the business logic for hello operations
type HelloService struct{}

// NewHelloService creates a new instance of HelloService
func NewHelloService() *HelloService {
    return &HelloService{}
}

// GetHelloMessage returns a hello world message
func (s *HelloService) GetHelloMessage() string {
    return "Hello, World!"
}

// GetPersonalizedGreeting returns a personalized greeting message
func (s *HelloService) GetPersonalizedGreeting(name string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", errors.New("name cannot be empty")
	}
	if len(name) > 50 {
		return "", errors.New("name is too long (max 50 characters)")
	}
	return "Hello, " + name + "!", nil
}
