package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rafli-lutfi/go-commerce/controllers"
	"github.com/rafli-lutfi/go-commerce/repository"
	"github.com/rafli-lutfi/go-commerce/services"
	"gorm.io/gorm"
)

type APIHandler struct {
	productHandler  controllers.ProductHandler
	categoryHandler controllers.CategoryHandler
	discountHandler controllers.DiscountHandler
	userHandler     controllers.UserHandler
}

func RunServer(db *gorm.DB, r *gin.Engine) {
	productRepository := repository.NewProductRepository(db)
	categoryRepository := repository.NewCategoryRepository(db)
	discountRepository := repository.NewDiscountRepository(db)
	userRepository := repository.NewUserRepository(db)

	productService := services.NewProductService(productRepository, categoryRepository)
	categoryService := services.NewCategoryService(categoryRepository)
	discountService := services.NewDiscountService(discountRepository)
	userService := services.NewUserService(userRepository)

	productHandler := controllers.NewProductHandler(productService)
	categoryHandler := controllers.NewCategoryHandler(categoryService)
	discountHandler := controllers.NewDiscountHandler(discountService)
	userHandler := controllers.NewUserHandler(userService)

	apiHandler := APIHandler{
		productHandler:  productHandler,
		categoryHandler: categoryHandler,
		discountHandler: discountHandler,
		userHandler:     userHandler,
	}

	api := r.Group("/api")

	products := api.Group("/product")
	products.GET("/:id", apiHandler.productHandler.GetProductByID)
	products.POST("/create", apiHandler.productHandler.NewProduct)
	products.PUT("/update", apiHandler.productHandler.UpdateProduct)
	products.DELETE("/delete", apiHandler.productHandler.DeleteProduct) //query id

	categories := api.Group("/category")
	categories.GET("/:id", apiHandler.categoryHandler.GetCategoryByID)
	categories.GET("", apiHandler.productHandler.GetAllProductByCategory) //query name
	categories.POST("/create", apiHandler.categoryHandler.CreateCategory)
	categories.PUT("/update", apiHandler.categoryHandler.UpdateCategory)
	categories.DELETE("/delete", apiHandler.categoryHandler.DeleteCategory) //query id

	discounts := api.Group("/discount")
	discounts.GET("/:id", apiHandler.discountHandler.GetDiscountByID)       //get discountByID
	discounts.POST("/create", apiHandler.discountHandler.CreateNewDiscount) //Create New Discount
	discounts.PUT("/update", apiHandler.discountHandler.UpdateDiscount)     //Update exist discount
	discounts.DELETE("/delete", apiHandler.discountHandler.DeleteDiscount)  // Delete Discount with query id

	users := api.Group("/user")
	users.GET("/:id", apiHandler.userHandler.GetUserByID)             //Get user by id
	users.POST("/login")                                              // user login
	users.POST("/register", apiHandler.userHandler.Register)          // add new user
	users.POST("/newAddress", apiHandler.userHandler.AddNewAddress)   //update user
	users.PUT("/update", apiHandler.userHandler.UpdateUser)           //update user
	users.PUT("/updateAddress", apiHandler.userHandler.UpdateAddress) //update user
}
