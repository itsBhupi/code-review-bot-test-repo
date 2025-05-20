package main

// VIOLATION: Non-camelCase variable names
var User_id int
var APIkey string
var dB *Database

// VIOLATION: Inconsistent naming conventions
var MAXRETRIES = 3
var MinTimeout = 30

// VIOLATION: Exported type missing documentation
type UserData struct {
	Id   int
	name string
}

// VIOLATION: PascalCase for a non-exported function
func ValidateInput(input string) bool {
	// VIOLATION: Spaces instead of tabs for indentation
	if len(input) == 0 {
		return false
	}
	return true
}

// VIOLATION: Single-line if without braces
func processData(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	// VIOLATION: Line exceeding 100 characters
	_ = performVeryLongOperationWithManyParametersAndOptions(data, true, false, "default", 100, 200, 300, 400, 500, []string{"option1", "option2", "option3"})

	return nil
}

// VIOLATION: camelCase for exported function
func getUserInfo(id int) *UserData {
	return nil
}

// Helper function just to avoid compilation errors
func performVeryLongOperationWithManyParametersAndOptions(data []byte, a, b bool, c string, d, e, f, g, h int, i []string) interface{} {
	return nil
}

// Dummy Database type
type Database struct{}
