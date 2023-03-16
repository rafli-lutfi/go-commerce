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
	products.GET("/:id", middleware.Authentication(), apiHandler.productHandler.GetProductByID)
	products.POST("/create", middleware.Authentication(), apiHandler.productHandler.NewProduct)
	products.PUT("/update", middleware.Authentication(), apiHandler.productHandler.UpdateProduct)
	products.DELETE("/delete", middleware.Authentication(), apiHandler.productHandler.DeleteProduct) //query id

	categories := api.Group("/category")
	categories.GET("/:id", middleware.Authentication(), apiHandler.categoryHandler.GetCategoryByID)
	categories.GET("", middleware.Authentication(), apiHandler.productHandler.GetAllProductByCategory) //query name
	categories.POST("/create", middleware.Authentication(), apiHandler.categoryHandler.CreateCategory)
	categories.PUT("/update", middleware.Authentication(), apiHandler.categoryHandler.UpdateCategory)
	categories.DELETE("/delete", middleware.Authentication(), apiHandler.categoryHandler.DeleteCategory) //query id

	discounts := api.Group("/discount")
	discounts.GET("/:id", middleware.Authentication(), apiHandler.discountHandler.GetDiscountByID)       //get discountByID
	discounts.POST("/create", middleware.Authentication(), apiHandler.discountHandler.CreateNewDiscount) //Create New Discount
	discounts.PUT("/update", middleware.Authentication(), apiHandler.discountHandler.UpdateDiscount)     //Update exist discount
	discounts.DELETE("/delete", middleware.Authentication(), apiHandler.discountHandler.DeleteDiscount)  // Delete Discount with query id

	users := api.Group("/user")
	users.GET("/:id", middleware.Authentication(), apiHandler.userHandler.GetUserByID) //Get user by id
	users.POST("/login", apiHandler.userHandler.Login)                                 // user login
	users.POST("/register", apiHandler.userHandler.Register)                           // add new user
	users.GET("/logout", apiHandler.userHandler.Logout)
	users.POST("/profile/newAddress", middleware.Authentication(), apiHandler.userHandler.AddNewAddress)   //update user
	users.PUT("/profile/update", middleware.Authentication(), apiHandler.userHandler.UpdateUser)           //update user
	users.PUT("/profile/updateAddress", middleware.Authentication(), apiHandler.userHandler.UpdateAddress) //update user

	orders := api.Group("/order")
	orders.GET("/myOrder/:id", middleware.Authentication(), apiHandler.orderHandler.GetOrderByID)
	orders.GET("/myOrder/item/:id", middleware.Authentication(), apiHandler.orderHandler.GetOrderItemByID)
	orders.GET("/myOrder/history", middleware.Authentication())
	orders.POST("/product", middleware.Authentication(), apiHandler.orderHandler.AddOrderItem)
	orders.POST("/payment", middleware.Authentication()) //add payment
	orders.PUT("/update", middleware.Authentication())   // Update Order
}
