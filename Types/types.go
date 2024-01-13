package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName      string             `json:"firstName" bson:"firstName"`
	LastName       string             `json:"lastName" bson:"lastName"`
	Email          string             `json:"email" bson:"email"`
	Password       string             `json:"password" bson:"password"`
	StudentID      uint16             `json:"studentID" bson:"studentID"`
	RegisterDate   time.Time          `json:"registerDate" bson:"registerDate"`
	LastUpdateDate time.Time          `json:"lastUpdateDate" bson:"lastUpdateDate"`
	ParentID       primitive.ObjectID `json:"parentID" bson:"parentID"`
	Status         bool               `json:"status" bson:"status"`
}

type Parent struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName      string             `json:"firstName" bson:"firstName"`
	LastName       string             `json:"lastName" bson:"lastName"`
	Email          string             `json:"email" bson:"email"`
	Phone          string             `json:"phone" bson:"phone"`
	Status         bool               `json:"status" bson:"status"`
	RegisterDate   time.Time          `json:"registerDate" bson:"registerDate"`
	LastUpdateDate time.Time          `json:"lastUpdateDate" bson:"lastUpdateDate"`
}

type Teacher struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName      string             `json:"firstName" bson:"firstName"`
	LastName       string             `json:"lastName" bson:"lastName"`
	Email          string             `json:"email" bson:"email"`
	Password       string             `json:"password" bson:"password"`
	Phone          string             `json:"phone" bson:"phone"`
	Status         bool               `json:"status" bson:"status"`
	Branch         string             `json:"branch" bson:"branch"`
	RegisterDate   time.Time          `json:"registerDate" bson:"registerDate"`
	LastUpdateDate time.Time          `json:"lastUpdateDate" bson:"lastUpdateDate"`
}

type Classroom struct {
	ID        string             `json:"id,omitempty" bson:"_id,omitempty"`
	TeacherID string             `json:"teacherID" bson:"teacherID"`
	Year      uint16             `json:"year" bson:"year"`
	GradeID   primitive.ObjectID `json:"gradeID" bson:"gradeID"`
	Section   string             `json:"section" bson:"section"`
	Status    bool               `json:"status" bson:"status"`
	Remarks   string             `json:"remarks" bson:"remarks"`
}

type Grade struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
}

type Course struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
}

type ClassroomStudent struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ClassroomID primitive.ObjectID `json:"classroomID" bson:"classroomID"`
	StudentID   primitive.ObjectID `json:"studentID" bson:"studentID"`
}

type Attandance struct {
	ID        string             `json:"id,omitempty" bson:"_id,omitempty"`
	StudentID primitive.ObjectID `json:"studentID" bson:"studentID"`
	Date      time.Time          `json:"date" bson:"date"`
	Status    bool               `json:"status" bson:"status"`
	Remarks   string             `json:"remarks" bson:"remarks"`
}
