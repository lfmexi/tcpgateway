package session

import "github.com/lfmexi/tcpgateway/session/model"

// Service interface that represents the session service
type Service interface {
	CreateSession(string, string) (*model.Session, error)
	DisableSession(*model.Session) error
}
