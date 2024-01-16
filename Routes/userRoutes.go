package routers

import (
	"github.com/gofiber/fiber/v2"
	MW "studentmanagement.com/Middlewares"
	controllers "studentmanagement.com/controllers"
)

func UserRoutes(app fiber.Router) {
	app.Get("/me", MW.AuthMiddleware, controllers.GetCurrentUser)
	app.Get("/user/:id", MW.TeacherAuthMiddleware, controllers.GetUser)
	app.Get("/", MW.AdminMiddleware, controllers.GetAllUsers)
	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)
}
