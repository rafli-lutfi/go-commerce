package services

import (
	"context"
	"errors"

	"github.com/rafli-lutfi/go-commerce/models"
	"github.com/rafli-lutfi/go-commerce/repository"
)

type ProductService interface {
	GetProductWithJoinByID(ctx context.Context, id int) (models.ProductInfo, error)
	CreateNewProduct(ctx context.Context, product models.Product) (int, error)
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, productID int) error
	GetAllProductByCategories(ctx context.Context, categoryName string) ([]models.ProductInfo, error)
}

type productService struct {
	productRepository  repository.ProductRepository
	categoryRepository repository.CategoryRepository
}

func NewProductService(productRepository repository.ProductRepository, categoryRepository repository.CategoryRepository) *productService {
	return &productService{productRepository, categoryRepository}
}

func (ps *productService) GetProductWithJoinByID(ctx context.Context, id int) (models.ProductInfo, error) {
	productDB, err := ps.productRepository.GetProductByID(ctx, id)
	if err != nil {
		return models.ProductInfo{}, err
	}
	if productDB.ID == 0 {
		return models.ProductInfo{}, errors.New("product not found")
	}

	product, err := ps.productRepository.JoinTableProduct(ctx, int(productDB.ID))
	if err != nil {
		return models.ProductInfo{}, err
	}

	return product, nil
}

func (ps *productService) CreateNewProduct(ctx context.Context, product models.Product) (int, error) {
	// Checking Category ID
	_, err := ps.categoryRepository.GetCategoryByID(ctx, int(product.CategoryID))
	if err != nil {
		return 0, errors.New("id_category is not exist")
	}

	// Create Product
	productID, err := ps.productRepository.CreateNewProduct(ctx, product)
	if err != nil {
		return 0, err
	}

	return productID, nil
}

func (ps *productService) UpdateProduct(ctx context.Context, product *models.Product) error {
	// Checking Product ID
	productDB, err := ps.productRepository.GetProductByID(ctx, int(product.ID))
	if err != nil {
		return err
	}

	if productDB.ID == 0 {
		return errors.New("id product not found")
	}

	// Update Product
	err = ps.productRepository.UpdateProduct(ctx, product)
	if err != nil {
		return err
	}

	return nil
}

func (ps *productService) DeleteProduct(ctx context.Context, productID int) error {
	// Checking Product ID
	productDB, err := ps.productRepository.GetProductByID(ctx, productID)
	if err != nil {
		return err
	}

	if productDB.ID == 0 {
		return errors.New("product not found")
	}

	// Delete Product
	err = ps.productRepository.DeleteProduct(ctx, productID)
	if err != nil {
		return err
	}

	return nil
}

func (ps *productService) GetAllProductByCategories(ctx context.Context, categoryName string) ([]models.ProductInfo, error) {
	// Checking Category Name
	category, err := ps.categoryRepository.GetCategoryByName(ctx, categoryName)
	if err != nil {
		return []models.ProductInfo{}, err
	}

	if category.ID == 0 {
		return []models.ProductInfo{}, errors.New("category name not found")
	}

	// Find Products By Category
	products, err := ps.productRepository.GetAllProductByCategories(ctx, int(category.ID))
	if err != nil {
		return []models.ProductInfo{}, err
	}

	return products, nil
}
