package routes

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"

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
	MFA      bool   `json:"mfa"`
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
	app.Post("/api/v1/mfa/face/register", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.FormValue("user_id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
		}
		faceImage, err := c.FormFile("face_image")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("User ID and face image are required")
		}
		user := User{
			ID: id,
		}
		fmt.Println(id)
		err = config.DB.Get(&user, "SELECT * FROM users WHERE id = $1", id)
		if err != nil {
			fmt.Print(err)
			return c.Status(fiber.StatusBadRequest).SendString("User not found")
		}
		if user.MFA {
			return c.Status(fiber.StatusBadRequest).SendString("User already has MFA enabled")
		}

		fileRead, err := faceImage.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Error opening face image")
		}
		err = SendImageToFaceRecognitionService(fileRead, strconv.Itoa(id))
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusBadRequest).SendString("Error sending image to face recognition service")
		}

		config.DB.MustExec(`UPDATE users SET mfa = true WHERE id = $1`, id)
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{"success": "true"})

	})
}

func SendImageToFaceRecognitionService(fileRead io.Reader, name string) error {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("image", "face.jpg")
	if err != nil {
		return err
	}
	_, err = io.Copy(fw, fileRead)
	if err != nil {
		return err
	}
	nw, err := w.CreateFormField("name")
	if err != nil {
		return err
	}
	_, err = io.WriteString(nw, name)
	if err != nil {
		return err
	}
	w.Close()
	req, err := http.NewRequest("POST", "http://localhost:5000/api/v1/register", &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("face recognition service returned non-200 status code")
	}
	fmt.Println(resp.StatusCode)
	return nil
}
