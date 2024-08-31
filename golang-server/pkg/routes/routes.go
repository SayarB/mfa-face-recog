package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mfa-face-recog/pkg/config"
	"github.com/mfa-face-recog/pkg/utils"
)

type UserRegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Post("/api/v1/register", func(c *fiber.Ctx) error {
		var req UserRegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}
		if req.Name == "" || req.Email == "" || req.Password == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Email and password are required")
		}
		alreadyExists := &User{
			ID: -1,
		}
		config.DB.Get(alreadyExists, "SELECT * FROM users WHERE email = $1", req.Email)
		fmt.Println("exists", alreadyExists.ID)
		if alreadyExists.ID != -1 {
			return c.Status(fiber.StatusBadRequest).SendString("Email already exists")
		}
		hashedPassword := utils.HashPassword(req.Password)
		config.DB.MustExec(`INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`, req.Name, req.Email, hashedPassword)
		config.DB.Get(alreadyExists, "SELECT * FROM users WHERE email = $1", req.Email)

		return c.Status(fiber.StatusCreated).JSON(&fiber.Map{"id": alreadyExists.ID, "name": alreadyExists.Name, "email": alreadyExists.Email})

	})
	app.Post("/api/v1/login", func(c *fiber.Ctx) error {
		var req UserLoginRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}
		if req.Email == "" || req.Password == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Email and password are required")
		}
		user := User{
			ID: -1,
		}
		config.DB.Get(&user, "SELECT * FROM users WHERE email = $1", req.Email)
		if user.ID == -1 {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid email or password")
		}
		if user.Password != utils.HashPassword(req.Password) {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid email or password")
		}
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{"success": "true"})
	})
}
