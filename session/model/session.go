package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Session represents an initialized session
type Session struct {
	ID         bson.ObjectId `bson:"_id"`
	Address    string        `bson:"address"`
	DeviceType string        `bson:"device_type"`
	CreatedAt  time.Time     `bson:"created_at"`
	UpdatedAt  time.Time     `bson:"updated_at"`
	Enabled    bool          `bson:"enabled"`
}

// NewSession creates a new session model
func NewSession(id bson.ObjectId, address string, deviceType string) *Session {
	return &Session{
		id,
		address,
		deviceType,
		time.Now(),
		time.Now(),
		true,
	}
}
