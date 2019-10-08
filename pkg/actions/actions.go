package actions

import (
	"fmt"
	"log"
	"net/http"

	"github.com/omaressameldin/lunch-roulette/internal/env"
	"github.com/omaressameldin/lunch-roulette/pkg/bot"
)

func HandleActions(bot *bot.Bot) {
	port := env.GetActionPort()
	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			payload, _ := unmarshalPayload(r)

			if mightCancel(payload, w) {
				return
			}
		}
	})
	log.Printf("listening on port %s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
