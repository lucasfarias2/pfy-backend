package main

import (
	"log"
	"packlify-cloud-backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	app := fiber.New()

	routes.SetupUsersRoutes(app)
	routes.SetupProjectsRoute(app)
	routes.SetupAuthRoutes(app)

	log.Println("Server started at :8080")
	app.Listen(":8080")
}
