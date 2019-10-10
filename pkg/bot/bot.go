package bot

import (
	"context"
	"log"

	"github.com/omaressameldin/lunch-roulette/internal/db"
	"github.com/omaressameldin/lunch-roulette/internal/env"
	"github.com/omaressameldin/lunch-roulette/pkg/commands"
	"github.com/shomali11/slacker"
)

type Bot struct {
	SlackBot *slacker.Slacker
	DB       *db.DB
}

func CreateBot(d *db.DB) (*Bot, error) {
	t := env.GetToken()

	return &Bot{
		SlackBot: slacker.NewClient(t),
		DB:       d,
	}, nil
}

func (b *Bot) StartListening() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	commands.AddInitCmd(b.DB, b.SlackBot)
	commands.AddFeedCmd(b.SlackBot)
	commands.AddDeleteCmd(b.SlackBot)

	log.Println("Listening...")
	err := b.SlackBot.Listen(ctx)
	if err != nil {
		return err
	}

	return nil
}
