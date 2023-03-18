package services

import (
	"context"
	"errors"

	"github.com/rafli-lutfi/go-commerce/models"
	"github.com/rafli-lutfi/go-commerce/repository"
)

type OrderService interface {
	GetOrderByID(ctx context.Context, orderID int) (models.Order, error)
	CreateNewOrder(ctx context.Context, orderItem models.OrderItem, userID int) (models.Order, error)
	UpdateOrder(ctx context.Context, updateItem models.UpdateItem) error
}
type orderService struct {
	orderRepository repository.OrderRepository
}

func NewOrderService(orderRepository repository.OrderRepository) *orderService {
	return &orderService{orderRepository}
}

func (os *orderService) GetOrderByID(ctx context.Context, orderID int) (models.Order, error) {
	order, err := os.orderRepository.GetOrderByID(ctx, orderID)
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}
func (os *orderService) CreateNewOrder(ctx context.Context, orderItem models.OrderItem, userID int) (models.Order, error) {
	Order, err := os.orderRepository.CreateNewOrder(ctx, orderItem, userID)
	if err != nil {
		return models.Order{}, err
	}

	return Order, nil
}

func (os *orderService) UpdateOrder(ctx context.Context, updateItem models.UpdateItem) error {
	// find order first
	order, err := os.orderRepository.GetOrderByID(ctx, int(updateItem.OrderID))
	if err != nil {
		return err
	}

	// get order item
	orderItem, err := os.orderRepository.GetOrderItemByID(ctx, int(updateItem.ID))
	if err != nil {
		return err
	}

	// check orderitem is exist or not in order table
	for _, item := range order.OrderItems {
		if item.ID == updateItem.ID {
			break
		}

		if order.OrderItems[len(order.OrderItems)-1].ID != updateItem.ID {
			return errors.New("item not found in order")
		}
	}

	// if update quantity is same with before
	if updateItem.Quantity == orderItem.Quantity {
		return errors.New("nothing changed")
	}

	err = os.orderRepository.UpdateOrder(ctx, updateItem, order, orderItem.Quantity)
	if err != nil {
		return err
	}

	return nil
}

// func (os *orderService) ConfirmOrder(ctx context.Context) {}
