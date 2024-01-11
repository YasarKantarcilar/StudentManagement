package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	StudentRoutes "studentmanagement.com/Routes/StudentRoutes"
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

	StudentGroup := app.Group("/students")

	StudentRoutes.StudentRoutes(StudentGroup)

	// start server on port 3000
	log.Fatal(app.Listen(":3000"))

}
