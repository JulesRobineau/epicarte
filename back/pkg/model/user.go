package model

import (
	"fmt"
	"gin-template/pkg/dto"
	"gin-template/pkg/model/enum"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	// ID is the primary key
	ID uint64 `gorm:"primarykey"`
	// AccountID is the foreign key to the account table
	AccountID uint64 `json:"account_id"`
	// Account is the account of the user
	Account *Account `json:"account" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// FirstName is the first name of the user
	FirstName string `json:"first_name;not null;size:120"`
	// LastName is the last name of the user
	LastName string `json:"last_name;not null;size:120"`
	// Role is the role of the user
	Role enum.Role `json:"role" gorm:"type:role;default:student"`
}

// TableName returns the name of the table
func (u *User) TableName() string {
	return "users_t"
}

type UserModel struct {
	Tx *gorm.DB
}

// NewUserModel creates a new user model
func NewUserModel(tx *gorm.DB) *UserModel {
	return &UserModel{Tx: tx}
}

// Find finds a user by id
func (u *UserModel) Find(model *User) *gorm.DB {
	return u.Tx.First(&model)
}

// FindByUserID finds an account by user id
func (u *UserModel) FindByUserID(userId uint64, model *User) *gorm.DB {
	return u.Tx.Preload("Account").First(model, userId)
}

// FindAll finds all users in the database and returns them paginated
func (u *UserModel) FindAll(models []User, params dto.UserQueryParams) *gorm.DB {
	return u.Tx.Scopes(func(db *gorm.DB) *gorm.DB {
		offset := (params.Page - 1) * params.PageSize
		return db.Offset(offset).Limit(params.PageSize)
	}).Order(
		fmt.Sprintf("%s %s", params.Sort, params.Order),
	).Find(&models)
}

// Update updates a user in the database
// and returns an error if there is one
func (u *UserModel) Update(model *User) error {
	u.Tx.Joins("Account").Updates(model)
	if u.Tx.Error != nil {
		return u.Tx.Error
	}

	return nil
}

// Delete deletes a user from the database
// and returns an error if there is one
func (u *UserModel) Delete(model User) error {
	return u.Tx.Select("User").Delete(&Account{ID: model.AccountID}).Error
}

// AddUserToSession adds a user to a session
func (u *UserModel) AddUserToSession(user *User, session *Session) error {
	return u.Tx.Model(session).Association("Students").Append(user)
}

// RemoveUserFromSession removes a user from a session
func (u *UserModel) RemoveUserFromSession(user *User, session *Session) error {
	return u.Tx.Model(session).Association("Students").Delete(user)
}

// Create creates a new user in the database
func (u *UserModel) Create(model *User) error {
	return u.Tx.Create(model).Error
}
