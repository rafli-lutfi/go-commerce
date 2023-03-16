package models

import "time"

type Order struct {
	ID         uint        `json:"id" gorm:"primaryKey"`
	UserID     uint        `json:"user_id" gorm:"not null"`
	Total      float64     `json:"total"`
	OrderItems []OrderItem `json:"order_items"`
	PaymentID  uint        `json:"payment_id"`
	Payment    Payment     `json:"payment"`
	CreatedAt  time.Time   `json:"-"`
	UpdatedAt  time.Time   `json:"-"`
}

type OrderItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	OrderID   uint      `json:"order_id" gorm:"not null"`
	ProductID uint      `json:"product_id" gorm:"not null"`
	Product   Product   `json:"-"`
	Price     float64   `json:"price" gorm:"not null"`
	Quantity  int       `json:"quantity" gorm:"default: 1"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type OrderItemInfo struct {
	ID           uint
	OrderID      uint
	ProductName  string
	ProductPrice float64
}
