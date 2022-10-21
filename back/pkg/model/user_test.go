package model

import (
	"fmt"
	"gin-template/pkg/dto"
	"testing"
)

// TestUserModel_Find finds a user by id
func TestUserModel_Find(t *testing.T) {
	// Setup the test database
	db := SetupTestDatabase()

	db.Create(&Account{
		Email:    "user.find@gmail.com",
		Username: "username.find",
		Password: "password",
		User: User{
			FirstName: "firstname.find",
			LastName:  "lastname.find",
			Role:      "user",
		},
	})

	// Try to get the user by Lastname
	userModel := UserModel{Tx: db}
	user := User{LastName: "lastname.find"}
	tx := userModel.Find(&user)
	if tx.Error != nil {
		t.Error("User with lastname is not found, but it should")
	}

	// Check if the user is found
	if user.FirstName != "firstname.find" {
		t.Error("User is not found")
	}
}

// TestUserModel_FindAll finds all users in the database and returns them
func TestUserModel_FindAll(t *testing.T) {
	// Setup the test database
	db := SetupTestDatabase()

	for i := 0; i < 10; i++ {
		db.Create(&Account{
			Email:    fmt.Sprintf("user%d.findall@gmail.com", i),
			Username: fmt.Sprintf("username%d.findall", i),
			Password: "password",
			User: User{
				FirstName: fmt.Sprintf("firstname%d.findall", i),
				LastName:  fmt.Sprintf("lastname%d.findall", i),
				Role:      "user",
			},
		})

	}

	// FindByEmailOrUsername all users
	userModel := UserModel{Tx: db}
	var users []User
	err := userModel.FindAll(users, dto.UserQueryParams{
		Page:     1,
		PageSize: 10,
		Sort:     "id",
		Order:    "asc",
	})
	if err.Error != nil {
		t.Error(err)
	}

	// Check if the user is found
	if err.RowsAffected == 0 {
		t.Error("Users are not found")
	}

	// Check if the user is found
	if err.RowsAffected != 10 {
		t.Error("Users number is not correct")
	}
}

func TestUserModel_Delete(t *testing.T) {
	// Setup the test database
	db := SetupTestDatabase()

	// Create a user
	db.Create(&Account{
		Email:    "user.delete@gmail.com",
		Username: "username.delete",
		Password: "password",
		User: User{
			FirstName: "firstname.delete",
			LastName:  "lastname.delete",
		},
	})

	// Try to get the account by email
	accountModel := AccountModel{Tx: db}
	account := Account{Email: "user.delete@gmail.com"}
	tx := accountModel.FindByEmailOrUsername(&account)
	if tx.Error != nil {
		t.Error("User with email is not found, but it should")
	}

	// Delete the user
	userModel := UserModel{Tx: db}
	err := userModel.Delete(account.User)
	if err != nil {
		t.Error(err)
	}

	// Check if the user is deleted
	tx = userModel.Find(&account.User)
	if tx.Error == nil {
		t.Error("User is not deleted")
	}
}
