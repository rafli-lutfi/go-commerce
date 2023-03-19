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
	ConfirmPayment(ctx context.Context, payment models.Payment) (models.Order, error)
	OrderHistory(ctx context.Context, userID int) ([]models.Order, error)
	CurrentOrder(ctx context.Context, userID int) (models.Order, error)
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

	if order.ID == 0 {
		return order, errors.New("record not found")
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

	if order.Payment.Status == "PAID" {
		return errors.New("order has already paid")
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

	// get order item
	orderItem, err := os.orderRepository.GetOrderItemByID(ctx, int(updateItem.ID))
	if err != nil {
		return err
	}

	// if update quantity is same with before
	if updateItem.Quantity == orderItem.Quantity {
		return errors.New("nothing changed")
	}

	err = os.orderRepository.UpdateOrder(ctx, updateItem, order, orderItem)
	if err != nil {
		return err
	}

	return nil
}

func (os *orderService) ConfirmPayment(ctx context.Context, payment models.Payment) (models.Order, error) {
	orderDB, err := os.orderRepository.CheckPaymentStatus(ctx, int(payment.ID))
	if err != nil {
		return models.Order{}, err
	}

	if orderDB.Payment.Status == "PAID" {
		return models.Order{}, errors.New("order already paid")
	}

	if float64(payment.Amount) < orderDB.Total {
		return models.Order{}, errors.New("payment amount is not enough")
	}

	err = os.orderRepository.ConfirmPayment(ctx, payment)
	if err != nil {
		return models.Order{}, err
	}

	order, err := os.orderRepository.GetOrderByID(ctx, int(orderDB.ID))
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (os *orderService) OrderHistory(ctx context.Context, userID int) ([]models.Order, error) {
	orders, err := os.orderRepository.OrdersHistory(ctx, userID)
	if err != nil {
		return []models.Order{}, err
	}

	return orders, nil
}

func (os *orderService) CurrentOrder(ctx context.Context, userID int) (models.Order, error) {
	order, err := os.orderRepository.CurrentOrder(ctx, userID)
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}
