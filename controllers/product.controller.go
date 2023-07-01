package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wpcodevo/golang-fiber-mysql/initializers"
	"github.com/wpcodevo/golang-fiber-mysql/models"
	"gorm.io/gorm"
)

func CreateProductHandler(c *fiber.Ctx) error {
	var payload *models.CreateProductSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := models.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	now := time.Now()
	newProduct := models.Product{
		Title:      payload.Title,
		Content:    payload.Content,
		CategoryID: payload.CategoryID, // Here's the change
		Published:  payload.Published,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	result := initializers.DB.Create(&newProduct)

	if result.Error != nil && strings.Contains(result.Error.Error(), "Duplicate entry") {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "Title already exist, please use another title"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"product": newProduct}})
}

func FindProducts(c *fiber.Ctx) error {
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var products []models.Product
	results := initializers.DB.Limit(intLimit).Offset(offset).Find(&products)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(products), "products": products})
}

func UpdateProduct(c *fiber.Ctx) error {
	productId := c.Params("productId")

	var payload *models.UpdateProductSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var product models.Product
	result := initializers.DB.First(&product, "id = ?", productId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No product with that Id exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	updates := make(map[string]interface{})
	if payload.Title != "" {
		updates["title"] = payload.Title
	}
	if payload.CategoryID != "" { // Here's the change
		updates["category_id"] = payload.CategoryID // Here's the change
	}
	if payload.Content != "" {
		updates["content"] = payload.Content
	}

	if payload.Published != nil {
		updates["published"] = payload.Published
	}

	updates["updated_at"] = time.Now()

	initializers.DB.Model(&product).Updates(updates)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"product": product}})
}

func FindProductById(c *fiber.Ctx) error {
	productId := c.Params("productId")

	var product models.Product
	result := initializers.DB.First(&product, "id = ?", productId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No product with that Id exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"product": product}})
}

func DeleteProduct(c *fiber.Ctx) error {
	productId := c.Params("productId")

	result := initializers.DB.Delete(&models.Product{}, "id = ?", productId)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No product with that Id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.SendStatus(fiber.StatusNoContent)
}



func AddProductToBasket(c *fiber.Ctx) error {
	var payload *models.AddProductToBasketSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	newBasketItem := models.BasketItem{
		BasketID:  payload.BasketID,
		ProductID: payload.ProductID,
		Quantity:  payload.Quantity,
	}

	result := initializers.DB.Create(&newBasketItem)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"basketItem": newBasketItem}})
}


