package controllers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// UserController handles HTTP requests for user operations
type UserController struct {
	DB *sql.DB
}

// NewUserController creates a new instance of UserController
func NewUserController(db *sql.DB) *UserController {
	return &UserController{DB: db}
}

// RegisterRoutes registers the routes for UserController
func (c *UserController) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.GET("", c.listUsers)
		users.GET("/:id", c.getUser)
		users.POST("", c.createUser)
	}
}

// listUsers handles GET /users
func (c *UserController) listUsers(ctx *gin.Context) {
	rows, err := c.DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		log.Error().Err(err).Msg("Failed to query users")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			log.Error().Err(err).Msg("Failed to scan user row")
			continue
		}
		users = append(users, map[string]interface{}{
			"id":    id,
			"name":  name,
			"email": email,
		})
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating user rows")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing users"})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// getUser handles GET /users/:id
func (c *UserController) getUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var name, email string
	err = c.DB.QueryRow("SELECT name, email FROM users WHERE id = $1", id).Scan(&name, &email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			log.Error().Err(err).Int("userID", id).Msg("Failed to query user")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":    id,
		"name":  name,
		"email": email,
	})
}

// createUser handles POST /users
func (c *UserController) createUser(ctx *gin.Context) {
	var input struct {
		Name  string `json:"name" binding:"required,min=2,max=100"`
		Email string `json:"email" binding:"required,email"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Additional validation
	if strings.TrimSpace(input.Name) == "" || !strings.Contains(input.Email, "@") {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Name and valid email are required"})
		return
	}

	// Start transaction
	tx, err := c.DB.Begin()
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	defer tx.Rollback()

	// Check if email already exists
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", input.Email).Scan(&count)
	if err != nil {
		log.Error().Err(err).Str("email", input.Email).Msg("Failed to check email existence")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	if count > 0 {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Insert new user
	var userID int
	err = tx.QueryRow(
		"INSERT INTO users (name, email, created_at) VALUES ($1, $2, $3) RETURNING id",
		input.Name,
		input.Email,
		time.Now(),
	).Scan(&userID)

	if err != nil {
		log.Error().Err(err).Msg("Failed to insert user")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Return created user
	ctx.JSON(http.StatusCreated, gin.H{
		"id":    userID,
		"name":  input.Name,
		"email": input.Email,
	})
}
