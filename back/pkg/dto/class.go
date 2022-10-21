package dto

type Class struct {
	// ID is the id of the class
	ID uint `json:"id"`
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
	// Name is the name of the class
	Name string `json:"name" binding:"omitempty,min=3,max=20,alphanum"`
	// Year is the year of the class
	Year string `json:"year" binding:"omitempty"`
}
