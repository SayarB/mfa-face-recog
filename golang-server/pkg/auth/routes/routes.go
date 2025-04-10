package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mfa-face-recog/pkg/auth/middlewares"
)

func RegisterRoutes(app *fiber.App) {
	AuthRoutes(app)
	MFARoutes(app)
	UserRoutes(app)
}

func RegisterMiddlewares(app *fiber.App) {
	app.Use("/api/v1/mfa/face/verify", middlewares.MFASessionMiddleware)
	app.Use("/api/v1/mfa/face/register/image", middlewares.MFARegisterSessionMiddleware)
	app.Use("/api/v1/mfa/register/sessiontoken", middlewares.MFAMiddleware)
	app.Use("/api/v1/mfa/sessiontoken", middlewares.MFAMiddleware)
	app.Use("/api/v1/user", middlewares.MFAMiddleware)
}
