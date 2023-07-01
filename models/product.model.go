package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID   string `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Name string `gorm:"type:varchar(100);not null;unique" json:"name"`
}

type Basket struct {
	ID        string       `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	UserID    string       `gorm:"type:char(36);not null" json:"userId,omitempty"`
	CreatedAt time.Time    `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time    `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
	Items     []BasketItem `gorm:"foreignKey:BasketID" json:"items,omitempty"`
}

type BasketItem struct {
	ID        string    `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	BasketID  string    `gorm:"type:char(36);not null" json:"basketId,omitempty"`
	ProductID string    `gorm:"type:char(36);not null" json:"productId,omitempty"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

type Product struct {
	ID         string    `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Title      string    `gorm:"type:varchar(255);uniqueIndex:idx_products_title,LENGTH(255);not null" json:"title,omitempty"`
	Content    string    `gorm:"not null" json:"content,omitempty"`
	Category   Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	CategoryID string    `gorm:"type:char(36)" json:"categoryId,omitempty"`
	Published  bool      `gorm:"default:false;not null" json:"published"`
	CreatedAt  time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt  time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

type CreateCategorySchema struct {
	Name string `json:"name" validate:"required"`
}

type UpdateCategorySchema struct {
	Name string `json:"name,omitempty"`
}

func (product *Product) BeforeCreate(tx *gorm.DB) (err error) {
	product.ID = uuid.New().String()
	return nil
}

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func ValidateStruct[T any](payload T) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

type CreateProductSchema struct {
	Title      string `json:"title" validate:"required"`
	Content    string `json:"content" validate:"required"`
	CategoryID string `json:"categoryId,omitempty"`
	Published  bool   `json:"published,omitempty"`
}

type UpdateProductSchema struct {
	Title      string `json:"title,omitempty"`
	Content    string `json:"content,omitempty"`
	CategoryID string `json:"categoryId,omitempty"`
	Published  *bool  `json:"published,omitempty"`
}

type AddProductToBasketSchema struct {
	BasketID  string `json:"basketId" validate:"required"`
	ProductID string `json:"productId" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,gte=1"`
}
