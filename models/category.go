package models

import "gorm.io/gorm"

type Category struct { // one to many with product
	gorm.Model
	Name     string `json:"name" gorm:"size:100; not null;unique"`
	Desc     string `json:"desc" gorm:"size:256;"`
	Products []Product
}
