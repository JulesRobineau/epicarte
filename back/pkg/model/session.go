package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	// ID is the session ID (UUID), overwriting the default ID field from gorm.Model
	ID uuid.UUID `gorm:"primary_key;type:uuid"`
	// Password is the password of the session (hashed)
	Password string `gorm:"type:varchar(255);not null"`
	// IsClosed is true if the session is closed
	IsClosed bool `gorm:"type:boolean;not null;default:false"`
	// ClassID is the ID of the class the session
	ClassID uint64 `gorm:"type:bigint;not null"`
	// Class is the class of the session
	Class *Class `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// Students is the list of students in the session
	Students []*User `gorm:"many2many:session_students;"`
}

// TableName overrides the default table name generated by GORM to be `sessions`
func (Session) TableName() string {
	return "sessions"
}

func (s *Session) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.NewV4()
	return
}

type SessionModel struct {
	Tx *gorm.DB
}

func NewSessionModel(tx *gorm.DB) *SessionModel {
	return &SessionModel{Tx: tx}
}

// GetByID gets a session by ID
func (m *SessionModel) GetByID(session *Session) *gorm.DB {
	return m.Tx.Preload("Class").First(&session)
}

// FindAll gets all sessions of a class
func (m *SessionModel) FindAll(classID uint64) ([]Session, error) {
	var sessions []Session
	m.Tx.Where("class_id = ?", classID).Find(&sessions)

	return sessions, nil
}

// Create creates a new session
func (m *SessionModel) Create(session *Session) error {
	return m.Tx.Create(session).Error
}

// Update updates a session
func (m *SessionModel) Update(session *Session) error {
	return m.Tx.Model(session).Updates(session).Error
}

// Close closes a session
func (m *SessionModel) Close(session *Session) *gorm.DB {
	return m.Tx.Model(session).Where(
		"id = ? AND class_id = ?", session.ID, session.ClassID,
	).Update("is_closed", true)
}

// Delete deletes a session
func (m *SessionModel) Delete(session *Session) error {
	return m.Tx.Delete(session).Error
}

// DeleteAll deletes all sessions of a class
func (m *SessionModel) DeleteAll(classID uint64) error {
	return m.Tx.Where("class_id = ?", classID).Delete(&Session{}).Error
}
