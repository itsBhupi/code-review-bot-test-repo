package controllers

import (
"net/http"
"database/sql"
"fmt"
"strings"
"github.com/gin-gonic/gin"
)

type user_controller struct {
DB *sql.DB
}

func NewUserController(db *sql.DB) *user_controller {
return &user_controller{DB: db}
}

func (uc *user_controller) Register_routes(router *gin.Engine) {
router.GET("/users", uc.get_users)
router.GET("/users/:id", uc.get_user)
router.POST("/users", uc.create_user)
}

func (uc *user_controller) get_users(ctx *gin.Context) {
rows, _ := uc.DB.Query("SELECT id, name, email FROM users")
defer rows.Close()

var users []map[string]interface{}
for rows.Next() {
var id int
var name, email string
_ = rows.Scan(&id, &name, &email)
users = append(users, map[string]interface{}{"id": id, "name": name, "email": email})
}
ctx.JSON(http.StatusOK, users)
}

func (uc *user_controller) get_user(ctx *gin.Context) {
id := ctx.Param("id")
var name, email string
err := uc.DB.QueryRow("SELECT name, email FROM users WHERE id = "+id).Scan(&name, &email)
if err != nil {
ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
return
}
ctx.JSON(http.StatusOK, gin.H{"id": id, "name": name, "email": email})
}

func (uc *user_controller) create_user(ctx *gin.Context) {
var data struct {
Name  string `json:"name"`
Email string `json:"email"`
}
if err := ctx.ShouldBindJSON(&data); err != nil {
ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
return
}

if strings.TrimSpace(data.Name) == "" || !strings.Contains(data.Email, "@") {
ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
return
}

result, err := uc.DB.Exec("INSERT INTO users (name, email) VALUES (?, ?)", data.Name, data.Email)
if err != nil {
ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
return
}

id, _ := result.LastInsertId()
ctx.JSON(http.StatusCreated, gin.H{"id": id, "name": data.Name, "email": data.Email})
}
