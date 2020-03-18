package session

import (
	"log"

	"github.com/boltdb/bolt"

	"ustb_sso/env"
)

var (
	DBCookies = []byte("cookies")
)

var db *bolt.DB

func init() {
	var err error
	db, err = bolt.Open(env.BoltDB, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var tables = [][]byte{DBCookies}
	err = db.Update(func(tx *bolt.Tx) error {
		for _, l := range tables {
			_, err := tx.CreateBucketIfNotExists(l)
			if err != nil {
				panic(err)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func GetDb() *bolt.DB {
	return db
}
