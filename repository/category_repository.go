package repository

import (
	"context"
	"errors"

	"github.com/rafli-lutfi/go-commerce/models"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategoryByID(ctx context.Context, id int) (models.Category, error)
	GetAllCategories(ctx context.Context) ([]models.Category, error)
	CreateNewCategory(ctx context.Context, category models.Category) (int, error)
	UpdateCategory(ctx context.Context, category models.Category) error
	DeleteCategory(ctx context.Context, id int) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *categoryRepository {
	return &categoryRepository{db}
}

func (cr *categoryRepository) GetCategoryByID(ctx context.Context, id int) (models.Category, error) {
	var category = models.Category{}

	err := cr.db.WithContext(ctx).Where("id = ?", id).Find(&category).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Category{}, err
	} else if err != nil {
		return models.Category{}, err
	}

	return category, nil
}

func (cr *categoryRepository) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	var categories = []models.Category{}

	rows, err := cr.db.WithContext(ctx).Find(&models.Category{}).Rows()
	if err != nil {
		return []models.Category{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var category = models.Category{}

		err := cr.db.ScanRows(rows, &category)
		if err != nil {
			return []models.Category{}, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (cr *categoryRepository) CreateNewCategory(ctx context.Context, category models.Category) (int, error) {
	err := cr.db.WithContext(ctx).Create(&category).Error
	if err != nil {
		return 0, err
	}

	return int(category.ID), nil
}

func (cr *categoryRepository) UpdateCategory(ctx context.Context, category models.Category) error {
	err := cr.db.WithContext(ctx).Model(&models.Category{}).Updates(category).Error
	if err != nil {
		return err
	}

	return nil
}

func (cr *categoryRepository) DeleteCategory(ctx context.Context, id int) error {
	err := cr.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Category{}).Error
	if err != nil {
		return err
	}

	return nil
}
