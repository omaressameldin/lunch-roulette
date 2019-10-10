package commands

import (
	"github.com/omaressameldin/lunch-roulette/internal/db"
	"github.com/shomali11/slacker"
)

func AddInitCmd(db *db.DB, bot *slacker.Slacker) {
	bot.Init(func() { organzieLunches(db, bot) })
}
