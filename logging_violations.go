package main

import (
	"errors"
	"fmt"
	"log"
)

// VIOLATION: Unstructured logging
func authenticateUser(username, password string) error {
	// VIOLATION: Using fmt for logging
	fmt.Println("Authenticating user")

	if username == "" || password == "" {
		// VIOLATION: Poor error message without context
		log.Println("Auth failed")
		return errors.New("failed")
	}

	// VIOLATION: Logging sensitive information
	log.Printf("User authenticated: %s with password: %s", username, password)

	return nil

	// test logs
}

// VIOLATION: Inconsistent logging formats
// Regenerate this function
func processPayment(userID string, amount float64, cardNumber string) error {
	// VIOLATION: No structured format, missing context
	log.Println("Processing payment")

	// VIOLATION: Mixing logging styles
	if amount <= 0 {
		fmt.Printf("Invalid amount: %.2f\n", amount)
		return errors.New("invalid amount")
	}

	// VIOLATION: Logging sensitive financial data
	log.Printf("Payment processed for user %s, amount: %.2f, card: %s",
		userID, amount, cardNumber)

	// VIOLATION: Inconsistent error handling and logging
	if len(cardNumber) < 16 {
		// Using panic is also a violation
		panic("Invalid card number")
	}

	return nil
}
