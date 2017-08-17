package database

import (
	"log"

	"github.com/boltdb/bolt"
)

var (
	// DB is the database reference
	DB *bolt.DB
)

// User defines the user model
type User struct {
	RCSID string
}

// Connect initializes the database
func Connect() {
	db, err := bolt.Open("tutescrew.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	DB = db

	err = DB.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("DB"))
		if err != nil {
			return err
		}
		_, err = root.CreateBucketIfNotExists([]byte("STUDENTS"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal("Unable to setup buckets")
	}
}

// AddStudent adds a new student to the database
func AddStudent(rcsid, did string) error {
	err := DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DB")).Bucket([]byte("STUDENTS"))
		err := b.Put([]byte(rcsid), []byte(did))
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// IsRegistered checks if a RCSID is registered
func IsRegistered(rcsid, did string) bool {
	var reg bool
	err := DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DB")).Bucket([]byte("STUDENTS"))
		v := b.Get([]byte(rcsid))
		if v != nil {
			reg = true
			return nil
		}
		reg = false
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return reg
}
