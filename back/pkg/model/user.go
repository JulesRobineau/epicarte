package model

import (
	"fmt"
	"gin-template/pkg/dto"
	"gin-template/pkg/model/enum"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint64    `gorm:"primarykey"`
	AccountID uint64    `json:"account_id"`
	Account   *Account  `json:"account" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FirstName string    `json:"first_name;not null;size:120"`
	LastName  string    `json:"last_name;not null;size:120"`
	Role      enum.Role `json:"role" gorm:"type:role;default:user"`
}

// TableName returns the name of the table
func (u *User) TableName() string {
	return "users_t"
}

type UserModel struct {
	Tx *gorm.DB
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
