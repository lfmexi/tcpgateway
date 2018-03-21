package repository

import (
	"bitbucket.org/challengerdevs/gpsdriver/session/model"
	"gopkg.in/mgo.v2"
)

const (
	collection = "sessions"
)

// NewMongoSessionRepository creates a new mongodb based session repository
func NewMongoSessionRepository(db *mgo.Database) SessionRepository {
	return &mongoSessionRepository{
		db,
	}
}

type mongoSessionRepository struct {
	db *mgo.Database
}

func (m *mongoSessionRepository) Insert(session *model.Session) error {
	return m.db.C(collection).Insert(session)
}

func (m *mongoSessionRepository) Update(sess *model.Session) error {
	return m.db.C(collection).UpdateId(sess.ID, sess)
}
