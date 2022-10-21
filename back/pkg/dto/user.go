package dto

import "gin-template/pkg/model/enum"

type User struct {
	// ID is the id of the user
	ID uint64 `json:"id"`
	// Username is the username of the user
	Username string `json:"username,omitempty"`
	// Email is the email of the user
	Email string `json:"email,omitempty"`
	// Password is the password of the user
	FirstName string `json:"first_name"`
	// LastName is the last name of the user
	LastName string `json:"last_name"`
	// Role is the role of the user
	Role enum.Role `json:"role"`
}

type CreateUser struct {
	// Username is the username of the user
	Username string `json:"username" binding:"required,min=3,max=20,alphanum"`
	// Email is the email of the user
	Email string `json:"email" binding:"required,email"`
	// FirstName is the first name of the user
	FirstName string `json:"first_name" binding:"required,min=2"`
	// LastName is the last name of the user
	LastName string `json:"last_name" binding:"required,min=2"`
	// Role is the role of the user
	Role string `json:"role" binding:"required,oneof=admin student"`
}

type UpdateUser struct {
	Id uint64 `json:"-" uri:"user_id" path:"user_id"`
	// Username is the username of the user
	Username string `json:"username" binding:"omitempty,min=3,max=20,alphanum"`
	// Email is the email of the user
	Email string `json:"email" binding:"omitempty,email"`
	// FirstName is the first name of the user
	FirstName string `json:"first_name" binding:"omitempty,min=2"`
	// LastName is the last name of the user
	LastName string `json:"last_name" binding:"omitempty,min=2"`
	// Role is the role of the user
	Role enum.Role `json:"role" binding:"omitempty,oneof=superadmin admin user"`
}

type UserQueryParams struct {
	// Page is the page number
	Page int `json:"page" form:"page,default=1" binding:"omitempty,min=1"`
	// PageSize is the page size
	PageSize int `json:"page_size" form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
	// Sort sort by field:
	// * id
	// * username
	// * email
	// * first_name
	// * last_name
	// * created_at
	// * updated_at
	Sort string `json:"sort" form:"sort,default=id" binding:"omitempty,oneof=id username email first_name last_name created_at updated_at"`
	// Order by field:
	// * asc: ascending order
	// * desc: descending order
	Order string `json:"order" form:"order,default=asc" binding:"omitempty,oneof=asc desc"`
}
