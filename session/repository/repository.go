package repository

import "bitbucket.org/challengerdevs/tcpgateway/session/model"

// SessionRepository is the session repository
type SessionRepository interface {
	Insert(*model.Session) error
	Update(*model.Session) error
}
