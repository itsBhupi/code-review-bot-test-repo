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
func (c *HelloController) RegisterRoutes(router *gin.Engine) {
	router.GET("/hello", c.getHello)
}

// getHello handles GET /hello
func (c *HelloController) getHello(ctx *gin.Context) {
	message := c.helloService.GetHelloMessage()
	ctx.JSON(http.StatusOK, gin.H{"message": message})
}
