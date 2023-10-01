package routes

import (
	"github.com/gofiber/fiber/v2"
)

func createProject(c *fiber.Ctx) error {
	return c.SendString("Project created!")
}

func editProject(c *fiber.Ctx) error {
	return c.SendString("Project updated!")
}

func deleteProject(c *fiber.Ctx) error {
	return c.SendString("Project deleted!")
}

func getUserProjects(c *fiber.Ctx) error {
	return c.SendString("This are the users projects: 1, 2, 3")
}

func SetupProjectsRoute(app *fiber.App) {
	projectsGroup := app.Group("/projects")
	projectsGroup.Get("/create", createProject)
	projectsGroup.Get("/delete", deleteProject)
	projectsGroup.Get("/edit", editProject)
	projectsGroup.Get("/get", getUserProjects)
}
