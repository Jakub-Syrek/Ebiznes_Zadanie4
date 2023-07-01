package models

import "github.com/jinzhu/gorm"

type Cart struct {
	gorm.Model
	Products []Product `gorm:"many2many:cart_products;"`
}