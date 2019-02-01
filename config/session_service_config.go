package config

import (
	"github.com/lfmexi/tcpgateway/session"
	"github.com/lfmexi/tcpgateway/session/repository"
)

func sessionRepository() repository.SessionRepository {
	return repository.NewMongoSessionRepository(sessionDb())
}

func sessionService() session.Service {
	return session.NewPersistentSessionProviderService(sessionRepository())
}
