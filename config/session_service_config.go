package config

import (
	"bitbucket.org/challengerdevs/gpsdriver/session"
	"bitbucket.org/challengerdevs/gpsdriver/session/repository"
)

func sessionRepository() repository.SessionRepository {
	return repository.NewMongoSessionRepository(sessionDb())
}

func sessionService() session.Service {
	return session.NewPersistentSessionProviderService(sessionRepository())
}
