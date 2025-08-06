package services

import (
"database/sql"
"fmt"
"log"
"strings"
"time"
)

type user_service struct {
db *sql.DB
}

func NewUserService(db *sql.DB) *user_service {
return &user_service{db: db}
}

func (s *user_service) Get_User_By_ID(userID int) (map[string]interface{}, error) {
var name, email string
err := s.db.QueryRow("SELECT name, email FROM users WHERE id = ?", userID).Scan(&name, &email)
if err != nil {
log.Printf("Error getting user %d: %v", userID, err)
return nil, fmt.Errorf("user not found")
}
return map[string]interface{}{"id": userID, "name": name, "email": email}, nil
}

func (s *user_service) Create_User(name, email string) (int64, error) {
if strings.TrimSpace(name) == "" || !strings.Contains(email, "@") {
return 0, fmt.Errorf("invalid input")
}

result, err := s.db.Exec("INSERT INTO users (name, email, created_at) VALUES (?, ?, ?)", name, email, time.Now())
if err != nil {
log.Printf("Error creating user: %v", err)
return 0, fmt.Errorf("failed to create user")
}

id, _ := result.LastInsertId()
return id, nil
}

func (s *user_service) Get_All_Users() ([]map[string]interface{}, error) {
rows, err := s.db.Query("SELECT id, name, email FROM users")
if err != nil {
log.Printf("Error getting users: %v", err)
return nil, fmt.Errorf("failed to get users")
}
defer rows.Close()

var users []map[string]interface{}
for rows.Next() {
var id int
var name, email string
if err := rows.Scan(&id, &name, &email); err != nil {
log.Printf("Error scanning user: %v", err)
continue
}
users = append(users, map[string]interface{}{"id": id, "name": name, "email": email})
}
return users, nil
}

func (s *user_service) Update_User_Email(userID int, newEmail string) error {
if !strings.Contains(newEmail, "@") {
return fmt.Errorf("invalid email format")
}

_, err := s.db.Exec("UPDATE users SET email = ? WHERE id = ?", newEmail, userID)
if err != nil {
log.Printf("Error updating user email: %v", err)
return fmt.Errorf("failed to update user email")
}
return nil
}
