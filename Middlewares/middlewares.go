package middlewares

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func AuthMiddleware(c *fiber.Ctx) error {
	_, err := JwtParse(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON("Unauthorized")
	}
	return c.Next()
}

func AdminMiddleware(c *fiber.Ctx) error {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	claims, err := JwtParse(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON("Unauthorized")
	}
	if claims["admin"] != true {
		return c.Status(fiber.StatusUnauthorized).JSON("Not an admin")
	}
	return c.Next()

}

func JwtParse(c *fiber.Ctx) (jwt.MapClaims, error) {

	tokenString := c.Cookies("jwt")
	secret := os.Getenv("JWT_SECRET_KEY")
	token, parseErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if parseErr != nil {
		return nil, parseErr
	}
	claims := token.Claims.(jwt.MapClaims)
	return claims, nil
}
