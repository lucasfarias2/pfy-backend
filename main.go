package main

import (
	"log"
	"packlify-cloud-backend/routes"
	"packlify-cloud-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	utils.ConnectDatabase()

	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Test from backend")
	})

	routes.SetupProjectRoutes(app)

	log.Println("Server started at :8080")
	app.Listen(":8080")
}
