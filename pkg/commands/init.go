package commands

import (
	"github.com/shomali11/slacker"
)

func AddInitCmd(bot *slacker.Slacker) {
	bot.Init(func() { organzieLunches(bot) })
}
