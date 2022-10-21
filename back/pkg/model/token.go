package model

import "gorm.io/gorm"

type Token struct {
	ID        int64          `json:"id" gorm:"primaryKey"`
	TokenID   string         `json:"token_id" gorm:"uniqueIndex:unique_idx_token;not null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName returns the name of the table
func (t *Token) TableName() string {
	return "tokens"
}

type TokenModel struct {
	Tx *gorm.DB
}

// FindToken finds a token by id
func (t *TokenModel) FindToken(model *Token) *gorm.DB {
	return t.Tx.Where("deleted_at is NULL").First(model)
}

// CreateToken creates a new token in the database
func (t *TokenModel) CreateToken(model *Token) error {
	return t.Tx.Create(model).Error
}

// DeleteToken deletes a token from the database
func (t *TokenModel) DeleteToken(model Token) error {
	return t.Tx.Where("token_id = ?", model.TokenID).Delete(&model).Error
}
