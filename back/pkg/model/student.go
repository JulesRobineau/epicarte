package model

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	// ID is the id of the student
	ID uint64 `json:"id" gorm:"primarykey"`
	// Email is the email of the student
	Email string `json:"email" gorm:"unique"`
	// FirstName is the first name of the student
	FirstName string `json:"first_name" gorm:"not null;size:120"`
	// LastName is the last name of the student
	LastName string `json:"last_name" gorm:"not null;size:120"`
	// ClassID is the foreign key to the class table
	ClassID uint64 `json:"class_id"`
	// Class is the class of the student
	Class *Class `json:"class" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName returns the name of the table
func (s *Student) TableName() string {
	return "students"
}

type StudentModel struct {
	Tx *gorm.DB
}

// NewStudentModel creates a new student model
func NewStudentModel(tx *gorm.DB) *StudentModel {
	return &StudentModel{Tx: tx}
}

// GetStudentsClass returns the class of the student
func (s *StudentModel) GetStudentsClass(classId uint64, student []Student) *gorm.DB {
	return s.Tx.Model(student).Where("class_id = ?", classId).Preload("Class").Find(student)
}

// AddStudentToClass adds a student to a class
func (s *StudentModel) AddStudentToClass(student *Student) error {
	return s.Tx.Association("Class").Append(student)
}

// RemoveStudentFromClass removes a student from a class
func (s *StudentModel) RemoveStudentFromClass(student *Student) error {
	return s.Tx.Association("Class").Delete(student)
}
