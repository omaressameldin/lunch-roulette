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
}

func OpenDB() *DB {
	dbName := env.GetDBName()

	var err error
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &DB{
		database: db,
	}
}

func (d *DB) CloseDB() {
	d.database.Close()
}

func (d *DB) AddBotChannel(channelID string) error {
	return d.database.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(channelID))
		return err
	})
}

func (d *DB) AddNextRoundDate(channelID string, t time.Time) error {
	if t.Before(time.Now()) {
		return fmt.Errorf("Next Round Date has to be in the future!")
	}

	err := d.database.Update(func(tx *bolt.Tx) error {
		b, err := getScheduleBucket(channelID, tx)
		if err != nil {
			return err
		}

		return b.Put([]byte(keyNextRound), []byte(t.Format(timeLayout)))
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) GetNextRoundDate(channelID string) (*time.Time, error) {
	var round []byte
	err := d.database.View(func(tx *bolt.Tx) error {
		b, err := getScheduleBucket(channelID, tx)
		if err != nil {
			return err
		}

		round = b.Get([]byte(keyNextRound))
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

func (d *DB) AddFrequencyPerMonth(channelID string, frequency int) error {
	err := d.database.Update(func(tx *bolt.Tx) error {
		b, err := getScheduleBucket(channelID, tx)
		if err != nil {
			return err
		}

		return b.Put([]byte(keyFrequencyPerMonth), []byte(strconv.Itoa(frequency)))
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) GetFrequencyPerMonth(channelID string) (*int, error) {
	var frequency []byte
	err := d.database.View(func(tx *bolt.Tx) error {
		b, err := getScheduleBucket(channelID, tx)
		if err != nil {
			return err
		}

		frequency = b.Get([]byte(keyFrequencyPerMonth))
		return nil
	})

	if err != nil {
		return nil, err
	}
	if string(frequency) == "" {
		return nil, nil
	}

	freqInt, err := strconv.Atoi(string(frequency))
	if err != nil {
		return nil, err
	}
	return &freqInt, nil
}

func (d *DB) AddGroupSize(channelID string, groupSize int) error {
	err := d.database.Update(func(tx *bolt.Tx) error {
		b, err := getScheduleBucket(channelID, tx)
		if err != nil {
			return err
		}

		return b.Put([]byte(keyGroupSize), []byte(strconv.Itoa(groupSize)))
	})

	if err != nil {
		return err
	}

	return nil
}

func (d *DB) GetGroupSize(channelID string) (*int, error) {
	var size []byte
	err := d.database.View(func(tx *bolt.Tx) error {
		b, err := getScheduleBucket(channelID, tx)
		if err != nil {
			return err
		}
		size = b.Get([]byte(keyGroupSize))

		return nil
	})

	if err != nil {
		return nil, err
	}
	if string(size) == "" {
		return nil, nil
	}

	sizeInt, err := strconv.Atoi(string(size))
	if err != nil {
		return nil, err
	}
	return &sizeInt, nil
}

func getScheduleBucket(channelID string, tx *bolt.Tx) (*bolt.Bucket, error) {
	bucketName := []byte(channelID)
	b := tx.Bucket(bucketName)
	if b == nil {
		return nil, bucketError
	}

	return b, nil
}
