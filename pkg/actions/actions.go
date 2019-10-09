package actions

import (
	"fmt"
	"log"
	"net/http"

	"github.com/omaressameldin/lunch-roulette/internal/env"
	"github.com/omaressameldin/lunch-roulette/pkg/bot"
	"github.com/omaressameldin/lunch-roulette/pkg/commands"
)

func HandleActions(bot *bot.Bot) {
	port := env.GetActionPort()
	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			payload, _ := unmarshalPayload(r)

			if mightCancel(payload, w) {
				return
			}
			for _, block := range payload.ActionCallback.BlockActions {
				switch block.BlockID {
				case commands.SelectChannelBlockID:
					{
						selectChannel(bot.DB, payload.ResponseURL, w, block.SelectedChannel)
					}
				case commands.FirstRoundStartBlockID:
					{
						setFirstRoundDate(bot.DB, payload.ResponseURL, w, block.SelectedDate)
					}
				case commands.FerquencyPerMonthBlockID:
					{
						setFrequencyPerMonth(bot.DB, payload.ResponseURL, w, block.Value)
					}
				case commands.GroupSizeBlockID:
					{
						setGroupSize(bot.SlackBot, bot.DB, payload.ResponseURL, w, block.Value)
					}
				}
			}
		}
	})
	log.Printf("listening on port %s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
