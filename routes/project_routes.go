package routes

import (
	"packlify-cloud-backend/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupProjectRoutes(app *fiber.App) {
	api := app.Group("/project")
	api.Post("/", controllers.CreateProject)
}
