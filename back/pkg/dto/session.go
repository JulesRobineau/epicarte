package dto

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type TinySession struct {
	// ID is the id of the session
	ID uuid.UUID `json:"id"`
	// IsClosed is true if the session is closed
	IsClosed bool `json:"is_closed"`
	// CreatedAt is the creation date of the session
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the last update date of the session
	UpdatedAt time.Time `json:"updated_at"`
}

type Session struct {
	TinySession
	Class TinyClass `json:"class"`
}

type SessionList struct {
	// Sessions is the list of sessions
	Sessions []TinySession `json:"sessions"`
}

type CreateSession struct {
	// Password is the password of the session
	Password string `json:"password" binding:"required,min=8,max=255"`
	// ClassID is the id of the class
	ClassID uint64 `json:"-" uri:"class_id" uri:"class_id"`
}

type JoinSession struct {
	// Password is the password of the session
	Password string `json:"password" binding:"required,min=8,max=255" path:"password" form:"password" query:"password"`
}
