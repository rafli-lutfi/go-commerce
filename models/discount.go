package models

import "gorm.io/gorm"

type Discount struct { // one to many with product
	gorm.Model
	Name            string  `json:"name" gorm:"size:100; not null;unique"`
	Desc            string  `json:"desc" gorm:"size:256;"`
	DiscountPercent float64 `json:"discount_percent" gorm:"not null"`
	Active          bool    `json:"active" gorm:"default:false"`
	Products        []Product
}

type DiscountInfo struct {
	Name            string  `json:"name"`
	Desc            string  `json:"desc"`
	DiscountPercent float64 `json:"discount_percent"`
	Active          bool    `json:"active"`
}
