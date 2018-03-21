package config

import (
	"gopkg.in/mgo.v2"
)

var db *mgo.Database

func sessionDb() *mgo.Database {
	if db == nil {
		sess, err := mgo.Dial(configuration.Database.Host)
		if err != nil {
			panic(err)
		}

		db = sess.DB(configuration.Database.Name)
	}

	return db
}
