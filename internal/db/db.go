package db

import (
	"log"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/omaressameldin/lunch-roulette/internal/env"
)

type Lunch struct {
	ChannelID         string     `pg:"channel_id,pk"`
	GroupSize         int        `pg:"group_size"`
	FrequencyPerMonth int        `pg:"frequency_per_month"`
	NextRoundDate     *time.Time `pg:"next_round_date"`
}

type LunchMember struct {
	ChannelID string `pg:"channel_id"`
	MemberID  string `pg:"member_id"`
}

func openDB() *pg.DB {
	url := env.GetDatabaseUrl()
	options, err := pg.ParseURL(url)
	if err != nil {
		log.Fatal(err)
	}
	return pg.Connect(options)
}

func AddLunchChannel(channelID string) error {
	db := openDB()
	defer db.Close()

	return db.Insert(&Lunch{
		ChannelID: channelID,
	})
}

func IsChannelLinked(channelID string) error {
	db := openDB()
	defer db.Close()

	return db.Select(&Lunch{
		ChannelID: channelID,
	})
}

func DeleteLunchChannel(channelID string) error {
	db := openDB()
	defer db.Close()

	return db.Delete(&Lunch{
		ChannelID: channelID,
	})
}

func GetLunchChannels() ([]Lunch, error) {
	db := openDB()
	defer db.Close()

	var lunches []Lunch
	err := db.Model(&lunches).Select()
	if err != nil {
		return nil, err
	}
	return lunches, nil
}

func AddNextRoundDate(channelID string, t time.Time) error {
	db := openDB()
	defer db.Close()
	lunch, err := GetLunchInfo(channelID)
	if err != nil {
		return err
	}

	lunch.NextRoundDate = &t
	return db.Update(lunch)
}

func AddFrequencyPerMonth(channelID string, frequency int) error {
	db := openDB()
	defer db.Close()

	lunch, err := GetLunchInfo(channelID)
	if err != nil {
		return err
	}
	lunch.FrequencyPerMonth = frequency
	return db.Update(lunch)
}

func AddGroupSize(channelID string, groupSize int) error {
	db := openDB()
	defer db.Close()
	lunch, err := GetLunchInfo(channelID)
	if err != nil {
		return err
	}
	lunch.GroupSize = groupSize
	return db.Update(lunch)
}

func GetLunchInfo(channelID string) (*Lunch, error) {
	db := openDB()
	defer db.Close()

	lunch := &Lunch{
		ChannelID: channelID,
	}
	err := db.Select(lunch)
	if err != nil {
		return nil, err
	}

	return lunch, err
}

func AddMembers(channelID string, members []string) error {
	db := openDB()
	defer db.Close()

	transacton, err := db.Begin()
	if err != nil {
		return err
	}
	for _, member := range members {
		err = db.Insert(&LunchMember{
			ChannelID: channelID,
			MemberID:  member,
		})
		if err != nil {
			return err
		}
	}
	err = transacton.Commit()
	if err != nil {
		return err
	}

	return transacton.Close()
}

func DeleteAllSelectedMembers(channelID string) error {
	db := openDB()
	defer db.Close()
	filter := func(q *orm.Query) (*orm.Query, error) {
		q = q.Where("channel_id = ?", channelID)
		return q, nil
	}
	var members []LunchMember
	_, err := db.Model(&members).Apply(filter).Delete()

	return err
}

func AllMembers(channelID string) ([]string, error) {
	db := openDB()
	defer db.Close()

	var members []LunchMember
	var memberIDs []string

	filter := func(q *orm.Query) (*orm.Query, error) {
		q = q.Where("channel_id = ?", channelID)
		return q, nil
	}
	err := db.Model(&members).Apply(filter).Select()
	if err != nil {
		return nil, err
	}
	for _, member := range members {
		memberIDs = append(memberIDs, member.MemberID)
	}

	return memberIDs, nil
}
