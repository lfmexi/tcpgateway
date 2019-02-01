package repository

import "github.com/lfmexi/tcpgateway/session/model"

// SessionRepository is the session repository
type SessionRepository interface {
	Insert(*model.Session) error
	Update(*model.Session) error
}
