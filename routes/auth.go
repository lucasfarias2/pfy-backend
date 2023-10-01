package routes

import (
	"github.com/gofiber/fiber/v2"
)

func register(c *fiber.Ctx) error {
	return c.SendString("User created!")
}

func login(c *fiber.Ctx) error {
	return c.SendString("User logged in!")
}

func logout(c *fiber.Ctx) error {
	return c.SendString("User logged out!")
}

func SetupAuthRoutes(app *fiber.App) {
	authGroup := app.Group("/auth")
	authGroup.Get("/register", register)
	authGroup.Get("/login", login)
	authGroup.Get("/logout", logout)
}
