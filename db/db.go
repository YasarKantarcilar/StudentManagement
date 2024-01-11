package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB connects to the database
var Client, _ = ConnectDB()

func ConnectDB() (*mongo.Client, error) {

	// Set your MongoDB connection string and database name
	connectionString := "mongodb://127.0.0.1:27017"

	// Set a timeout for the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new MongoDB client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	// Ping the MongoDB server to check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB")
	return client, nil
}

var StudentsCollection = Client.Database("schoolmanagement").Collection("students")
var TeachersCollection = Client.Database("schoolmanagement").Collection("teachers")
var CoursesCollection = Client.Database("schoolmanagement").Collection("courses")
var ClassesCollection = Client.Database("schoolmanagement").Collection("classes")
var EnrollmentsCollection = Client.Database("schoolmanagement").Collection("enrollments")
var AssignmentsCollection = Client.Database("schoolmanagement").Collection("assignments")
var GradesCollection = Client.Database("schoolmanagement").Collection("grades")
