package controllers

import (
	"net/http"

	"code-review-bot-test-repo/services"
	"github.com/gin-gonic/gin"
)

// HelloController handles HTTP requests for hello operations
type HelloController struct {
	helloService *services.HelloService
}

// NewHelloController creates a new instance of HelloController
func NewHelloController(helloService *services.HelloService) *HelloController {
	return &HelloController{
		helloService: helloService,
	}
}

// RegisterRoutes registers the routes for HelloController
func (c *HelloController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/hello", c.getHello)
	router.GET("/hello/:name", c.getPersonalizedHello)
}

// getHello handles GET /hello
func (c *HelloController) getHello(ctx *gin.Context) {
	message := c.helloService.GetHelloMessage()
	ctx.JSON(http.StatusOK, gin.H{"message": message})
}

// getPersonalizedHello handles GET /hello/:name
func (c *HelloController) getPersonalizedHello(ctx *gin.Context) {
	name := ctx.Param("name")
	
	message, err := c.helloService.GetPersonalizedGreeting(name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"greeting": message,
	})
}
