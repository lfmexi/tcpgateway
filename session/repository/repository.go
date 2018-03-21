package repository

import "bitbucket.org/challengerdevs/gpsdriver/session/model"

// SessionRepository is the session repository
type SessionRepository interface {
	Insert(*model.Session) error
	Update(*model.Session) error
}
