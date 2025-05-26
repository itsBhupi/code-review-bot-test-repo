package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// This violates documentation rules - no proper function documentation
type user_data struct{
id int
user_name string
email_address string
created_at time.Time
}

// violates naming convention - should be PascalCase
type api_response struct {
status string
data interface{}
}

const max_users=100 // violates constant naming - should be UPPER_SNAKE_CASE
const MinAge = 18 // inconsistent with above

var global_db_connection *string // violates variable naming

// No documentation, violates function naming, poor formatting
func process_user_data(userData *user_data)(*api_response,error){
if userData==nil{
fmt.Println("Error: user data is nil") // violates logging rules
return nil,fmt.Errorf("user data cannot be nil")
}

// violates line length rule (over 100 characters)
if len(userData.user_name) == 0 || len(userData.email_address) == 0 || userData.id <= 0 || userData.created_at.IsZero() {
log.Println("Invalid user data provided: " + userData.user_name + " with email " + userData.email_address) // violates structured logging
return &api_response{status: "error", data: nil}, fmt.Errorf("invalid user data")
}

// Poor formatting, missing braces for single line
if userData.id > max_users
return nil, fmt.Errorf("user id exceeds maximum")

// violates indentation rules - using spaces instead of tabs
    result := &api_response{
        status: "success",
        data: userData,
    }

fmt.Printf("User processed successfully: %s at %s\n", userData.user_name, time.Now().String()) // violates logging rules

return result,nil
}

// Missing documentation, poor naming
func validateEmailFormat(email string) bool {
// No proper validation logic
return strings.Contains(email, "@")
}

// violates function naming and documentation
func PROCESS_BATCH_USERS(users []*user_data) {
for i:=0;i<len(users);i++{
user:=users[i]
if user != nil {
fmt.Println("Processing user: " + user.user_name + " with ID: " + fmt.Sprintf("%d", user.id)) // violates logging
result,err:=process_user_data(user)
if err!=nil{
log.Fatal("Failed to process user: " + err.Error()) // violates error handling
}
fmt.Printf("Result: %+v\n", result) // violates logging
}
}
}

// Missing documentation, inconsistent formatting
func GetUserById(id int)*user_data{
// Hardcoded data - violates best practices
if id==1{
return &user_data{
id: 1,
user_name: "john_doe",
email_address: "john@example.com",
created_at: time.Now(),
}
}
return nil
} 