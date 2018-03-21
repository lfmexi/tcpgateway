package session

import (
	"time"

	"bitbucket.org/challengerdevs/gpsdriver/session/model"
	"bitbucket.org/challengerdevs/gpsdriver/session/repository"
	"gopkg.in/mgo.v2/bson"
)

type persistentSessionProviderService struct {
	sessionRepository repository.SessionRepository
}

// NewPersistentSessionProviderService creates a new session provider service
func NewPersistentSessionProviderService(repo repository.SessionRepository) Service {
	return &persistentSessionProviderService{
		repo,
	}
}

func (p *persistentSessionProviderService) CreateSession(address string, deviceType string) (*model.Session, error) {
	sess := model.NewSession(bson.NewObjectId(), address, deviceType)

	if err := p.sessionRepository.Insert(sess); err != nil {
		return nil, err
	}

	return sess, nil
}

func (p *persistentSessionProviderService) DisableSession(sess *model.Session) error {
	sess.UpdatedAt = time.Now()
	sess.Enabled = false

	return p.sessionRepository.Update(sess)
}
