package models

import "time"

type OrderDetail struct {
	ID         uint        `json:"id" gorm:"primaryKey"`
	UserID     uint        `json:"user_id" gorm:"not null"`
	Total      float64     `json:"total"`
	OrderItems []OrderItem `json:"order_items"`
	Payment    Payment     `json:"payment"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	OrderID   uint      `json:"order_id" gorm:"not null"`
	ProductID uint      `json:"product_id" gorm:"not null"`
	Price     float64   `json:"price" gorm:"not null"`
	Quantity  int       `json:"quantity" gorm:"default: 1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
