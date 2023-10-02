package routes

import (
	"packlify-cloud-backend/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupProjectRoutes(app *fiber.App) {
	projectGroup := app.Group("/project")
	projectGroup.Post("/", controllers.CreateProject)
	projectGroup.Get("/", controllers.GetAllProjects)
}
