package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rafli-lutfi/go-commerce/controllers"
	"github.com/rafli-lutfi/go-commerce/middleware"
	"github.com/rafli-lutfi/go-commerce/repository"
	"github.com/rafli-lutfi/go-commerce/services"
	"gorm.io/gorm"
)

type APIHandler struct {
	productHandler  controllers.ProductHandler
	categoryHandler controllers.CategoryHandler
	discountHandler controllers.DiscountHandler
	userHandler     controllers.UserHandler
	orderHandler    controllers.OrderHandler
}

func RunServer(db *gorm.DB, r *gin.Engine) {
	productRepository := repository.NewProductRepository(db)
	categoryRepository := repository.NewCategoryRepository(db)
	discountRepository := repository.NewDiscountRepository(db)
	userRepository := repository.NewUserRepository(db)
	orderRepository := repository.NewOrderRepository(db)

	productService := services.NewProductService(productRepository, categoryRepository)
	categoryService := services.NewCategoryService(categoryRepository)
	discountService := services.NewDiscountService(discountRepository)
	userService := services.NewUserService(userRepository)
	orderService := services.NewOrderService(orderRepository)

	productHandler := controllers.NewProductHandler(productService)
	categoryHandler := controllers.NewCategoryHandler(categoryService)
	discountHandler := controllers.NewDiscountHandler(discountService)
	userHandler := controllers.NewUserHandler(userService)
	orderHandler := controllers.NewOrderHandler(orderService)

	apiHandler := APIHandler{
		productHandler:  productHandler,
		categoryHandler: categoryHandler,
		discountHandler: discountHandler,
		userHandler:     userHandler,
		orderHandler:    orderHandler,
	}

	api := r.Group("/api")

	products := api.Group("/product")
	products.GET("/:id", apiHandler.productHandler.GetProductByID)

	protectedProduct := products.Group("", middleware.Authentication())
	protectedProduct.POST("/create", apiHandler.productHandler.NewProduct)
	protectedProduct.PUT("/update", apiHandler.productHandler.UpdateProduct)
	protectedProduct.DELETE("/delete", apiHandler.productHandler.DeleteProduct) //query id

	categories := api.Group("/category")
	categories.GET("", apiHandler.productHandler.GetAllProductByCategory) //query name

	protectedCategory := categories.Group("", middleware.Authentication())
	protectedCategory.GET("/:id", apiHandler.categoryHandler.GetCategoryByID)
	protectedCategory.POST("/create", apiHandler.categoryHandler.CreateCategory)
	protectedCategory.PUT("/update", apiHandler.categoryHandler.UpdateCategory)
	protectedCategory.DELETE("/delete", apiHandler.categoryHandler.DeleteCategory) //query id

	discounts := api.Group("/discount", middleware.Authentication())
	discounts.GET("/:id", apiHandler.discountHandler.GetDiscountByID)       //get discountByID
	discounts.POST("/create", apiHandler.discountHandler.CreateNewDiscount) //Create New Discount
	discounts.PUT("/update", apiHandler.discountHandler.UpdateDiscount)     //Update exist discount
	discounts.DELETE("/delete", apiHandler.discountHandler.DeleteDiscount)  // Delete Discount with query id

	users := api.Group("/user")
	users.POST("/login", apiHandler.userHandler.Login)       // user login
	users.POST("/register", apiHandler.userHandler.Register) // add new user

	protectedUser := users.Group("", middleware.Authentication())
	protectedUser.GET("/:id", apiHandler.userHandler.GetUserByID) //Get user by id
	protectedUser.GET("/logout", apiHandler.userHandler.Logout)
	protectedUser.POST("/profile/newAddress", apiHandler.userHandler.AddNewAddress)   //update user
	protectedUser.PUT("/profile/update", apiHandler.userHandler.UpdateUser)           //update user
	protectedUser.PUT("/profile/updateAddress", apiHandler.userHandler.UpdateAddress) //update user

	orders := api.Group("/order", middleware.Authentication())
	orders.GET("/myOrder", apiHandler.orderHandler.ActiveOrder)
	orders.GET("/myOrder/:id", apiHandler.orderHandler.GetOrderByID)
	orders.GET("/myOrder/history", apiHandler.orderHandler.OrderHistory)
	orders.POST("/product", apiHandler.orderHandler.AddOrderItem)
	orders.POST("/payment/confirm", apiHandler.orderHandler.ConfirmOrder) //add payment
	orders.PUT("/update", apiHandler.orderHandler.UpdateOrder)            // Update Order
}
