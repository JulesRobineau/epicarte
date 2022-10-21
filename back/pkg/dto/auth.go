package dto

import "time"

type Login struct {
	// Username is the username of the user
	Username string `json:"username" binding:"omitempty,min=3,alphanum"`
	// Email is the email of the user
	Email string `json:"email" binding:"omitempty,email"`
	// Password is the password of the user
	Password string `json:"password" binding:"required"`
}

type Register struct {
	// Username is the username of the user
	Username string `json:"username" binding:"required,min=3,max=20,alphanum"`
	// Email is the email of the user
	Email string `json:"email" binding:"required,email"`
	// Password is the password of the user
	Password string `json:"password" binding:"required,min=6"`
	// FirstName is the first name of the user
	FirstName string `json:"first_name" binding:"required,min=2"`
	// LastName is the last name of the user
	LastName string `json:"last_name" binding:"required,min=2"`
}

type ChangePassword struct {
	// UserID is the id of the user
	UserId uint64 `json:"-"`
	// OldPassword is the old password of the user
	OldPassword string `json:"old_password" binding:"required"`
	// NewPassword is the new password of the user
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type AuthResponse struct {
	// AccessToken is the access token of the user
	AccessToken string `json:"access_token"`
	// RefreshToken is the refresh token of the user
	RefreshToken string `json:"refresh_token"`
	// ExpiresIn is the expiration time of the access token
	ExpiresIn time.Time `json:"expires_in"`
}
