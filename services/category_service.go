package services

import (
	"context"

	"github.com/rafli-lutfi/go-commerce/models"
	"github.com/rafli-lutfi/go-commerce/repository"
)

type CategoryService interface {
	GetCategoryByID(ctx context.Context, categoryID int) (models.Category, error)
	GetAllCategories(ctx context.Context) ([]models.Category, error)
	CreateNewCategory(ctx context.Context, category models.Category) (int, error)
	UpdateCategory(ctx context.Context, category models.Category) error
	DeleteCategory(ctx context.Context, categoryID int) error
}

type categoryService struct {
	categoryRepository repository.CategoryRepository
}

func NewCategoryService(categoryRepository repository.CategoryRepository) *categoryService {
	return &categoryService{categoryRepository}
}

func (cs *categoryService) GetCategoryByID(ctx context.Context, categoryID int) (models.Category, error) {
	category, err := cs.categoryRepository.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return models.Category{}, err
	}

	return category, err
}

func (cs *categoryService) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	categories, err := cs.categoryRepository.GetAllCategories(ctx)
	if err != nil {
		return []models.Category{}, err
	}

	return categories, nil
}

func (cs *categoryService) CreateNewCategory(ctx context.Context, category models.Category) (int, error) {
	categoryID, err := cs.categoryRepository.CreateNewCategory(ctx, category)
	if err != nil {
		return 0, err
	}

	return categoryID, nil
}

func (cs *categoryService) UpdateCategory(ctx context.Context, category models.Category) error {
	err := cs.categoryRepository.UpdateCategory(ctx, category)
	if err != nil {
		return err
	}

	return nil
}

func (cs *categoryService) DeleteCategory(ctx context.Context, categoryID int) error {
	_, err := cs.categoryRepository.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return err
	}

	// what happend when data already soft delete ?

	err = cs.categoryRepository.DeleteCategory(ctx, categoryID)
	if err != nil {
		return err
	}

	return nil
}
