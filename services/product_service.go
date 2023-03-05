package services

import (
	"context"
	"errors"

	"github.com/rafli-lutfi/go-commerce/models"
	"github.com/rafli-lutfi/go-commerce/repository"
)

type ProductService interface {
	GetProductByID(ctx context.Context, id int) (models.Product, error)
	CreateNewProduct(ctx context.Context, product models.NewProduct) (int, error)
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, productID int) error
	GetAllProductByCategories(ctx context.Context, categoryID int) ([]models.Product, error)
}

type productService struct {
	productRepository  repository.ProductRepository
	categoryRepository repository.CategoryRepository
}

func NewProductService(productRepository repository.ProductRepository, categoryRepository repository.CategoryRepository) *productService {
	return &productService{productRepository, categoryRepository}
}

func (ps *productService) GetProductByID(ctx context.Context, id int) (models.Product, error) {
	product, err := ps.productRepository.GetProductByID(ctx, id)
	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func (ps *productService) CreateNewProduct(ctx context.Context, product models.NewProduct) (int, error) {
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
	_, err := ps.productRepository.GetProductByID(ctx, int(product.ID))
	if err != nil {
		return err
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
	_, err := ps.productRepository.GetProductByID(ctx, productID)
	if err != nil {
		return err
	}

	// Delete Product
	err = ps.productRepository.DeleteProduct(ctx, productID)
	if err != nil {
		return err
	}

	return nil
}

func (ps *productService) GetAllProductByCategories(ctx context.Context, categoryID int) ([]models.Product, error) {
	// Checking Category ID
	_, err := ps.categoryRepository.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return []models.Product{}, err
	}

	// Find Products By Category
	products, err := ps.productRepository.GetAllProductByCategories(ctx, categoryID)
	if err != nil {
		return []models.Product{}, err
	}

	return products, nil
}
