package studentcontrollers

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	T "studentmanagement.com/Types"
	db "studentmanagement.com/db"
)

func SignJWT(ID string, Admin bool) (string, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	secret := os.Getenv("JWT_SECRET_KEY")

	claims := jwt.MapClaims{
		"id":    ID,
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, nil

}

func GetStudents(c *fiber.Ctx) error {
	var students []T.Student
	cursor, err := db.StudentsCollection.Find(c.Context(), bson.M{})
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if err := cursor.All(c.Context(), &students); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(students)
}

func GetStudent(c *fiber.Ctx) error {
	var student T.Student
	id := c.Params("id")
	studentMongoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	findErr := db.StudentsCollection.FindOne(c.Context(), bson.M{"_id": studentMongoID}).Decode(&student)
	if findErr != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(student)
}

func CreateStudent(c *fiber.Ctx) error {
	var payload T.Student

	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 8)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	emailFilter := bson.M{"email": payload.Email}
	count, err := db.StudentsCollection.CountDocuments(c.Context(), emailFilter)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if count > 0 {
		return c.SendStatus(fiber.StatusConflict)
	}

	payload.Password = string(hashedPassword)
	payload.StudentID = uint16(checkStudentID(c))

	if count > 0 {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	insertResult, err := db.StudentsCollection.InsertOne(c.Context(), payload)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusCreated).JSON(insertResult.InsertedID)
}

func UpdateStudent(c *fiber.Ctx) error {
	id := c.Params("id")
	studentMongoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	var payload T.Student
	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	update := bson.M{
		"$set": bson.M{
			"firstName": payload.FirstName,
			"lastName":  payload.LastName,
			"email":     payload.Email,
			"studentID": payload.StudentID,
		},
	}
	_, err = db.StudentsCollection.UpdateOne(c.Context(), bson.M{"_id": studentMongoID}, update)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}

func DeleteStudent(c *fiber.Ctx) error {
	id := c.Params("id")
	studentMongoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	_, err = db.StudentsCollection.DeleteOne(c.Context(), bson.M{"_id": studentMongoID})
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}

func LoginStudent(c *fiber.Ctx) error {
	var payload T.Student
	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	var student T.Student
	findErr := db.StudentsCollection.FindOne(c.Context(), bson.M{"email": payload.Email}).Decode(&student)
	if findErr != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	err := bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(payload.Password))
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	token, err := SignJWT(student.ID.Hex(), false)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Cookie(&fiber.Cookie{
		Name:    "jwt",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 72),
	})

	return c.SendStatus(fiber.StatusOK)
}

func randomStudentID() uint16 {
	randomStudentID := rand.Intn(9999) + 1
	return uint16(randomStudentID)
}

func checkStudentID(c *fiber.Ctx) uint16 {
	randomID := randomStudentID()
	findStudentByID := bson.M{"studentID": randomID}
	count, err := db.StudentsCollection.CountDocuments(c.Context(), findStudentByID)
	if err != nil {
		return checkStudentID(c)
	}
	if count > 0 {
		return checkStudentID(c)
	}
	return randomID
}
