package models

import "gorm.io/gorm"

type Brand struct {
	gorm.Model
	Name     string    `gorm:"unique;not null;size:100"`
	Products []Product `gorm:"foreignKey:BrandID"`
}

type Category struct {
	gorm.Model
	Name          string        `gorm:"unique;not null;size:100"`
	Subcategories []Subcategory `gorm:"foreignKey:CategoryID"`
}

type Subcategory struct {
	gorm.Model
	Name       string    `gorm:"not null;size:100"`
	CategoryID uint      `gorm:"not null"`
	Products   []Product `gorm:"foreignKey:SubcategoryID"`
}

type Product struct {
	gorm.Model
	Name          string  `gorm:"not null;size:255"`
	BrandID       uint    `gorm:"not null"`
	SubcategoryID uint    `gorm:"not null"`
	Price         float64 `gorm:"not null"`
	ImagePath     string  `gorm:"not null;size:255"`
	Description   string  `gorm:"type:text"`
	IsOnSale      bool    `gorm:"default:false"`
	SalePrice     float64
}
