package repository

import (
	"context"
	"errors"

	"github.com/rafli-lutfi/go-commerce/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	GetOrderByID(ctx context.Context, orderID int) (models.Order, error)
	CreateNewOrder(ctx context.Context, orderItem models.OrderItem, userID int) (models.Order, error)
	JoinTableOrderItem(ctx context.Context, orderItem *models.OrderItem) (models.OrderItem, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *orderRepository {
	return &orderRepository{db}
}

func (or *orderRepository) GetOrderByID(ctx context.Context, orderID int) (models.Order, error) {
	var order models.Order

	err := or.db.WithContext(ctx).Preload("OrderItems").Preload("Payment").Where("id = ? ", orderID).Find(&order).Error
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (or *orderRepository) CreateNewOrder(ctx context.Context, orderItem models.OrderItem, userID int) (models.Order, error) {
	newOrder := models.Order{}

	txErr := or.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// find product in table product
		if err := tx.Find(&models.Product{}, orderItem.ProductID).Error; err != nil {
			tx.Rollback()
			return err
		}

		// update stock in table product
		if err := tx.Model(&models.Product{}).Where("id = ?", orderItem.ProductID).Update("quantity", gorm.Expr("quantity - ?", orderItem.Quantity)).Error; err != nil {
			tx.Rollback()
			return err
		}

		// find order is exist or not
		if err := tx.Debug().Preload("Payment", "status = ?", "UNPAID").Where("user_id = ?", userID).First(&newOrder).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			// order not exist and make a new order
			newOrder.UserID = uint(userID)
			newOrder.Total = 0
			newOrder.Payment = models.Payment{
				Name:   "CASH",
				Amount: 0,
				Status: "UNPAID",
			}

			if err := tx.Create(&newOrder).Error; err != nil {
				tx.Rollback()
				return err
			}

		} else if err != nil {
			tx.Rollback()
			return err
		}

		// fmt.Println(orderItem.Price, "=", productPrice, "-", discountPrice)

		// insert orderItem to table
		// newOrderItem := models.OrderItem{
		// 	OrderID:   newOrder.ID,
		// 	ProductID: orderItem.ProductID,
		// 	Quantity:  orderItem.Quantity,
		// }
		orderItem.OrderID = newOrder.ID

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return err
		}

		err := tx.Preload("Product.Discount").Find(&orderItem).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		discountPrice := orderItem.Product.Discount.DiscountPercent * orderItem.Product.Price
		productPrice := orderItem.Product.Price

		orderItem.Price = float64(orderItem.Quantity) * (productPrice - discountPrice)

		if err := tx.Model(&models.OrderItem{}).Where("id = ?", orderItem.ID).Update("price", orderItem.Price).Error; err != nil {
			tx.Rollback()
			return err
		}

		err = tx.Model(&newOrder).Association("OrderItems").Append(&orderItem)
		if err != nil {
			return err
		}

		newOrder.Total += orderItem.Price

		if err := tx.Model(&models.Order{}).Where("id = ?", newOrder.ID).Update("total", newOrder.Total).Error; err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})

	if txErr != nil {
		return models.Order{}, txErr
	}

	return newOrder, nil
}

func (or *orderRepository) AddOrderItems(ctx context.Context, orderItems []models.OrderItem, orderID int) (models.Order, error) {
	Order, err := or.GetOrderByID(ctx, orderID)
	if err != nil {
		return models.Order{}, err
	}

	err = or.db.WithContext(ctx).Model(&Order).Association("OrderItems").Append(orderItems)
	if err != nil {
		return models.Order{}, err
	}

	return Order, nil
}

func (or *orderRepository) UpdateOrder(ctx context.Context, orderItems []models.OrderItem) (models.Order, error) {
	for _, item := range orderItems {
		err := or.db.WithContext(ctx).Model(&models.OrderItem{}).Where("order_id = ?", item.OrderID).Updates(item).Error
		if err != nil {
			return models.Order{}, err
		}
	}
	return models.Order{}, nil
}

func (or *orderRepository) JoinTableOrderItem(ctx context.Context, orderItem *models.OrderItem) (models.OrderItem, error) {
	err := or.db.WithContext(ctx).Preload("Product.Discount").Find(&orderItem).Error
	if err != nil {
		return models.OrderItem{}, err
	}

	return *orderItem, nil
}
