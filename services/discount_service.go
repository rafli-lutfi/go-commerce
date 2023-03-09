package services

import (
	"context"
	"errors"

	"github.com/rafli-lutfi/go-commerce/models"
	"github.com/rafli-lutfi/go-commerce/repository"
)

type DiscountService interface {
	GetDiscountByID(ctx context.Context, discountID int) (models.DiscountInfo, error)
	CreateDiscount(ctx context.Context, disocunt models.Discount) (int, error)
	UpdateDiscount(ctx context.Context, discount *models.Discount) error
	DeleteDiscount(ctx context.Context, discountID int) error
}

type discountService struct {
	discountRepository repository.DiscountRepository
}

func NewDiscountService(discountRepository repository.DiscountRepository) *discountService {
	return &discountService{discountRepository}
}

func (ds *discountService) GetDiscountByID(ctx context.Context, discountID int) (models.DiscountInfo, error) {
	discountDB, err := ds.discountRepository.GetDiscountByID(ctx, discountID)
	if err != nil {
		return models.DiscountInfo{}, err
	}

	if discountDB.ID == 0 {
		return models.DiscountInfo{}, errors.New("discount not found")
	}

	discount := models.DiscountInfo{
		Name:            discountDB.Name,
		Desc:            discountDB.Desc,
		DiscountPercent: discountDB.DiscountPercent,
		Active:          discountDB.Active,
	}

	return discount, nil
}

func (ds *discountService) CreateDiscount(ctx context.Context, disocunt models.Discount) (int, error) {
	discountID, err := ds.discountRepository.CreateDiscount(ctx, disocunt)
	if err != nil {
		return 0, err
	}

	return discountID, nil
}

func (ds *discountService) UpdateDiscount(ctx context.Context, discount *models.Discount) error {
	// checking discount ID
	discountDB, err := ds.discountRepository.GetDiscountByID(ctx, int(discount.ID))
	if err != nil {
		return err
	}

	if discountDB.ID == 0 {
		return errors.New("discount id not found")
	}

	err = ds.discountRepository.UpdateDiscount(ctx, discount)
	if err != nil {
		return err
	}

	return nil
}
func (ds *discountService) DeleteDiscount(ctx context.Context, discountID int) error {
	discountDB, err := ds.discountRepository.GetDiscountByID(ctx, discountID)
	if err != nil {
		return err
	}

	if discountDB.ID == 0 {
		return errors.New("discount id not found")
	}

	err = ds.discountRepository.DeleteDiscount(ctx, discountID)
	if err != nil {
		return err
	}

	return nil
}
