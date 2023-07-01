package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/wpcodevo/golang-fiber-mysql/controllers"
	"github.com/wpcodevo/golang-fiber-mysql/initializers"
)

func init() {
	initializers.ConnectDB()
}

func main() {
	app := fiber.New()
	micro := fiber.New()

	app.Mount("/api", micro)
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	}))

	micro.Route("/products", func(router fiber.Router) {
		router.Post("/", controllers.CreateProductHandler)
		router.Get("", controllers.FindProducts)
	})
	micro.Route("/products/:productId", func(router fiber.Router) {
		router.Delete("", controllers.DeleteProduct)
		router.Get("", controllers.FindProductById)
		router.Patch("", controllers.UpdateProduct)
	})
	micro.Get("/healthchecker", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Welcome to Golang, Fiber, SQLite, and GORM",
		})
	})

	log.Fatal(app.Listen(":8000"))
}
