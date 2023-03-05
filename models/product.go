package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string  `json:"name" gorm:"size:100; not null; unique"`
	Desc        string  `json:"desc" gorm:"size:256; not null"`
	Price       float64 `json:"price" gorm:"not null"`
	CategoryID  uint    `json:"category_id" gorm:"not null"`
	DiscountID  uint    `json:"discount_id" gorm:"default:0"`
	InventoryID uint    `json:"inventory_id" gorm:"not null"`
}

type NewProduct struct {
	Name       string  `json:"name"`
	Desc       string  `json:"desc"`
	Price      float64 `json:"price"`
	Quantity   int     `json:"quantity"`
	CategoryID uint    `json:"category_id"`
	DiscountID uint    `json:"discount_id"`
}
