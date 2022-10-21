package model

import (
	"fmt"
	"gin-template/pkg/dto"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	ID       uint64 `gorm:"primarykey"`
	Email    string `json:"email" gorm:"uniqueIndex:unique_idx_email;not null"`
	Username string `json:"username" gorm:"uniqueIndex:unique_idx_username;not null;size:80"`
	Password string `json:"password" gorm:"not null"`
	User     User   `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;re"`
}

// TableName returns the name of the table
func (a *Account) TableName() string {
	return "accounts"
}

type AccountModel struct {
	Tx *gorm.DB
}

// Find finds an account in the database
func (a *AccountModel) Find(model *Account) *gorm.DB {
	return a.Tx.Preload("User").First(model)
}

// FindByEmailOrUsername finds a user by email or username
func (a *AccountModel) FindByEmailOrUsername(model *Account) *gorm.DB {
	return a.Tx.Where(
		"email = ?", model.Email,
	).Or("username = ?", model.Username).Preload("User").First(model)
}

// FindAll finds all accounts in the database and returns them
func (a *AccountModel) FindAll(models *[]Account, params dto.UserQueryParams) *gorm.DB {
	return a.Tx.Scopes(func(db *gorm.DB) *gorm.DB {
		offset := (params.Page - 1) * params.PageSize
		return db.Offset(offset).Limit(params.PageSize)
	}).Order(
		fmt.Sprintf("%s %s", params.Sort, params.Order),
	).Preload("User").Find(&models)
}

// Create creates a new account in the database
func (a *AccountModel) Create(model *Account) error {
	return a.Tx.Create(model).Error
}

// Update updates an account in the database
func (a *AccountModel) Update(model *Account) error {
	return a.Tx.Select("User").Updates(model).Error
}

// Delete deletes a account from the database
func (a *AccountModel) Delete(model Account) error {
	return a.Tx.Select("User").Delete(&model).Error
}
