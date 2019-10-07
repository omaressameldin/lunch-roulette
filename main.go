package main

import (
	"log"

	"github.com/omaressameldin/lunch-roulette/internal/db"
	"github.com/omaressameldin/lunch-roulette/internal/env"
	"github.com/omaressameldin/lunch-roulette/pkg/bot"
)

func main() {
	env.ValidateEnvKeys()

	d := db.OpenDB()
	defer d.CloseDB()

	b, err := bot.CreateBot(d)
	if err != nil {
		log.Fatal(err)
	}
	b.StartListening()
}
