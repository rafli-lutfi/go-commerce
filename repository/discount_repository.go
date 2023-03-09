package repository

import (
	"context"
	"errors"

	"github.com/rafli-lutfi/go-commerce/models"
	"gorm.io/gorm"
)

type DiscountRepository interface {
	GetDiscountByID(ctx context.Context, discountID int) (models.Discount, error)
	CreateDiscount(ctx context.Context, discount models.Discount) (int, error)
	UpdateDiscount(ctx context.Context, discount *models.Discount) error
	DeleteDiscount(ctx context.Context, discountID int) error
}

type discountRepository struct {
	db *gorm.DB
}

func NewDiscountRepository(db *gorm.DB) *discountRepository {
	return &discountRepository{db}
}

func (dr *discountRepository) GetDiscountByID(ctx context.Context, discountID int) (models.Discount, error) {
	var discount = models.Discount{}

	err := dr.db.WithContext(ctx).Where("id = ?", discountID).First(&discount).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Discount{}, nil
	} else if err != nil {
		return models.Discount{}, err
	}

	return discount, nil
}

func (dr *discountRepository) CreateDiscount(ctx context.Context, discount models.Discount) (int, error) {
	err := dr.db.WithContext(ctx).Create(&discount).Error
	if err != nil {
		return 0, err
	}

	return int(discount.ID), nil
}

func (dr *discountRepository) UpdateDiscount(ctx context.Context, discount *models.Discount) error {
	err := dr.db.WithContext(ctx).Model(&models.Discount{}).Where("id = ?", discount.ID).Updates(discount).Error
	if err != nil {
		return err
	}

	return nil
}

func (dr *discountRepository) DeleteDiscount(ctx context.Context, discountID int) error {
	err := dr.db.WithContext(ctx).Where("id = ?", discountID).Delete(&models.Discount{}).Error
	if err != nil {
		return err
	}

	return nil
}
