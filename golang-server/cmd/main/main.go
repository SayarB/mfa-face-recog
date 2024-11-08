package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/mfa-face-recog/pkg/auth/config"
	"github.com/mfa-face-recog/pkg/auth/routes"
)

func main() {
	godotenv.Load(".env")
	config.ConnectDB()
	defer config.DB.Close()
	app := fiber.New(fiber.Config{
		Immutable: true,
	})

	app.Use(cors.New())
	routes.RegisterMiddlewares(app)
	routes.RegisterRoutes(app)
	app.Listen(":8000")
}
