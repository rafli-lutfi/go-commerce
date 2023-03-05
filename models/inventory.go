package models

import "gorm.io/gorm"

type Inventory struct { // one to one with product
	gorm.Model
	Quantity int     `json:"quantity" gorm:"not null"`
	Product  Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
