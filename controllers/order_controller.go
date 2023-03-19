package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rafli-lutfi/go-commerce/models"
	"github.com/rafli-lutfi/go-commerce/services"
)

type OrderHandler interface {
	GetOrderByID(c *gin.Context)
	AddOrderItem(c *gin.Context)
	UpdateOrder(c *gin.Context)
	ConfirmOrder(c *gin.Context)
	OrderHistory(c *gin.Context)
	ActiveOrder(c *gin.Context)
}

type orderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) *orderHandler {
	return &orderHandler{orderService}
}

func (oh *orderHandler) GetOrderByID(c *gin.Context) {
	var ctx = c.Request.Context()
	var id = c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "id order is empty",
		})
		return
	}

	orderID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "id is not integer",
		})
		return
	}

	order, err := oh.orderService.GetOrderByID(ctx, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success get order",
		"data":    order,
	})
}

func (oh *orderHandler) AddOrderItem(c *gin.Context) {
	var ctx = c.Request.Context()
	var id, _ = c.Get("userID")
	userIDInt := id.(int)

	orderItem := models.OrderItem{}

	err := c.ShouldBindJSON(&orderItem)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "something wrong with the body request",
		})
		return
	}

	result, err := oh.orderService.CreateNewOrder(ctx, orderItem, userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order_detail": result,
	})
}

func (oh *orderHandler) UpdateOrder(c *gin.Context) {
	var ctx = c.Request.Context()
	updateItem := models.UpdateItem{}

	err := c.ShouldBindJSON(&updateItem)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "something wrong with the body request",
		})
		return
	}

	err = oh.orderService.UpdateOrder(ctx, updateItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "update success",
	})
}

func (oh *orderHandler) ConfirmOrder(c *gin.Context) {
	var ctx = c.Request.Context()
	var payment = models.Payment{}

	err := c.ShouldBindJSON(&payment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "something wrong with the body request",
		})
		return
	}

	if payment.ID == 0 || payment.Amount == 0 || payment.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "payment amount or name is empty",
		})
		return
	}

	order, err := oh.orderService.ConfirmPayment(ctx, payment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "payment success",
		"data":    order,
	})
}

func (oh *orderHandler) OrderHistory(c *gin.Context) {
	var ctx = c.Request.Context()
	var id, _ = c.Get("userID")
	userIDInt := id.(int)

	orders, err := oh.orderService.OrderHistory(ctx, userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success get orders",
		"data":    orders,
	})
}

func (oh *orderHandler) ActiveOrder(c *gin.Context) {
	var ctx = c.Request.Context()
	var id, _ = c.Get("userID")
	userIDInt := id.(int)

	order, err := oh.orderService.CurrentOrder(ctx, userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success get order",
		"data":    order,
	})
}
