package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mfa-face-recog/pkg/auth/config"
)

func UserRoutes(app *fiber.App) {
	app.Get("/api/v1/user", func(c *fiber.Ctx) error {

		id := c.Locals("user_id").(int)
		user := User{
			ID: id,
		}
		err := config.DB.Get(&user, "SELECT * FROM users WHERE id = $1", id)
		if err != nil {
			fmt.Print(err)
			return c.Status(fiber.StatusBadRequest).SendString("User not found")
		}
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{"id": user.ID, "name": user.Name, "email": user.Email})
	})
}
