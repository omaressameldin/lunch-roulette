package commands

import (
	"log"

	"github.com/omaressameldin/lunch-roulette/internal/db"
	"github.com/omaressameldin/lunch-roulette/internal/env"
	"github.com/omaressameldin/lunch-roulette/internal/utils"
	"github.com/shomali11/slacker"
)

func AddInitCmd(db *db.DB, bot *slacker.Slacker) {
	bot.Init(organizeLunch(db, bot))
}

func organizeLunch(d *db.DB, bot *slacker.Slacker) func() {
	return func() {
		bot.Client().ConnectRTM()
		for {
			waitForRound(d)
		}
	}
}

func waitForRound(d *db.DB) {
	freq := env.GetRoundFrequencyPerMonth()
	firstRound := env.GetFirstRoundDate()
	latestRound, err := d.GetRound(env.TimeLayout)
	if err != nil {
		log.Fatal(err)
	}

	if latestRound == nil {
		latestRound = &firstRound
	}
	// add freq weeks to next round
	err = utils.SleepTill(*latestRound)
	d.CreateRound(env.TimeLayout, latestRound.AddDate(0, 0, 7*freq))
	if err != nil {
		log.Fatal(err)
	}
}
