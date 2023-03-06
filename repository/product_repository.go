package repository

import (
	"context"
	"errors"

	"github.com/rafli-lutfi/go-commerce/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetProductByID(ctx context.Context, id int) (models.Product, error)
	CreateNewProduct(ctx context.Context, newProduct models.Product) (int, error)
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, id int) error
	GetAllProductByCategories(ctx context.Context, idCategory int) ([]models.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{db}
}

func (pr *productRepository) GetProductByID(ctx context.Context, id int) (models.Product, error) {
	product := models.Product{}

	err := pr.db.WithContext(ctx).Where("id = ?", id).First(&product).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Product{}, nil
	} else if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func (pr *productRepository) CreateNewProduct(ctx context.Context, product models.Product) (int, error) {
	// Create product record
	if err := pr.db.WithContext(ctx).Create(&product).Error; err != nil {
		return 0, err
	}

	return int(product.ID), nil
}

func (pr *productRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	err := pr.db.WithContext(ctx).Model(&models.Product{}).Where("id = ?", product.ID).Updates(product).Error
	if err != nil {
		return err
	}

	return nil
}

func (pr *productRepository) DeleteProduct(ctx context.Context, id int) error {
	err := pr.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Product{}).Error
	if err != nil {
		return err
	}

	return nil
}

// get all product by categories
func (pr *productRepository) GetAllProductByCategories(ctx context.Context, idCategory int) ([]models.Product, error) {
	productsByCategory := []models.Product{}

	rows, err := pr.db.WithContext(ctx).Where("id_category = ?", idCategory).Find(&models.Category{}).Rows()
	if err != nil {
		return []models.Product{}, err
	}
	defer rows.Close()

	for rows.Next() {
		product := models.Product{}

		err := pr.db.WithContext(ctx).ScanRows(rows, &product)
		if err != nil {
			return []models.Product{}, err
		}

		productsByCategory = append(productsByCategory, product)
	}

	return productsByCategory, nil
}
