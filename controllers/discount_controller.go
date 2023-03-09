package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rafli-lutfi/go-commerce/models"
	"github.com/rafli-lutfi/go-commerce/services"
)

type DiscountHandler interface {
	GetDiscountByID(c *gin.Context)
	CreateNewDiscount(c *gin.Context)
	UpdateDiscount(c *gin.Context)
	DeleteDiscount(c *gin.Context)
}

type discountHandler struct {
	discountService services.DiscountService
}

func NewDiscountHandler(discountService services.DiscountService) *discountHandler {
	return &discountHandler{discountService}
}

func (dh *discountHandler) GetDiscountByID(c *gin.Context) {
	var ctx = c.Request.Context()
	var id = c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "query param is empty",
			"data":    nil,
		})
		return
	}

	discountID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "invalid id type",
			"data":    nil,
		})
		return
	}

	discount, err := dh.discountService.GetDiscountByID(ctx, discountID)
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
		"message": "success get discount",
		"data":    discount,
	})
}

func (dh *discountHandler) CreateNewDiscount(c *gin.Context) {
	var ctx = c.Request.Context()
	var discount = models.Discount{}

	err := c.ShouldBindJSON(&discount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "something wrong with your request body",
			"data":    nil,
		})
		return
	}

	discountID, err := dh.discountService.CreateDiscount(ctx, discount)
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
		"message": "success create new discount",
		"data": gin.H{
			"discount_id": discountID,
		},
	})
}

func (dh *discountHandler) UpdateDiscount(c *gin.Context) {
	var ctx = c.Request.Context()
	var product = models.Discount{}

	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "error body request",
			"data":    nil,
		})
		return
	}

	err = dh.discountService.UpdateDiscount(ctx, &product)
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
		"message": "success update discount",
		"data":    "",
	})
}

func (dh *discountHandler) DeleteDiscount(c *gin.Context) {
	var ctx = c.Request.Context()
	var id = c.Query("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "id query param is empty",
			"data":    nil,
		})
		return
	}

	discountID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "invalid id type",
			"data":    nil,
		})
		return
	}

	err = dh.discountService.DeleteDiscount(ctx, discountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"staus":   http.StatusOK,
		"message": "success delete discount",
		"data":    "",
	})
}
