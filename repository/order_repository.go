package repository

import (
	"context"
	"errors"

	"github.com/rafli-lutfi/go-commerce/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	GetOrderByID(ctx context.Context, orderID int) (models.Order, error)
	GetOrderItemByID(ctx context.Context, orderItemID int) (models.OrderItem, error)
	CreateNewOrder(ctx context.Context, orderItem models.OrderItem, userID int) (models.Order, error)
	UpdateOrder(ctx context.Context, updateItem models.UpdateItem, order models.Order, orderItemBefore models.OrderItem) error
	ConfirmPayment(ctx context.Context, payment models.Payment) error
	OrdersHistory(ctx context.Context, userID int) ([]models.Order, error)
	CurrentOrder(ctx context.Context, userID int) (models.Order, error)
	JoinTableOrderItem(ctx context.Context, orderItem *models.OrderItem) (models.OrderItem, error)
	CheckPaymentStatus(ctx context.Context, paymentID int) (models.Order, error)
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

func (or *orderRepository) CheckPaymentStatus(ctx context.Context, paymentID int) (models.Order, error) {
	var order models.Order

	err := or.db.WithContext(ctx).Preload("Payment").Where("payment_id = ?", paymentID).Find(&order).Error
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (or *orderRepository) CreateNewOrder(ctx context.Context, orderItem models.OrderItem, userID int) (models.Order, error) {
	order := models.Order{}

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
		if check := tx.Debug().Joins("Payment", or.db.Where(&models.Payment{Status: "UNPAID"})).Where("user_id = ?", userID).Find(&order); check.RowsAffected == 0 {
			newOrder := models.Order{}

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

			order = newOrder

		} else if check.Error != nil {
			tx.Rollback()
			return check.Error
		}

		if err := tx.Preload("OrderItems").Find(&order).Error; err != nil {
			tx.Rollback()
			return err
		}

		// update item if already exist
		for _, item := range order.OrderItems {
			if item.ProductID == orderItem.ProductID {
				tx.Rollback()
				return errors.New("item already in order")
			}
		}

		orderItem.OrderID = order.ID

		// Create New order item
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

		// update order item with new price
		if err := tx.Model(&models.OrderItem{}).Where("id = ?", orderItem.ID).Update("price", orderItem.Price).Error; err != nil {
			tx.Rollback()
			return err
		}

		// append order item into order
		err = tx.Model(&order).Association("OrderItems").Append(&orderItem)
		if err != nil {
			return err
		}

		order.Total += orderItem.Price

		// update order total payment
		if err := tx.Model(&models.Order{}).Where("id = ?", order.ID).Update("total", order.Total).Error; err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})

	if txErr != nil {
		return models.Order{}, txErr
	}

	return order, nil
}

func (or *orderRepository) GetOrderItemByID(ctx context.Context, orderItemID int) (models.OrderItem, error) {
	var orderItem models.OrderItem

	err := or.db.Find(&orderItem, orderItemID).Error
	if err != nil {
		return models.OrderItem{}, err
	}

	return orderItem, nil
}

func (or *orderRepository) UpdateOrder(ctx context.Context, updateItem models.UpdateItem, order models.Order, orderItemBefore models.OrderItem) error {
	quantityBefore := orderItemBefore.Quantity
	productID := orderItemBefore.ProductID

	return or.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		product := models.Product{}
		if err := tx.Preload("Discount").Find(&product, productID).Error; err != nil {
			tx.Rollback()
			return err
		}
		discountPrice := product.Price * product.Discount.DiscountPercent

		var priceUpdate float64

		if updateItem.Quantity == 0 {
			// update quantity product
			err := tx.Model(models.Product{}).Where("id = ?", productID).Update("quantity", gorm.Expr("quantity + ?", quantityBefore)).Error
			if err != nil {
				tx.Rollback()
				return err
			}

			// deleted order item
			if err := tx.Where("id = ?", updateItem.ID).Delete(&models.OrderItem{}).Error; err != nil {
				tx.Rollback()
				return err
			}

			order.Total -= float64(quantityBefore) * (product.Price - discountPrice)

		} else {
			priceBefore := float64(quantityBefore) * (product.Price - discountPrice)
			priceUpdate = float64(updateItem.Quantity) * (product.Price - discountPrice)

			if updateItem.Quantity > quantityBefore {
				err := tx.Model(models.Product{}).Where("id = ?", productID).Update("quantity", gorm.Expr("quantity - ?", (updateItem.Quantity-quantityBefore))).Error
				if err != nil {
					tx.Rollback()
					return err
				}
				order.Total += priceUpdate - priceBefore
			}

			if updateItem.Quantity < quantityBefore {
				err := tx.Model(models.Product{}).Where("id = ?", productID).Update("quantity", gorm.Expr("quantity + ?", (quantityBefore-updateItem.Quantity))).Error
				if err != nil {
					tx.Rollback()
					return err
				}
				order.Total -= priceBefore - priceUpdate
			}

		}

		if err := tx.Model(&models.OrderItem{}).Where("id = ?", updateItem.ID).Updates(models.OrderItem{Price: priceUpdate, Quantity: updateItem.Quantity}).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Model(&models.Order{}).Where("id = ?", order.ID).Update("total", order.Total).Error; err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})
}

func (or *orderRepository) ConfirmPayment(ctx context.Context, payment models.Payment) error {
	return or.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		err := or.db.Model(&models.Payment{}).Where("id  = ?", payment.ID).Updates(
			&models.Payment{
				Name:   payment.Name,
				Amount: payment.Amount,
				Status: "PAID",
			}).Error

		if err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})
}

func (or *orderRepository) OrdersHistory(ctx context.Context, userID int) ([]models.Order, error) {
	var orders []models.Order

	rows, err := or.db.WithContext(ctx).Model(&models.Order{}).Joins("Payment", or.db.Where("status = ?", "PAID")).Where("user_id = ?", userID).Rows()
	if err != nil {
		return []models.Order{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var order models.Order

		err := or.db.ScanRows(rows, &order)
		if err != nil {
			return []models.Order{}, err
		}

		err = or.db.WithContext(ctx).Preload("OrderItems").Find(&order).Error
		if err != nil {
			return []models.Order{}, err
		}

		orders = append(orders, order)
	}
	return orders, nil
}

func (or *orderRepository) CurrentOrder(ctx context.Context, userID int) (models.Order, error) {
	var order models.Order

	err := or.db.Joins("Payment", or.db.Where("status = ?", "UNPAID")).Where("user_id = ?", userID).Find(&order).Error
	if err != nil {
		return models.Order{}, err
	}

	err = or.db.Preload("OrderItems").Find(&order).Error
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (or *orderRepository) JoinTableOrderItem(ctx context.Context, orderItem *models.OrderItem) (models.OrderItem, error) {
	err := or.db.WithContext(ctx).Preload("Product.Discount").Find(&orderItem).Error
	if err != nil {
		return models.OrderItem{}, err
	}

	return *orderItem, nil
}
