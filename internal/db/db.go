package db

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/omaressameldin/lunch-roulette/internal/env"
)

type DB struct {
	database *bolt.DB
	bucket   []byte
}

func OpenDB() *DB {
	dbName := env.GetDBName()
	bucketName := env.GetDBBucket()
	bucket := []byte(bucketName)
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		return err
	})

	return &DB{
		database: db,
		bucket:   bucket,
	}
}

func (d *DB) CloseDB() {
	d.database.Close()
}

func (d *DB) AddNextRoundDate(layout string, t time.Time) error {
	err := d.database.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(d.bucket)
		return b.Put([]byte(latestRoundKey), []byte(t.Format(layout)))
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) GetNextRoundDate(layout string) (*time.Time, error) {
	var round []byte
	err := d.database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(d.bucket))
		round = b.Get([]byte(latestRoundKey))
		return nil
	})
	if err != nil {
		return nil, err
	}
	if string(round) == "" {
		return nil, nil
	}

	t, err := time.Parse(layout, string(round))
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (d *DB) AddBotChannel(channelID string) error {
	err := d.database.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(d.bucket)
		return b.Put([]byte(botChannelKey), []byte(channelID))
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) GetBotChannel() (*string, error) {
	var channelID []byte
	err := d.database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(d.bucket))
		channelID = b.Get([]byte(botChannelKey))
		return nil
	})
	if err != nil {
		return nil, err
	}
	stringifiedChannelID := string(channelID)
	if stringifiedChannelID == "" {
		return nil, nil
	}

	return &stringifiedChannelID, nil
}
