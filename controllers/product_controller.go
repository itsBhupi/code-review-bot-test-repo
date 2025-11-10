package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"code-review-bot-test-repo/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ProductController handles HTTP requests for product operations
type ProductController struct {
	productService *services.ProductService
}

// NewProductController creates a new instance of ProductController
func NewProductController(productService *services.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

// RegisterRoutes registers the routes for ProductController
func (c *ProductController) RegisterRoutes(router *gin.RouterGroup) {
	products := router.Group("/products")
	{
		products.GET("", c.listProducts)
		products.GET("/:id", c.getProduct)
		products.GET("/sku/:sku", c.getProductBySKU)
		products.POST("", c.createProduct)
		products.PUT("/:id", c.updateProduct)
		products.PATCH("/:id", c.updateProduct)
		products.DELETE("/:id", c.deleteProduct)
	}
}

// listProducts handles GET /products
func (c *ProductController) listProducts(ctx *gin.Context) {
	// Parse query parameters for pagination
	limitStr := ctx.DefaultQuery("limit", "50")
	offsetStr := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	products, err := c.productService.ListProducts(ctx.Request.Context(), limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list products")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve products",
		})
		return
	}

	// Return empty array instead of null if no products
	if products == nil {
		products = []services.Product{}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products": products,
		"limit":    limit,
		"offset":   offset,
		"count":    len(products),
	})
}

// getProduct handles GET /products/:id
func (c *ProductController) getProduct(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	product, err := c.productService.GetProductByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
		} else {
			log.Error().Err(err).Int("productID", id).Msg("Failed to get product")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve product",
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// getProductBySKU handles GET /products/sku/:sku
func (c *ProductController) getProductBySKU(ctx *gin.Context) {
	sku := ctx.Param("sku")
	if sku == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "SKU is required",
		})
		return
	}

	product, err := c.productService.GetProductBySKU(ctx.Request.Context(), sku)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
		} else {
			log.Error().Err(err).Str("sku", sku).Msg("Failed to get product by SKU")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve product",
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// createProduct handles POST /products
func (c *ProductController) createProduct(ctx *gin.Context) {
	var input services.CreateProductInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input: " + err.Error(),
		})
		return
	}

	product, err := c.productService.CreateProduct(ctx.Request.Context(), input)
	if err != nil {
		if errors.Is(err, services.ErrProductAlreadyExists) {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "Product with this SKU already exists",
			})
		} else if errors.Is(err, services.ErrInvalidProductData) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			log.Error().Err(err).Msg("Failed to create product")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create product",
			})
		}
		return
	}

	ctx.JSON(http.StatusCreated, product)
}

// updateProduct handles PUT/PATCH /products/:id
func (c *ProductController) updateProduct(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	var input services.UpdateProductInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input: " + err.Error(),
		})
		return
	}

	product, err := c.productService.UpdateProduct(ctx.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
		} else if errors.Is(err, services.ErrInvalidProductData) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			log.Error().Err(err).Int("productID", id).Msg("Failed to update product")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update product",
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// deleteProduct handles DELETE /products/:id
func (c *ProductController) deleteProduct(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	err = c.productService.DeleteProduct(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
		} else {
			log.Error().Err(err).Int("productID", id).Msg("Failed to delete product")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to delete product",
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}

