package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	routers "studentmanagement.com/Routes"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// fiber application
	app := fiber.New()

	// middleware
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, ngrok-skip-browser-warning",
		AllowCredentials: true,
	}))
	app.Use(logger.New())

	// routes
	usersRoutes := app.Group("/users")

	routers.UserRoutes(usersRoutes)

	// start server on port 3000
	log.Fatal(app.Listen(":3000"))

}
