package middlewares

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mfa-face-recog/pkg/auth/config"
	"github.com/mfa-face-recog/pkg/auth/utils"
)

type MFASession struct {
	ID          int        `db:"id" json:"id"`
	UserID      int        `db:"user_id" json:"user_id"`
	PosVerified int        `db:"pos_verified" json:"pos_verified"`
	NegVerified int        `db:"neg_verified" json:"neg_verified"`
	Match       bool       `db:"match" json:"match"`
	Used        bool       `db:"used" json:"used"`
	UsedAt      *time.Time `db:"used_at" json:"used_at"`
	CreatedAt   *time.Time `db:"created_at" json:"created_at"`
}
type RegisterSession struct {
	ID        int        `db:"id" json:"id"`
	UserID    int        `db:"user_id" json:"user_id"`
	Used      bool       `db:"used" json:"used"`
	UsedAt    *time.Time `db:"used_at" json:"used_at"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
}

func AuthMiddleware(c *fiber.Ctx) error {
	bearerToken := c.Get("Authorization")
	if bearerToken == "" {
		fmt.Println("No token")
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"success": "false", "message": "Unauthorized"})
	}

	token := strings.Split(bearerToken, "BEARER ")[1]

	valid, err := utils.VerifyAccessToken(token)
	if err != nil {
		fmt.Println("Error verifying token", err)
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"success": "false", "message": "Unauthorized"})
	}
	if !valid {
		fmt.Println("Token not valid")
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"success": "false", "message": "Unauthorized"})
	}

	userId, err := utils.GetClaimFromToken(token, os.Getenv("JWT_SECRET"), "id")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"success": "false", "message": "Unauthorized"})
	}

	c.Locals("user_id", userId)

	return c.Next()
}

func MFAMiddleware(c *fiber.Ctx) error {
	bearerToken := c.Get("Authorization")
	if bearerToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"success": "false", "message": "Unauthorized"})
	}

	token := strings.Split(bearerToken, "BEARER ")[1]

	valid, err := utils.VerifyMFAToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"success": "false", "message": "Unauthorized"})
	}
	if !valid {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"success": "false", "message": "Unauthorized"})
	}

	userId, err := utils.GetClaimFromToken(token, os.Getenv("JWT_MFA_TOKEN_SECRET"), "id")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"success": "false", "message": "Unauthorized"})
	}

	c.Locals("user_id", userId)

	return c.Next()
}

func MFASessionMiddleware(c *fiber.Ctx) error {
	bearerToken := c.Get("Authorization")
	if bearerToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"success": "false", "message": "Unauthorized"})
	}

	token := strings.Split(bearerToken, "BEARER ")[1]

	valid, err := utils.VerifyMFASession(token)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"success": "false", "message": "Unauthorized"})
	}
	if !valid {
		fmt.Println("token not valid")
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"success": "false", "message": "Unauthorized"})
	}

	userId, err := utils.GetClaimFromToken(token, os.Getenv("JWT_MFA_SESSION_SECRET"), "user_id")
	if err != nil {
		fmt.Println("user id not found in session token")
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"success": "false", "message": "Unauthorized"})
	}
	sessionId, err := utils.GetClaimFromToken(token, os.Getenv("JWT_MFA_SESSION_SECRET"), "id")
	if err != nil {
		fmt.Println("session_id not found on token")
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"success": "false", "message": "Unauthorized"})
	}

	session := MFASession{
		ID: sessionId,
	}
	err = config.DB.Get(&session, "SELECT * from mfa_sessions where id=$1", sessionId)
	if err != nil {
		fmt.Printf("error while session fetch %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"success": "false", "message": "Session not found"})
	}

	if session.Used {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"success": "false", "message": "Session already expired"})
	}

	c.Locals("user_id", userId)
	c.Locals("session_id", sessionId)

	return c.Next()
}

func MFARegisterSessionMiddleware(c *fiber.Ctx) error {
	bearerToken := c.Get("Authorization")
	if bearerToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"success": "false", "message": "Unauthorized"})
	}

	token := strings.Split(bearerToken, "BEARER ")[1]

	valid, err := utils.VerifyMFARegisterSession(token)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"success": "false", "message": "Unauthorized"})
	}
	if !valid {
		fmt.Println("token not valid")
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"success": "false", "message": "Unauthorized"})
	}

	userId, err := utils.GetClaimFromToken(token, os.Getenv("JWT_MFA_REGISTER_SESSION_SECRET"), "user_id")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"success": "false", "message": "Unauthorized"})
	}
	sessionId, err := utils.GetClaimFromToken(token, os.Getenv("JWT_MFA_REGISTER_SESSION_SECRET"), "id")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"success": "false", "message": "Unauthorized"})
	}

	session := RegisterSession{
		ID: sessionId,
	}
	fmt.Printf("session id = %d", sessionId)
	err = config.DB.Get(&session, "SELECT * from register_session where id=$1", sessionId)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"success": "false", "message": "Session not found"})
	}

	if session.Used {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"success": "false", "message": "Session already expired"})
	}

	c.Locals("user_id", userId)
	c.Locals("session_id", sessionId)

	return c.Next()
}
