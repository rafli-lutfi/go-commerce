package services

import (
	"context"

	"github.com/rafli-lutfi/go-commerce/models"
	"github.com/rafli-lutfi/go-commerce/repository"
)

type OrderService interface {
	GetOrderByID(ctx context.Context, orderID int) (models.Order, error)
	CreateNewOrder(ctx context.Context, orderItem models.OrderItem, userID int) (models.Order, error)
	GetOrderItemByID(ctx context.Context, orderItem models.OrderItem) (models.OrderItem, error)
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

func (os *orderService) GetOrderItemByID(ctx context.Context, orderItem models.OrderItem) (models.OrderItem, error) {
	_, err := os.orderRepository.JoinTableOrderItem(ctx, &orderItem)
	if err != nil {
		return models.OrderItem{}, err
	}

	return orderItem, nil
}

// func (os *orderService) AddItemsIntoExistingOrder(ctx context.Context, orderItems []models.OrderItem) {
// }
// func (os *orderService) UpdateOrder(ctx context.Context)  {}
// func (os *orderService) ConfirmOrder(ctx context.Context) {}
