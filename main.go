package main

import (
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Struktura modelu Produktu
type Produkt struct {
	gorm.Model
	Nazwa  string
	Cena   float64
	Skład  string
	Sklepy []Sklep `gorm:"many2many:produkty_sklepy;"`
}

// Struktura modelu Sklepu
type Sklep struct {
	gorm.Model
	Nazwa    string
	Produkty []Produkt `gorm:"many2many:produkty_sklepy;"`
}

func main() {
	// Inicjalizacja bazy danych SQLite
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
	
	// Automatyczne migracje tabel
	db.AutoMigrate(&Produkt{}, &Sklep{})

	// Inicjalizacja frameworku Echo
	e := echo.New()

	// Utworzenie kontrolera dla Produktów
	produktController := &ProduktController{
		db: db,
	}

	// Zdefiniowanie endpointów dla Produktów
	e.GET("/produkty", produktController.GetProdukty)
	e.POST("/produkty", produktController.CreateProdukt)
	e.GET("/produkty/:id", produktController.GetProdukt)
	e.PUT("/produkty/:id", produktController.UpdateProdukt)
	e.DELETE("/produkty/:id", produktController.DeleteProdukt)

	// Uruchomienie serwera na porcie 8080
	e.Start(":8080")
}
