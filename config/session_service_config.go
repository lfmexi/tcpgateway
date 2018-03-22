package config

import (
	"bitbucket.org/challengerdevs/tcpgateway/session"
	"bitbucket.org/challengerdevs/tcpgateway/session/repository"
)

func sessionRepository() repository.SessionRepository {
	return repository.NewMongoSessionRepository(sessionDb())
}

func sessionService() session.Service {
	return session.NewPersistentSessionProviderService(sessionRepository())
}
