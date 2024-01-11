package studentroutes

import (
	"github.com/gofiber/fiber/v2"
	StudentControllers "studentmanagement.com/Controllers/StudentControllers"
	MW "studentmanagement.com/Middlewares"
)

func StudentRoutes(app fiber.Router) {
	app.Post("/login", StudentControllers.LoginStudent)
	app.Use(MW.AuthMiddleware)
	app.Get("/:id", StudentControllers.GetStudent)
	app.Use(MW.AdminMiddleware)
	app.Get("/", StudentControllers.GetStudents)
	app.Post("/", StudentControllers.CreateStudent)
	app.Put("/:id", StudentControllers.UpdateStudent)
	app.Delete("/:id", StudentControllers.DeleteStudent)
}
