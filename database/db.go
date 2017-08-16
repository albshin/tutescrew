package database

import (
	"log"

	"github.com/boltdb/bolt"
)

var (
	// DB is the database singleton
	DB *bolt.DB
)

// User defines the user model
type User struct {
	RCSID string
}

// Connect initializes the database
func Connect() {
	DB, err := bolt.Open("tutescrew.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	err = DB.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("DB"))
		if err != nil {
			return err
		}
		_, err = root.CreateBucketIfNotExists([]byte("USERS"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal("Unabe to setup buckets")
	}
}

// AddUser adds a new user to the database
func AddUser(db *bolt.DB, did, rcsid string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DB")).Bucket([]byte("USERS"))
		err := b.Put([]byte(did), []byte(rcsid))
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
