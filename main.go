package main

import (
	"github.com/Jakub-Syrek/Ebiznes_Zadanie4/models"
	"github.com/labstack/echo/v4"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

var DB *gorm.DB

func main() {
	DB = initDB()

	e := echo.New()

	// Endpointy
	e.GET("/products", GetProducts)
	e.POST("/products", CreateProduct)
	e.GET("/products/:id", GetProduct)
	e.PUT("/products/:id", UpdateProduct)
	e.DELETE("/products/:id", DeleteProduct)

	// Uruchamiamy serwer
	e.Start(":8000")
}

func initDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Product{}, &Category{})

	return db
}