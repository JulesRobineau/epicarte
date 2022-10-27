package auth

import (
	"gin-template/config"
	"gin-template/pkg/dto"
	"gin-template/pkg/model"
	"gin-template/pkg/model/enum"
	error2 "gin-template/utils/error"
	"gin-template/utils/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Login login a user from dto.Login struct and return a token
func Login(db *gorm.DB, req dto.Login, jwtConfig config.JwtConfig) (*dto.AuthResponse, error) {
	userModel := model.AccountModel{Tx: db}
	// find a by email or username
	a := model.Account{Username: req.Username, Email: req.Email}
	if err := userModel.FindByEmailOrUsername(&a).Error; err != nil {
		return nil, error2.FromDatabaseError(err)
	}

	// check password is correct
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(req.Password))
	if err != nil {
		return nil, error2.BadRequestError("password is incorrect", nil)
	}

	// generate token
	token, err := jwt.GenerateTokens(a.ID, a.User.Role, jwtConfig)
	if err != nil {
		return nil, error2.InternalServerError("", err)
	}

	// save token id to database
	tokenModel := model.TokenModel{Tx: db}
	if err = tokenModel.CreateToken(&model.Token{TokenID: token.TokenID}); err != nil {
		return nil, error2.FromDatabaseError(err)
	}

	// return token
	return &dto.AuthResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.ExpiresIn,
	}, nil
}

// Logout logout a user from token id and return an error if there is one
func Logout(db *gorm.DB, tokenId string) error {
	// delete token from database
	tokenModel := model.TokenModel{Tx: db}
	return tokenModel.DeleteToken(model.Token{TokenID: tokenId})
}

// Register register a user from dto.Register and return a token
func Register(db *gorm.DB, req dto.Register, jwtConfig config.JwtConfig) (*dto.AuthResponse, error) {
	// Generate password hash
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, error2.InternalServerError("", nil)
	}

	// Create Account
	accountModel := model.AccountModel{Tx: db}
	a := model.Account{
		Email:    req.Email,
		Username: req.Username,
		Password: string(hash),
		User: model.User{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Role:      enum.STUDENT,
		},
	}
	if err = accountModel.Create(&a); err != nil {
		return nil, error2.FromDatabaseError(err)
	}

	// Generate token
	token, err := jwt.GenerateTokens(a.ID, a.User.Role, jwtConfig)
	if err != nil {
		return nil, error2.InternalServerError("", err)
	}

	// Save token id to database
	tokenModel := model.TokenModel{Tx: db}
	if err = tokenModel.CreateToken(&model.Token{TokenID: token.TokenID}); err != nil {
		return nil, error2.FromDatabaseError(err)
	}

	// Return token
	return &dto.AuthResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.ExpiresIn,
	}, nil
}

// ChangePassword change a user password from request dto.ChangePassword and return a token
func ChangePassword(db *gorm.DB, req dto.ChangePassword) error {
	// find u by id
	userModel := model.AccountModel{Tx: db}
	u := model.Account{User: model.User{ID: req.UserId}}

	if err := userModel.FindByEmailOrUsername(&u).Error; err != nil {
		return error2.FromDatabaseError(err)
	}

	// check old password is correct
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.OldPassword))
	if err != nil {
		return error2.BadRequestError("password is incorrect", nil)
	}

	// generate new password hash
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return error2.InternalServerError("", nil)
	}

	// update password
	u.Password = string(hash)
	if err = userModel.Update(&u); err != nil {
		return error2.FromDatabaseError(err)
	}

	return nil
}

// RefreshToken refresh a user token from refresh token and return a token
func RefreshToken(db *gorm.DB, jwtConfig config.JwtConfig, token string) (*dto.AuthResponse, error) {
	// parse token
	claims, err := jwt.ParseToken(token, jwtConfig.Secret)
	if err != nil {
		return nil, error2.InternalServerError("", err)
	}

	// generate new token
	newToken, err := jwt.GenerateTokens(claims.UserId, claims.Role, jwtConfig)
	if err != nil {
		return nil, error2.InternalServerError("", err)
	}

	// save token id to database
	tokenModel := model.TokenModel{Tx: db}
	if err = tokenModel.CreateToken(&model.Token{TokenID: newToken.TokenID}); err != nil {
		return nil, error2.FromDatabaseError(err)
	}

	return &dto.AuthResponse{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		ExpiresIn:    newToken.ExpiresIn,
	}, nil
}

// GetAccount get a user from user id and return a dto.Account
func GetAccount(db *gorm.DB, userId uint64) (*dto.User, error) {
	accountModel := model.AccountModel{Tx: db}
	a := model.Account{User: model.User{ID: userId}}
	if err := accountModel.Find(&a).Error; err != nil {
		return nil, error2.FromDatabaseError(err)
	}

	return &dto.User{
		ID:        a.User.ID,
		FirstName: a.User.FirstName,
		LastName:  a.User.LastName,
		Email:     a.Email,
		Username:  a.Username,
		Role:      a.User.Role,
	}, nil
}
