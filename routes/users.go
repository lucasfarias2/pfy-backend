package routes

import (
	"github.com/gofiber/fiber/v2"
)

func createUser(c *fiber.Ctx) error {
	return c.SendString("User created!")
}

func getUser(c *fiber.Ctx) error {
	return c.SendString("User details!")
}

// SetupUserRoutes sets up all the user routes
func SetupUsersRoutes(app *fiber.App) {
	usersGroup := app.Group("/users")
	usersGroup.Get("/create", createUser)
	usersGroup.Get("/get", getUser)
}
