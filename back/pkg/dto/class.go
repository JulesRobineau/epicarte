package dto

type TinyClass struct {
	// ID is the id of the class
	ID uint64 `json:"id"`
	// Name is the name of the class
	Name string `json:"name"`
	// Year is the year of the class
	Year string `json:"year"`
}

type CreateClass struct {
	// Name is the name of the class
	Name string `json:"name" binding:"required,min=3,max=20,alphanum"`
	// Year is the year of the class
	Year string `json:"year" binding:"required"`
}

type UpdateClass struct {
	// ID is the id of the class
	Id uint64 `json:"-" uri:"class_id" path:"class_id"`
	// Name is the name of the class
	Name string `json:"name" binding:"omitempty,min=3,max=20,alphanum"`
	// Year is the year of the class
	Year string `json:"year" binding:"omitempty"`
}

type Class struct {
	TinyClass
	// Students is the list of students in the class
	Students []Student `json:"students"`
}

type ClassList struct {
	// Classes is the list of classes
	Classes []TinyClass `json:"classes"`
}

type Student struct {
	// ID is the id of the student
	ID uint64 `json:"id,omitempty"`
	// Email is the email of the student
	Email string `json:"email" binding:"required,email"`
	// FirstName is the first name of the student
	FirstName string `json:"first_name" binding:"required,min=2"`
	// LastName is the last name of the student
	LastName string `json:"last_name" binding:"required,min=2"`
}

type AddStudentToClass struct {
	// ClassId is the id of the class
	ClassId uint64 `json:"-" uri:"class_id" path:"class_id"`
	// Student is the student to add to the class
	Student Student `json:"student" binding:"required"`
}
