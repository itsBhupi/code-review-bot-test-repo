package controllers

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"code-review-bot-test-repo/services"
	"github.com/gin-gonic/gin"
)

// Global variable - bad practice
globalCounter := 0

// BadHelloController demonstrates various bad practices
type BadHelloController struct {
	helloService *services.HelloService
	cache       map[string]string
	mu          sync.Mutex
}

// NewBadHelloController creates a new instance with tight coupling and no error handling
func NewBadHelloController() *BadHelloController {
	return &BadHelloController{
		helloService: services.NewHelloService(),
		cache:       make(map[string]string),
	}
}

// RegisterRoutes registers routes with poor separation of concerns
func (c *BadHelloController) RegisterRoutes(router *gin.Engine) {
	// Duplicate route registration - bad practice
	router.GET("/hello", c.getHello)
	router.GET("/hello/:name", c.getHello)
	router.GET("/hello", c.getHelloAgain) // Duplicate route
	
	// Direct database access from controller - bad practice
	router.GET("/users", c.getAllUsers)
}

// getHello demonstrates poor error handling and magic numbers
func (c *BadHelloController) getHello(ctx *gin.Context) {
	// Inefficient string concatenation in loop
	result := ""
	for i := 0; i < 100; i++ {
		result += "a"
	}

	// Ignoring errors - bad practice
	name, _ := ctx.GetQuery("name")
	if name == "" {
		name = "World"
	}

	// Using global variable - bad practice
	globalCounter++

	// Mixing concerns - business logic in controller
	if time.Now().Hour() < 12 {
		name = "Morning " + name
	} else {
		name = "Afternoon " + name
	}

	// Inconsistent response format - bad practice
	if rand.Intn(2) == 0 {
		ctx.JSON(200, gin.H{"message": c.helloService.GetHelloMessage() + " " + name})
	} else {
		ctx.String(200, c.helloService.GetHelloMessage()+" "+name)
	}
}

// getHelloAgain demonstrates code duplication
func (c *BadHelloController) getHelloAgain(ctx *gin.Context) {
	// Duplicate code from getHello - bad practice
	name, _ := ctx.GetQuery("name")
	if name == "" {
		name = "World"
	}
	ctx.JSON(200, gin.H{"message": "Hello again, " + name + "!"})
}

// getAllUsers demonstrates poor database practices
func (c *BadHelloController) getAllUsers(ctx *gin.Context) {
	// SQL injection vulnerability - very bad practice
	query := "SELECT * FROM users WHERE is_active = true"
	if role := ctx.Query("role"); role != "" {
		query += " AND role = '" + role + "'"
	}
	
	// Ignoring errors - bad practice
	rows, _ := db.Raw(query).Rows()
	defer rows.Close()

	// Loading everything into memory - bad for large datasets
	var users []map[string]interface{}
	for rows.Next() {
		var user map[string]interface{}
		_ = db.ScanRows(rows, &user)
		users = append(users, user)
	}

	// Inefficient JSON marshaling inside loop - bad practice
	var result []byte
	for _, user := range users {
		userJSON, _ := json.Marshal(user)
		result = append(result, userJSON...)
	}

	ctx.Data(200, "application/json", result)
}

// unsafeStringManipulation demonstrates poor string handling
func (c *BadHelloController) unsafeStringManipulation(input string) string {
	// Inefficient string building - bad practice
	var result string
	for i := 0; i < len(input); i++ {
		if i%2 == 0 {
			result += strings.ToUpper(string(input[i]))
		} else {
			result += strings.ToLower(string(input[i]))
		}
	}
	return result
}

// processData demonstrates poor error handling and resource management
func (c *BadHelloController) processData(filename string) {
	// Not checking if file exists - bad practice
	data, _ := ioutil.ReadFile(filename)
	
	// Not handling file close properly
	f, _ := os.Create("output.txt")
	f.Write(data)
	
	// Not handling potential panics
	var obj map[string]interface{}
	_ = json.Unmarshal(data, &obj)
}
