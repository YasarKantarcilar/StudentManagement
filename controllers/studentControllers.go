package controllers

import (
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	T "studentmanagement.com/Types"
	collections "studentmanagement.com/collections"
)

func GetStudent(c *fiber.Ctx) error {
	id := c.Params("id")
	studentObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid ID")
	}

	singleResult, err := collections.StudentsCollection.Find(c.Context(), bson.M{"_id": studentObjectID})
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON("Student not found")
	}

	var student T.Student

	decodeErr := singleResult.Decode(&student)
	if decodeErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Error decoding student")
	}

	return c.JSON(student)
}

func GetAllStudents(c *fiber.Ctx) error {
	var students []T.Student
	cursor, err := collections.StudentsCollection.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Error fetching students")
	}

	cursorErr := cursor.All(c.Context(), &students)
	if cursorErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Error decoding students")
	}
	return c.JSON(students)
}

func CreateStudent(c *fiber.Ctx) error {
	userID := c.Params("id")
	var student T.Student

	if err := c.BodyParser(student); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Error parsing student")
	}

	student.StudentID = checkStudentID(c)
	student.LastUpdateDate = time.Now()
	student.RegisterDate = time.Now()
	student.Status = true
	student.ParentID = primitive.NewObjectID()

	InsertOne, err := collections.StudentsCollection.InsertOne(c.Context(), student)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Error creating student")
	}
	insertID := InsertOne.InsertedID

	studentObjectID, err := primitive.ObjectIDFromHex(insertID.(primitive.ObjectID).Hex())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid IDs")
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid IDs")
	}

	updateResult, err := collections.UsersCollection.UpdateOne(c.Context(), bson.M{"_id": userObjectID}, bson.M{"$set": bson.M{"roleID": studentObjectID}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Error updating user")
	}

	if updateResult.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON("User not found")
	}

	return c.JSON(student)
}

func randomStudentID() uint16 {
	randomStudentID := rand.Intn(9999) + 1
	return uint16(randomStudentID)
}

func checkStudentID(c *fiber.Ctx) uint16 {
	randomID := randomStudentID()
	findStudentByID := bson.M{"studentID": randomID}
	count, err := collections.StudentsCollection.CountDocuments(c.Context(), findStudentByID)
	if err != nil {
		checkStudentID(c)
	}
	if count > 0 {
		checkStudentID(c)
	}
	return randomID
}
