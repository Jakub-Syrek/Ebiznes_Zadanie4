package main

import (
	"github.com/Jakub-Syrek/Ebiznes_Zadanie4/models"
	"github.com/labstack/echo/v4"
	"github.com/jinzhu/gorm"
	"net/http"
)

var DB *gorm.DB

func GetProducts(c echo.Context) error {
	var products []models.Product
	DB.Preload("Category").Find(&products)

	return c.JSON(http.StatusOK, products)
}

func CreateProduct(c echo.Context) error {
	product := new(Product)
	if err := c.Bind(product); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := DB.Create(product).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, product)
}

func GetProduct(c echo.Context) error {
	id := c.Param("id")

	var product Product
	if err := DB.Preload("Category").First(&product, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, product)
}

func UpdateProduct(c echo.Context) error {
	id := c.Param("id")

	product := new(Product)
	if err := c.Bind(product); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := DB.Model(product).Where("id = ?", id).Updates(product).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, product)
}

func DeleteProduct(c echo.Context) error {
	id := c.Param("id")

	if err := DB.Where("id = ?", id).Delete(&Product{}).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
