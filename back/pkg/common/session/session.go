package session

import (
	"fmt"
	"gin-template/pkg/dto"
	"gin-template/pkg/model"
	error2 "gin-template/utils/error"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// GetSessionsByClassID gets all sessions of a class
func GetSessionsByClassID(tx *gorm.DB, classID uint64) (*dto.SessionList, error) {
	sessionModel := model.NewSessionModel(tx)
	sessions, err := sessionModel.FindAll(classID)
	if err != nil {
		return nil, error2.FromDatabaseError(err)
	}

	sessionDtos := make([]dto.TinySession, 0)
	for _, session := range sessions {
		sessionDtos = append(sessionDtos, dto.TinySession{
			ID:        session.ID,
			IsClosed:  session.IsClosed,
			CreatedAt: session.CreatedAt,
			UpdatedAt: session.UpdatedAt,
		})
	}

	return &dto.SessionList{Sessions: sessionDtos}, nil
}

// GetSessionByID gets a session by ID
func GetSessionByID(tx *gorm.DB, sessionID uuid.UUID) (*dto.Session, error) {
	sessionModel := model.NewSessionModel(tx)

	session := model.Session{ID: sessionID}
	if err := sessionModel.GetByID(&session).Error; err != nil {
		return nil, error2.FromDatabaseError(err)
	}

	return &dto.Session{
		TinySession: dto.TinySession{
			ID:        session.ID,
			IsClosed:  session.IsClosed,
			CreatedAt: session.CreatedAt,
			UpdatedAt: session.UpdatedAt,
		},
		Class: dto.TinyClass{
			ID:   session.Class.ID,
			Name: session.Class.Name,
			Year: session.Class.Year,
		},
	}, nil
}

// CreateSession creates a new session for a class
func CreateSession(tx *gorm.DB, session dto.CreateSession) (*dto.TinySession, error) {
	sessionModel := model.NewSessionModel(tx)

	s := model.Session{
		ClassID:  session.ClassID,
		Password: session.Password,
		IsClosed: false,
	}
	if err := sessionModel.Create(&s); err != nil {
		return nil, error2.FromDatabaseError(err)
	}

	return &dto.TinySession{
		ID:        s.ID,
		IsClosed:  s.IsClosed,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}, nil
}

// CloseSession closes a session
func CloseSession(tx *gorm.DB, classID uint64, sessionID uuid.UUID) error {
	sessionModel := model.NewSessionModel(tx)

	s := model.Session{ID: sessionID, ClassID: classID, IsClosed: true}
	tx = sessionModel.Close(&s)
	if tx.Error != nil {
		return error2.FromDatabaseError(tx.Error)
	}

	if tx.RowsAffected == 0 {
		return error2.NotFoundError(fmt.Sprintf("session '%s' not found for class '%d'", sessionID, classID))
	}

	return nil
}

// DeleteSession deletes a session
func DeleteSession(tx *gorm.DB, sessionID uuid.UUID) error {
	sessionModel := model.NewSessionModel(tx)

	s := model.Session{ID: sessionID}
	if err := sessionModel.Delete(&s); err != nil {
		return error2.FromDatabaseError(err)
	}

	return nil
}

// DeleteSessionsByClassID deletes all sessions of a class
func DeleteSessionsByClassID(tx *gorm.DB, classID uint64) error {
	sessionModel := model.NewSessionModel(tx)

	if err := sessionModel.DeleteAll(classID); err != nil {
		return error2.FromDatabaseError(err)
	}

	return nil
}

// AddStudentToSession adds a student to a session
func AddStudentToSession(tx *gorm.DB, sessionID uuid.UUID, studentID, classID uint64) error {
	//userModel := model.NewUserModel(tx)
	//
	//if err := userModel.AddUserToSession(&model.User{ID: studentID, ClassID: &classID}, &model.Session{
	//	ID:      sessionID,
	//	ClassID: classID,
	//}); err != nil {
	//	return error2.FromDatabaseError(err)
	//}

	return nil
}

// RemoveStudentFromSession removes a student from a session
func RemoveStudentFromSession(tx *gorm.DB, sessionID uuid.UUID, studentID, classID uint64) error {
	//userModel := model.NewUserModel(tx)
	//
	//if err := userModel.RemoveUserFromSession(&model.User{ID: studentID, ClassID: &classID}, &model.Session{
	//	ID:      sessionID,
	//	ClassID: classID,
	//}); err != nil {
	//	return error2.FromDatabaseError(err)
	//}

	return nil
}

// VerifySession verifies a session password
func VerifySession(tx *gorm.DB, sessionID uuid.UUID, password string) (*model.Session, error) {
	sessionModel := model.NewSessionModel(tx)

	s := model.Session{ID: sessionID}
	tx = sessionModel.GetByID(&s)
	if tx.Error != nil {
		return nil, error2.FromDatabaseError(tx.Error)
	}

	if s.IsClosed {
		return nil, error2.BadRequestError("session is closed", nil)
	}

	if s.Password != password {
		return nil, error2.BadRequestError("incorrect password", nil)
	}

	return &s, nil
}
