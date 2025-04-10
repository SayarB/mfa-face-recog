package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mfa-face-recog/pkg/auth/config"
	"github.com/mfa-face-recog/pkg/auth/utils"
)

type UserRegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type User struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	MFA      bool    `json:"mfa"`
	Pub      *string `json:"pub"`
}

type UserLoginRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	PublicKey string `json:"public_key"`
}

func AuthRoutes(app *fiber.App) {
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
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"success": "false", "message": "Email already exists"})
		}
		hashedPassword := utils.HashPassword(req.Password)
		config.DB.MustExec(`INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`, req.Name, req.Email, hashedPassword)
		config.DB.Get(alreadyExists, "SELECT * FROM users WHERE email = $1", req.Email)

		return c.Status(fiber.StatusCreated).JSON(&fiber.Map{"success": "true", "id": alreadyExists.ID, "name": alreadyExists.Name, "email": alreadyExists.Email})
	})
	app.Post("/api/v1/login", func(c *fiber.Ctx) error {
		var req UserLoginRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}
		user := User{}
		err := config.DB.Get(&user, "SELECT * FROM users WHERE email = $1", req.Email)
		if err != nil {
			fmt.Printf("user not found - %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"success": "false", "message": "Invalid email or password"})
		}
		if user.Password != utils.HashPassword(req.Password) {
			fmt.Println("password mismatch")
			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"success": "false", "message": "Invalid email or password"})
		}
		// save the public key in the database
		token, err := utils.CreateMFAToken(user.Email, user.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"success": "false", "message": "Error creating token"})
		}

		return c.Status(fiber.StatusOK).JSON(&fiber.Map{"success": "true", "id": user.ID, "mfa_enabled": user.MFA, "message": "Login successful", "token": token})
	})
}
