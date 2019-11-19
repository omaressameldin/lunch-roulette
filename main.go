package main

import (
	"log"

	"github.com/omaressameldin/lunch-roulette/internal/env"
	"github.com/omaressameldin/lunch-roulette/pkg/actions"
	"github.com/omaressameldin/lunch-roulette/pkg/bot"
)

func main() {
	env.ValidateEnvKeys()

	b, err := bot.CreateBot()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		actions.HandleActions(b)
	}()
	b.StartListening()
}
