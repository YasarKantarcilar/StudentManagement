package collections

import (
	client "studentmanagement.com/db"
)

var UsersCollection = client.Client.Database("schoolmanagement").Collection("users")
var StudentsCollection = client.Client.Database("schoolmanagement").Collection("students")
var TeachersCollection = client.Client.Database("schoolmanagement").Collection("teachers")
var CoursesCollection = client.Client.Database("schoolmanagement").Collection("courses")
var ClassesCollection = client.Client.Database("schoolmanagement").Collection("classes")
var EnrollmentsCollection = client.Client.Database("schoolmanagement").Collection("enrollments")
var AssignmentsCollection = client.Client.Database("schoolmanagement").Collection("assignments")
var GradesCollection = client.Client.Database("schoolmanagement").Collection("grades")
