package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Kontroler dla Produktów
type ProduktController struct {
	db *gorm.DB
}

// Pobierz wszystkie Produkty
func (c *ProduktController) GetProdukty(ctx echo.Context) error {
	var produkty []Produkt
	result := c.db.Preload("Sklepy").Find(&produkty)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, result.Error)
	}
	return ctx.JSON(http.StatusOK, produkty)
}

// Pobierz pojedynczy Produkt
func (c *ProduktController) GetProdukt(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid ID")
	}

	var produkt Produkt
	result := c.db.Preload("Sklepy").First(&produkt, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return ctx.JSON(http.StatusNotFound, "Product not found")
		}
		return ctx.JSON(http.StatusInternalServerError, result.Error)
	}
	return ctx.JSON(http.StatusOK, produkt)
}

// Utwórz nowy Produkt
func (c *ProduktController) CreateProdukt(ctx echo.Context) error {
	produkt := new(Produkt)
	if err := ctx.Bind(produkt); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid data")
	}

	result := c.db.Create(&produkt)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, result.Error)
	}
	return ctx.JSON(http.StatusCreated, produkt)
}

// Zaktualizuj Produkt
func (c *ProduktController) UpdateProdukt(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid ID")
	}

	var produkt Produkt
	result := c.db.First(&produkt, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return ctx.JSON(http.StatusNotFound, "Product not found")
		}
		return ctx.JSON(http.StatusInternalServerError, result.Error)
	}

	if err := ctx.Bind(&produkt); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid data")
	}

	result = c.db.Save(&produkt)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, result.Error)
	}
	return ctx.JSON(http.StatusOK, produkt)
}

// Usuń Produkt
func (c *ProduktController) DeleteProdukt(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid ID")
	}

	var produkt Produkt
	result := c.db.First(&produkt, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return ctx.JSON(http.StatusNotFound, "Product not found")
		}
		return ctx.JSON(http.StatusInternalServerError, result.Error)
	}

	result = c.db.Delete(&produkt)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, result.Error)
	}
	return ctx.NoContent(http.StatusNoContent)
}
