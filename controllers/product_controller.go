package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rafli-lutfi/go-commerce/models"
	"github.com/rafli-lutfi/go-commerce/services"
)

type ProductHandler interface {
	NewProduct(c *gin.Context)
	GetProductByID(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
}

type productHandler struct {
	productService services.ProductService
}

func NewProductHandler(productService services.ProductService) *productHandler {
	return &productHandler{productService}
}

func (ph *productHandler) NewProduct(c *gin.Context) {
	var product = models.Product{}
	var ctx = c.Request.Context()

	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if product.Name == "" || product.Price == 0 || product.DiscountID == 0 || product.CategoryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "there is data/datas that should not empty",
			"data":    nil,
		})
		return
	}

	productID, err := ph.productService.CreateNewProduct(ctx, product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "success create new product",
		"data": gin.H{
			"id": productID,
		},
	})
}

func (ph *productHandler) GetProductByID(c *gin.Context) {
	var ctx = c.Request.Context()
	var id = c.Param("id")

	productID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "invalid id type",
			"data":    nil,
		})
		return
	}

	product, err := ph.productService.GetProductByID(ctx, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success get product",
		"data":    product,
	})
}

func (ph *productHandler) UpdateProduct(c *gin.Context) {
	var ctx = c.Request.Context()
	var product = models.Product{}

	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	err = ph.productService.UpdateProduct(ctx, &product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success update data product",
		"data":    "",
	})
}

func (ph *productHandler) DeleteProduct(c *gin.Context) {
	var ctx = c.Request.Context()
	var id = c.Param("id")

	productID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "invalid id type",
			"data":    nil,
		})
		return
	}

	err = ph.productService.DeleteProduct(ctx, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success delete product",
		"data":    "",
	})
}
