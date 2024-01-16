package middlewares

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	T "studentmanagement.com/Types"
	collections "studentmanagement.com/collections"
)

func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Cookies("jwt")

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	secret := os.Getenv("JWT_SECRET_KEY")
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.Next()
}

func TeacherAuthMiddleware(c *fiber.Ctx) error {

	claims, err := JwtParse(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON("Unauthorized")
	}

	if claims["role"] != "teacher" {
		return c.Next()
	}
	if claims["role"] == "admin" {
		return c.Next()
	}
	return c.Status(fiber.StatusUnauthorized).JSON("Unauthorized")
}

func AdminMiddleware(c *fiber.Ctx) error {
	var user T.User
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	claims, err := JwtParse(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON("Unauthorized")
	}
	if claims["role"] != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON("not admin")
	}
	userObjectID, err := primitive.ObjectIDFromHex(claims["id"].(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid IDs")
	}
	userDoc := collections.UsersCollection.FindOne(c.Context(), bson.M{"_id": userObjectID}).Decode(&user)
	if userDoc != nil {
		return c.Status(fiber.StatusNotFound).JSON("User not found")
	}
	if user.Role != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON("not admin")
	}
	log.Println(claims)
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

func SignJWT(ID string, Role string) (string, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	secret := os.Getenv("JWT_SECRET_KEY")
	claims := jwt.MapClaims{
		"id":   ID,
		"role": Role,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, nil

}
