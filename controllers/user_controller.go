package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rafli-lutfi/go-commerce/models"
	"github.com/rafli-lutfi/go-commerce/services"
)

type UserHandler interface {
	GetUserByID(c *gin.Context)
	Register(c *gin.Context)
	AddNewAddress(c *gin.Context)
	UpdateUser(c *gin.Context)
	UpdateAddress(c *gin.Context)
}

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *userHandler {
	return &userHandler{userService}
}

func (uh *userHandler) GetUserByID(c *gin.Context) {
	var ctx = c.Request.Context()
	var id = c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "there is no id in param",
			"data":    nil,
		})
		return
	}

	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "invalid id type",
			"data":    nil,
		})
		return
	}

	user, err := uh.userService.GetUserByID(ctx, userID)
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
		"message": "success get user",
		"data":    user,
	})

}

func (uh *userHandler) Login(c *gin.Context) {
	var ctx = c.Request.Context()
	var form models.Login

	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "request body is empty",
			"data":    nil,
		})
		return
	}

	userID, err := uh.userService.Login(ctx, &form)
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
		"message": "success login",
		"data":    "",
	})
}

func (uh *userHandler) Register(c *gin.Context) {
	var ctx = c.Request.Context()
	var form models.NewUser

	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "request body is empty",
			"data":    nil,
		})
		return
	}

	userID, err := uh.userService.Register(ctx, form)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":   http.StatusInternalServerError,
			"messagge": err.Error(),
			"data":     nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "success create new user",
		"data": gin.H{
			"user_id": userID,
		},
	})
}

func (uh *userHandler) AddNewAddress(c *gin.Context) {
	var ctx = c.Request.Context()
	var form models.NewAddress

	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "request body is empty",
			"data":    nil,
		})
		return
	}

	err = uh.userService.AddNewAddress(ctx, form)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":   http.StatusInternalServerError,
			"messagge": err.Error(),
			"data":     nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "success add new user",
		"data":    "",
	})
}

func (uh *userHandler) UpdateUser(c *gin.Context) {
	var ctx = c.Request.Context()
	var form models.User

	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "request body is empty",
			"data":    nil,
		})
		return
	}

	err = uh.userService.UpdateUser(ctx, &form)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":   http.StatusInternalServerError,
			"messagge": err.Error(),
			"data":     nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success update user",
		"data":    "",
	})
}

func (uh *userHandler) UpdateAddress(c *gin.Context) {
	var ctx = c.Request.Context()
	var form models.UserAddress

	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "request body is empty",
			"data":    nil,
		})
		return
	}

	err = uh.userService.UpdateUserAddress(ctx, &form)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":   http.StatusInternalServerError,
			"messagge": err.Error(),
			"data":     nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success update address",
		"data":    "",
	})
}
