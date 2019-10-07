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
			payload, err := unmarshalPayload(r)
			if err != nil {
				sendCancelResponse(w, err.Error())
				return
			}

			if mightCancel(payload, w) {
				return
			}

			switch payload.CallbackID {
			}

		}
	})
	log.Printf("listening on port %s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
