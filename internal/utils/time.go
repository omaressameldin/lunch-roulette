package utils

import (
	"errors"
	"time"
)

func timeTill(t time.Time) (*time.Duration, error) {
	n := time.Now()
	d := t.Sub(n)
	if n.After(t) {
		return nil, errors.New("Time has to be in the future!")
	}
	return &d, nil
}

func SleepTill(t time.Time) error {
	now := time.Now()
	if t.After(now) {
		d, err := timeTill(t)
		if err != nil {
			return err
		}
		time.Sleep(*d)
	}
	return nil

}
