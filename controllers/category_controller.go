package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rafli-lutfi/go-commerce/models"
	"github.com/rafli-lutfi/go-commerce/services"
)

type CategoryHandler interface {
	CreateCategory(c *gin.Context)
	GetCategoryByID(c *gin.Context)
}

type categoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) *categoryHandler {
	return &categoryHandler{categoryService}
}

func (ch *categoryHandler) CreateCategory(c *gin.Context) {
	var category = models.Category{}

	err := c.ShouldBindJSON(&category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if category.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "category name is empty",
			"data":    nil,
		})
		return
	}

	categoryID, err := ch.categoryService.CreateNewCategory(c.Request.Context(), category)
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
		"message": "success create new category",
		"data": gin.H{
			"category_id": categoryID,
		},
	})
}

func (ch *categoryHandler) GetCategoryByID(c *gin.Context) {
	var ctx = c.Request.Context()

	id := c.Param("id")
	categoryID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "invalid id type",
			"data":    nil,
		})
		return
	}

	category, err := ch.categoryService.GetCategoryByID(ctx, categoryID)
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
		"message": "success get category",
		"data":    category,
	})
}
