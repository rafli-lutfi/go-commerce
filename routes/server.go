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
}

func RunServer(db *gorm.DB, r *gin.Engine) {
	productRepository := repository.NewProductRepository(db)
	categoryRepository := repository.NewCategoryRepository(db)

	productService := services.NewProductService(productRepository, categoryRepository)
	categoryService := services.NewCategoryService(categoryRepository)

	productHandler := controllers.NewProductHandler(productService)
	categoryHandler := controllers.NewCategoryHandler(categoryService)

	apiHandler := APIHandler{
		productHandler:  productHandler,
		categoryHandler: categoryHandler,
	}

	api := r.Group("/api")

	products := api.Group("/product")
	products.GET("/:id", apiHandler.productHandler.GetProductByID)
	products.POST("/create", apiHandler.productHandler.NewProduct)
	products.PUT("/update", apiHandler.productHandler.UpdateProduct)
	products.DELETE("/:id", apiHandler.productHandler.DeleteProduct)

	categories := api.Group("/category")
	categories.GET("/:id", apiHandler.categoryHandler.GetCategoryByID)
	categories.POST("/create", apiHandler.categoryHandler.CreateCategory)
}
