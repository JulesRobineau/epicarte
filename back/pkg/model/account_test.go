package model

import (
	"fmt"
	"gin-template/pkg/dto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"testing"
)

// SetupTestDatabase sets up the test database
func SetupTestDatabase() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Paris",
		"localhost",
		"postgres",
		"postgres",
		"postgres",
		5432,
	)
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Error,
			},
		),
	})

	// Truncate all tables
	db.Migrator().DropTable(&User{})
	db.Migrator().DropTable(&Account{})
	db.Migrator().DropTable(&Token{})

	// Migrate the schema
	db.AutoMigrate(
		User{},
		Token{},
		Account{},
	)

	return db
}

func TestAccountModel_Find(t *testing.T) {
	// Setup the test database
	db := SetupTestDatabase()

	// insert a user account by batch
	for i := 0; i < 10; i++ {
		// Create a user
		db.Create(&Account{
			Email:    fmt.Sprintf("user%d@gmail.com", i),
			Username: fmt.Sprintf("username%d", i),
			Password: "password",
			User: User{
				FirstName: fmt.Sprintf("firstname%d", i),
				LastName:  fmt.Sprintf("lastname%d", i),
			},
		})
	}

	// FindByEmailOrUsername the user by email
	accountModel := AccountModel{Tx: db}
	account := Account{Email: "user3@gmail.com"}
	err := accountModel.FindByEmailOrUsername(&account)
	if err.Error != nil {
		t.Error(err)
	}

	// Check if the user is found
	if account.ID == 0 {
		t.Error("User with email is not found")
	}

	// FindByEmailOrUsername the user by username
	account = Account{Username: "username3"}
	err = accountModel.FindByEmailOrUsername(&account)
	if err.Error != nil {
		t.Error(err)
	}

	// Check if the user is found
	if account.ID == 0 {
		t.Error("User is username is not found")
	}
}

func TestAccountModel_FindAll(t *testing.T) {
	// Setup the test database
	db := SetupTestDatabase()

	// insert a user account by batch
	for i := 0; i < 10; i++ {
		// Create a user
		db.Create(&Account{
			Email:    fmt.Sprintf("user%d@gmail.com", i),
			Username: fmt.Sprintf("username%d", i),
			Password: "password",
			User: User{
				FirstName: fmt.Sprintf("firstname%d", i),
				LastName:  fmt.Sprintf("lastname%d", i),
			},
		})
	}

	// FindByEmailOrUsername all users
	accountModel := AccountModel{Tx: db}
	var accounts []Account
	err := accountModel.FindAll(&accounts, dto.UserQueryParams{
		Page:     1,
		PageSize: 10,
		Sort:     "id",
		Order:    "desc",
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

func TestAccountModel_Create(t *testing.T) {
	// Setup the test database
	db := SetupTestDatabase()

	// Try to get the user by email
	accountModel := AccountModel{Tx: db}
	account := Account{Email: "user.create@gmail.com"}
	tx := accountModel.FindByEmailOrUsername(&account)
	if tx.Error == nil {
		t.Error("User with email is found, but it should not")
	}

	// Create a user
	account = Account{
		Email:    "user.create@gmail.com",
		Username: "username.create",
		Password: "password",
		User: User{
			FirstName: "firstname.create",
			LastName:  "lastname.create",
		},
	}
	err := accountModel.Create(&account)
	if err != nil {
		t.Error(err)
	}

	// Try to get the user by email
	account = Account{Email: "user.create@gmail.com"}
	tx = accountModel.FindByEmailOrUsername(&account)
	if tx.Error != nil {
		t.Error("User with email is found, but it should not")
	}

	// Check if the user is found
	if account.Username != "username.create" {
		t.Error("User is not created")
	}
}

func TestAccountModel_Update(t *testing.T) {
	// Setup the test database
	db := SetupTestDatabase()

	// Create a user
	db.Create(&Account{
		Email:    "user.update@gmail.com",
		Username: "username.update",
		Password: "password",
		User: User{
			FirstName: "firstname.update",
			LastName:  "lastname.update",
		},
	})

	// Try to get the user by email
	accountModel := AccountModel{Tx: db}
	account := Account{Email: "user.update@gmail.com"}
	tx := accountModel.FindByEmailOrUsername(&account)
	if tx.Error != nil {
		t.Error("User with email is not found")
	}

	// Update the user
	account.Password = "newpassword"
	err := accountModel.Update(&account)
	if err != nil {
		t.Error(err)
	}

	// Try to get the user by email
	account = Account{Email: "user.update@gmail.com"}
	tx = accountModel.FindByEmailOrUsername(&account)
	if tx.Error != nil {
		t.Error("User with email is not found")
	}

	// Check if the user has been updated
	if account.Password != "newpassword" {
		t.Error("User is not updated")
	}
}

func TestAccountModel_Delete(t *testing.T) {
	// Setup the test database
	db := SetupTestDatabase()

	// Create a user
	tx := db.Create(&Account{
		Email:    "user.delete@gmail.com",
		Username: "username.delete",
		Password: "password",
		User: User{
			FirstName: "firstname.delete",
			LastName:  "lastname.delete",
		},
	})
	if tx.Error != nil {
		t.Error(tx.Error)
	}

	// Try to get the user by email
	accountModel := AccountModel{Tx: db}
	account := Account{Email: "user.delete@gmail.com"}
	tx = accountModel.FindByEmailOrUsername(&account)
	if tx.Error != nil {
		t.Error("User with email is not found, but it should")
	}

	if tx.RowsAffected == 0 {
		t.Error("User with email is not found")
	}

	// Delete the user
	userId := account.User.ID
	err := accountModel.Delete(account)
	if err != nil {
		t.Error(err)
	}

	// Try to get the user by email
	account = Account{Email: "user.delete@gmail.com"}
	tx = accountModel.FindByEmailOrUsername(&account)
	if tx.Error == nil {
		t.Error("User with email is found, but it should not")
	}

	user := User{ID: userId}
	tx = db.Find(&user)
	if tx.Error != nil {
		t.Error("User is found, but it should not")
	}
}
