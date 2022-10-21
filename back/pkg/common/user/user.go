package user

import (
	"gin-template/pkg/dto"
	"gin-template/pkg/model"
	error2 "gin-template/utils/error"
	"gorm.io/gorm"
)

// GetUser gets a user from the database by id
func GetUser(db *gorm.DB, userId uint64) (*dto.User, error) {
	userModel := model.UserModel{Tx: db}
	// get user from database
	u := model.User{ID: userId}
	if err := userModel.Find(&u).Error; err != nil {
		return nil, error2.FromDatabaseError(err)
	}
	return &dto.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
	}, nil
}

// GetUserList gets a list of users from the database
func GetUserList(db *gorm.DB, req dto.UserQueryParams) ([]dto.User, error) {
	accountModel := model.AccountModel{Tx: db}
	// get users from database
	var users []model.Account
	if err := accountModel.FindAll(&users, req).Error; err != nil {
		return nil, error2.FromDatabaseError(err)
	}

	// convert to dto
	userList := make([]dto.User, 0)
	for _, u := range users {
		userList = append(userList, dto.User{
			ID:        u.User.ID,
			Email:     u.Email,
			Username:  u.Username,
			FirstName: u.User.FirstName,
			LastName:  u.User.LastName,
			Role:      u.User.Role,
		})
	}
	return userList, nil
}

// SuperAdminUpdateUser updates a user in the database by a super admin
func SuperAdminUpdateUser(db *gorm.DB, req dto.UpdateUser) error {
	userModel := model.UserModel{Tx: db}
	// update user in database
	u := model.User{}
	if err := userModel.FindByUserID(req.Id, &u).Error; err != nil {
		return error2.FromDatabaseError(err)
	}

	// update account in database
	accountModel := model.AccountModel{Tx: db}
	a := model.Account{ID: u.AccountID, Email: req.Email, Username: req.Username}
	if err := accountModel.Update(&a); err != nil {
		return error2.FromDatabaseError(err)
	}

	// update user in database
	u = model.User{ID: req.Id, FirstName: req.FirstName, LastName: req.LastName, Role: req.Role}
	if err := userModel.Update(&u); err != nil {
		return error2.FromDatabaseError(err)
	}

	return nil
}

// UpdateUser updates a user in the database
func UpdateUser(db *gorm.DB, req dto.UpdateUser) error {
	userModel := model.UserModel{Tx: db}
	// update user in database
	u := model.User{}
	if err := userModel.FindByUserID(req.Id, &u).Error; err != nil {
		return error2.FromDatabaseError(err)
	}

	// update account in database
	accountModel := model.AccountModel{Tx: db}
	a := model.Account{ID: u.AccountID, Email: req.Email, Username: req.Username}
	if err := accountModel.Update(&a); err != nil {
		return error2.FromDatabaseError(err)
	}

	// update user in database
	u = model.User{ID: req.Id, FirstName: req.FirstName, LastName: req.LastName}
	if err := userModel.Update(&u); err != nil {
		return error2.FromDatabaseError(err)
	}
	return nil
}
