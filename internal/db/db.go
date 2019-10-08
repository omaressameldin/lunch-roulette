package db

import (
	"fmt"
	"log"
	"strconv"
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

func (d *DB) AddNextRoundDate(t time.Time) error {
	if t.Before(time.Now()) {
		return fmt.Errorf("Next Round Date has to be in the future!")
	}

	err := d.database.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(d.bucket)
		return b.Put([]byte(latestRoundKey), []byte(t.Format(timeLayout)))
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) GetNextRoundDate() (*time.Time, error) {
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

	t, err := time.Parse(timeLayout, string(round))
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

func (d *DB) AddFrequencyPerMonth(frequency int) error {
	err := d.database.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(d.bucket)
		return b.Put([]byte(frequencyPerMonthKey), []byte(strconv.Itoa(frequency)))
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) GetFrequencyPerMonth() (*int, error) {
	var frequency []byte
	err := d.database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(d.bucket))
		frequency = b.Get([]byte(frequencyPerMonthKey))
		return nil
	})
	if err != nil {
		return nil, err
	}

	freqInt, err := strconv.Atoi(string(frequency))
	if err != nil {
		return nil, err
	}
	return &freqInt, nil
}

func (d *DB) AddGroupSize(groupSize int) error {
	err := d.database.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(d.bucket)
		return b.Put([]byte(groupSizeKey), []byte(strconv.Itoa(groupSize)))
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) GetGroupSize() (*int, error) {
	var frequency []byte
	err := d.database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(d.bucket))
		frequency = b.Get([]byte(groupSizeKey))
		return nil
	})
	if err != nil {
		return nil, err
	}

	freqInt, err := strconv.Atoi(string(frequency))
	if err != nil {
		return nil, err
	}
	return &freqInt, nil
}
