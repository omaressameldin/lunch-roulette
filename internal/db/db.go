package db

import (
	"encoding/binary"
	"errors"
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

func (d *DB) GetBotChannels() ([]string, error) {
	var channels []string
	err := d.database.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			channels = append(channels, string(name))
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func (d *DB) AddNextRoundDate(channelID string, t time.Time) error {
	if t.Before(time.Now()) {
		return fmt.Errorf("Next Round Date has to be in the future!")
	}

	return d.database.Update(func(tx *bolt.Tx) error {
		b, err := getScheduleBucket(channelID, tx)
		if err != nil {
			return err
		}

		return b.Put([]byte(keyNextRound), []byte(t.Format(timeLayout)))
	})
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
	if frequency < 1 || frequency > 30 {
		return fmt.Errorf("Frequency has to be between 1 and 30! ")
	}

	return d.database.Update(func(tx *bolt.Tx) error {
		b, err := getScheduleBucket(channelID, tx)
		if err != nil {
			return err
		}

		return b.Put([]byte(keyFrequencyPerMonth), []byte(strconv.Itoa(frequency)))
	})
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
	if groupSize < 2 {
		return fmt.Errorf("Group size can not be less than 2! ")
	}

	return d.database.Update(func(tx *bolt.Tx) error {
		b, err := getScheduleBucket(channelID, tx)
		if err != nil {
			return err
		}

		return b.Put([]byte(keyGroupSize), []byte(strconv.Itoa(groupSize)))
	})
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

func (d *DB) AddMembers(channelID string, members []string) error {
	return d.database.Batch(func(tx *bolt.Tx) error {
		b, err := getScheduleBucket(channelID, tx)
		if err != nil {
			return err
		}

		membersBucket, err := b.CreateBucketIfNotExists([]byte(keyMembers))
		for _, member := range members {
			id64, _ := membersBucket.NextSequence()
			key := itob(int(id64))
			b.Put(key, []byte(member))
		}

		return nil
	})
}

func (d *DB) DeleteAllSelectedMembers(channelID string) error {
	return d.database.Update(func(tx *bolt.Tx) error {
		b, err := getScheduleBucket(channelID, tx)
		if err != nil {
			return err
		}

		return b.DeleteBucket([]byte(keyMembers))
	})
}

func (d *DB) AllMembers(channelID string) ([]string, error) {
	var members []string
	err := d.database.View(func(tx *bolt.Tx) error {
		b, err := getScheduleBucket(channelID, tx)
		if err != nil {
			return err
		}
		membersBucket := b.Bucket([]byte(keyMembers))
		if membersBucket == nil {
			return errors.New("Members bucket does not exist!")
		}
		c := membersBucket.Cursor()
		for k, member := c.First(); k != nil; k, member = c.Next() {
			members = append(members, string(member))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return members, nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func getScheduleBucket(channelID string, tx *bolt.Tx) (*bolt.Bucket, error) {
	bucketName := []byte(channelID)
	b := tx.Bucket(bucketName)
	if b == nil {
		return nil, bucketError
	}

	return b, nil
}
