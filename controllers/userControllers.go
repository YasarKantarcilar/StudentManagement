package controllers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	MW "studentmanagement.com/Middlewares"
	T "studentmanagement.com/Types"
	"studentmanagement.com/collections"
)

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid ID")
	}

	singleResult, err := collections.UsersCollection.Find(c.Context(), bson.M{"_id": userObjectID})
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON("User not found")
	}

	var user T.User

	decodeErr := singleResult.Decode(&user)
	if decodeErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Error decoding user")
	}

	return c.JSON(user)

}

func GetCurrentUser(c *fiber.Ctx) error {
	var user T.User
	claims, err := MW.JwtParse(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON("Unauthorized")
	}
	log.Println(claims)

	userObjectID, err := primitive.ObjectIDFromHex(claims["id"].(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid IDs")
	}
	singleResult := collections.UsersCollection.FindOne(c.Context(), bson.M{"_id": userObjectID})
	if singleResult.Err() != nil {
		return c.Status(fiber.StatusNotFound).JSON("User not found")
	}

	decodeErr := singleResult.Decode(&user)
	if decodeErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Error decoding user")
	}

	return c.JSON(user)
}

func GetAllUsers(c *fiber.Ctx) error {
	users, err := collections.UsersCollection.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Error decoding user")
	}

	defer users.Close(c.Context())

	var allUsers []T.User

	cursorErr := users.All(c.Context(), &allUsers)
	if cursorErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Error decoding user")
	}

	return c.Status(fiber.StatusOK).JSON(allUsers)
}

func UpdateUser(c *fiber.Ctx) error {
	var payload T.User

	parseErr := c.BodyParser(&payload)
	if parseErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Error parsing request")
	}

	id := c.Params("id")
	userObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid ID")
	}

	updateResult, err := collections.UsersCollection.UpdateOne(c.Context(), bson.M{"_id": userObjectID}, bson.M{"$set": bson.M{
		"role": payload.Role,
	}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Error updating user")
	}

	if updateResult.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON("User not found")
	}

	return c.Status(fiber.StatusOK).JSON("User updated")
}

func Register(c *fiber.Ctx) error {
	var payload T.User

	err := c.BodyParser(&payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Payload missing data")
	}

	userCount, err := collections.UsersCollection.CountDocuments(c.Context(), bson.M{"email": payload.Email})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Error checking for existing user")
	}

	if userCount > 0 {
		return c.Status(fiber.StatusBadRequest).JSON("User already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Error hashing password")
	}
	payload.Password = string(hashedPassword)
	payload.Status = true
	payload.RegisterDate = time.Now()
	payload.LastUpdateDate = time.Now()
	payload.Role = "student"

	jwt, err := MW.SignJWT(payload.ID.Hex(), payload.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Error signing JWT")
	}

	c.Cookie(&fiber.Cookie{
		Name:    "jwt",
		Value:   jwt,
		Expires: time.Now().Add(time.Hour * 24),
	})

	_, err = collections.UsersCollection.InsertOne(c.Context(), payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Error creating user")
	}

	return c.Status(fiber.StatusCreated).JSON("User created")
}

func Login(c *fiber.Ctx) error {
	var payload T.User

	var existingUser T.User

	err := c.BodyParser(&payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Payload missing data")
	}

	singleResult := collections.UsersCollection.FindOne(c.Context(), bson.M{"email": payload.Email}).Decode(&existingUser)
	if singleResult == mongo.ErrNoDocuments {
		return c.Status(fiber.StatusBadRequest).JSON("User Doesn't exists")
	}

	passwordErr := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(payload.Password))
	if passwordErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Incorrect Password")
	}

	jwt, err := MW.SignJWT(existingUser.ID.Hex(), existingUser.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Error signing JWT")
	}

	c.Cookie(&fiber.Cookie{
		Name:    "jwt",
		Value:   jwt,
		Expires: time.Now().Add(time.Hour * 24),
	})

	return c.Status(fiber.StatusOK).JSON("Logged in")
}
