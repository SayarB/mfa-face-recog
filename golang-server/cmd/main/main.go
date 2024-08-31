package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mfa-face-recog/pkg/config"
	"github.com/mfa-face-recog/pkg/routes"
)

func main() {
	config.ConnectDB()
	defer config.DB.Close()
	app := fiber.New(fiber.Config{
		Immutable: true,
	})
	routes.RegisterRoutes(app)
	app.Listen(":3000")

}
