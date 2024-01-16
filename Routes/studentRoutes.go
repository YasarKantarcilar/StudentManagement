package routers

import (
	"github.com/gofiber/fiber/v2"
	MW "studentmanagement.com/Middlewares"
)

func StudentRoutes(app fiber.Router) {
	app.Use(MW.AuthMiddleware)
	app.Use(MW.AdminMiddleware)
}
