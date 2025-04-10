package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mfa-face-recog/pkg/auth/config"
	"github.com/mfa-face-recog/pkg/auth/middlewares"
	"github.com/mfa-face-recog/pkg/auth/utils"
)

type SessionStatus struct {
	IsComplete bool `json:"isComplete"`
	IsSuccess  bool `json:"isSuccess"`
	IsFailed   bool `json:"isFailed"`
}

func MFARoutes(app *fiber.App) {
	app.Post("/api/v1/mfa/face/register/image", func(c *fiber.Ctx) error {
		id := c.Locals("user_id").(int)
		sessionId := c.Locals("session_id").(int)
		pubKey := c.FormValue("public_key")

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
		err = RegisterImageToFaceRecognitionService(fileRead, strconv.Itoa(id))
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusBadRequest).SendString("Error sending image to face recognition service")
		}
		config.DB.MustExec(`UPDATE users SET pub = $1 WHERE id = $2`, pubKey, user.ID)
		config.DB.MustExec(`UPDATE register_session SET used = $1 WHERE id = $2`, true, sessionId)
		// config.DB.MustExec(`UPDATE users SET mfa = true WHERE id = $1`, id)
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{"success": "true"})

	})
	app.Get("/api/v1/mfa/register/sessiontoken", func(c *fiber.Ctx) error {
		id := c.Locals("user_id").(int)
		user := User{
			ID: id,
		}
		fmt.Println(id)
		err := config.DB.Get(&user, "SELECT * FROM users WHERE id = $1", id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("User not found")
		}
		session, err := utils.CreateMFARegisterSession(id)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Error creating session token")
		}
		fmt.Printf("created\n sessionID: %s sessionToken: %s\n", session.ID, session.Token)
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{"success": "true", "session_token": session.Token, "session_id": session.ID})
	})
	app.Get("/api/v1/mfa/session/:id/status", func(c *fiber.Ctx) error {
		id := c.Params("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid session ID")
		}
		session := &middlewares.MFASession{
			ID: idInt,
		}
		err = config.DB.Get(session, "SELECT * FROM mfa_sessions WHERE id = $1", idInt)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Session not found")
		}

		if session.Used {
			return c.Status(fiber.StatusBadRequest).JSON(&SessionStatus{IsComplete: true, IsSuccess: session.Match, IsFailed: !session.Match})
		}
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{"isComplete": false, "isSuccess": false, "isFailed": false})
	})
	app.Get("/api/v1/mfa/register/session/:id/status", func(c *fiber.Ctx) error {
		id := c.Params("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid session ID")
		}
		session := &middlewares.RegisterSession{
			ID: idInt,
		}
		err = config.DB.Get(session, "SELECT * FROM register_session WHERE id = $1", idInt)
		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusBadRequest).SendString("Session not found")
		}

		if session.Used {
			return c.Status(fiber.StatusBadRequest).JSON(&SessionStatus{IsComplete: true, IsSuccess: true, IsFailed: false})
		}
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{"isComplete": false, "isSuccess": false, "isFailed": false})
	})
	app.Get("/api/v1/mfa/sessiontoken", func(c *fiber.Ctx) error {
		id := c.Locals("user_id").(int)
		user := User{
			ID: id,
		}
		fmt.Println(id)
		err := config.DB.Get(&user, "SELECT * FROM users WHERE id = $1", id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("User not found")
		}
		session, err := utils.CreateMFASession(id)
		fmt.Printf("created\n sessionID: %s sessionToken: %s\n", session.ID, session.Token)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Error creating session token")
		}
		if user.Pub == nil {
			return c.Status(fiber.StatusBadRequest).SendString("Public key not found")
		}
		encToken, err := utils.Encrypt(session.Token, *user.Pub)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Error encrypting session token")
		}
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{"success": "true", "session_token": encToken, "session_id": session.ID})
	})
	app.Post("api/v1/mfa/face/verify", func(c *fiber.Ctx) error {
		id := c.Locals("user_id").(int)
		sessionId := c.Locals("session_id").(int)

		fmt.Println(id)

		faceImage, err := c.FormFile("face_image")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("User ID and face image are required")
		}
		user := User{
			ID: id,
		}

		session := middlewares.MFASession{
			ID: sessionId,
		}

		err = config.DB.Get(&user, "SELECT * FROM users WHERE id = $1", id)
		if err != nil {
			fmt.Print(err)
			return c.Status(fiber.StatusBadRequest).SendString("User not found")
		}

		err = config.DB.Get(&session, "SELECT * FROM mfa_sessions WHERE id = $1", sessionId)

		if err != nil {
			fmt.Print(err)
			return c.Status(fiber.StatusBadRequest).SendString("User not found")
		}
		fileRead, err := faceImage.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Error opening face image")
		}
		verify, err := VerifyImageOnFaceRecognitionService(fileRead, strconv.Itoa(id))
		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusBadRequest).SendString("Error verifying image")
		}
		fmt.Println(verify)

		if verify.Verified {
			if session.PosVerified < 2 {
				config.DB.MustExec(`UPDATE mfa_sessions SET pos_verified = $1, neg_verified = 0 WHERE id = $2`, session.PosVerified+1, sessionId)
			} else {
				config.DB.MustExec(`UPDATE mfa_sessions SET match = $1, used = $2, used_at = CURRENT_TIMESTAMP WHERE id = $3`, true, true, sessionId)
			}
			return c.Status(fiber.StatusOK).JSON(&fiber.Map{"verified": "true"})
		}
		if session.NegVerified < 4 {
			config.DB.MustExec(`UPDATE mfa_sessions SET neg_verified = $1 WHERE id = $2`, session.NegVerified+1, sessionId)
		} else {
			config.DB.MustExec(`UPDATE mfa_sessions SET match = $1, used = $2, used_at = CURRENT_TIMESTAMP WHERE id = $3`, false, true, sessionId)
		}
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{"verified": "false"})
	})
}

type VerifyResponse struct {
	Verified bool `json:"success"`
}
type VerifyServiceResponse struct {
	Status    string  `json:"status"`
	Verified  bool    `json:"verified"`
	Distance  float64 `json:"distance"`
	Threshold float64 `json:"threshold"`
}

func VerifyImageOnFaceRecognitionService(fileRead io.Reader, name string) (*VerifyResponse, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("image", "face.jpg")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(fw, fileRead)
	if err != nil {
		return nil, err
	}

	nw, err := w.CreateFormField("name")
	if err != nil {
		return nil, err
	}
	_, err = io.WriteString(nw, name)
	if err != nil {
		return nil, err
	}
	w.Close()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/face-recognition", os.Getenv("FACE_RECOGNITION_SERVICE_URL")), &b)
	if err != nil {
		return nil, err
		// return c.Status(fiber.StatusBadRequest).SendString("Error opening face image")
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("face recognition service returned non-200 status code")
		// return c.Status(fiber.StatusBadRequest).SendString("Error opening face image")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data VerifyServiceResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	if data.Verified {
		return &VerifyResponse{Verified: true}, nil
	}
	return &VerifyResponse{Verified: false}, nil
}

func RegisterImageToFaceRecognitionService(fileRead io.Reader, name string) error {
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/register", os.Getenv("FACE_RECOGNITION_SERVICE_URL")), &b)
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
