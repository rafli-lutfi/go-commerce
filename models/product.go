package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name       string  `json:"name" gorm:"size:100; not null; unique"`
	Desc       string  `json:"desc" gorm:"size:256; not null"`
	Price      float64 `json:"price" gorm:"not null"`
	Quantity   int     `json:"quantity" gorm:"not null"`
	CategoryID uint    `json:"category_id" gorm:"not null"`
	DiscountID uint    `json:"discount_id" gorm:"default:1"`
}

type ProductInfo struct {
	ID              uint    `json:"id"`
	Name            string  `json:"name"`
	Desc            string  `json:"desc"`
	Price           float64 `json:"price"`
	Quantity        int     `json:"quantity"`
	CategoryName    string  `json:"category_name"`
	DiscountName    string  `json:"discount_name"`
	DiscountPercent float64 `json:"discount_percent"`
	Active          bool    `json:"active"`
}
