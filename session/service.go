package session

import "bitbucket.org/challengerdevs/tcpgateway/session/model"

// Service interface that represents the session service
type Service interface {
	CreateSession(string, string) (*model.Session, error)
	DisableSession(*model.Session) error
}
