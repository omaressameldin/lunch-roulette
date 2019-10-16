package actions

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/omaressameldin/lunch-roulette/internal/env"
	"github.com/omaressameldin/lunch-roulette/pkg/bot"
	"github.com/omaressameldin/lunch-roulette/pkg/commands"
)

func HandleActions(bot *bot.Bot) {
	port := env.GetActionPort()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Nothing to do here"))
	})

	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			payload, _ := unmarshalPayload(r)

			if mightCancel(payload, w) {
				return
			}
			for _, block := range payload.ActionCallback.BlockActions {
				switch block.BlockID {
				case commands.SelectDeletedBlockID:
					{
						deleteChannel(bot.DB, payload.ResponseURL, w, block.SelectedChannel)
					}
				case commands.ExcludeChannelBlockID:
					{
						excludeChannel(bot.DB, payload.ResponseURL, w, block.SelectedChannel)
					}
				case commands.ExcludeMemberBlockID:
					{
						excludeMember(
							bot.DB,
							bot.SlackBot,
							payload.ResponseURL,
							w,
							block.SelectedUser,
							block.ActionID,
						)
					}
				case commands.SelectChannelBlockID:
					{
						selectChannel(bot.DB, payload.ResponseURL, w, block.SelectedChannel)
					}
				case commands.FirstRoundStartBlockID:
					{
						setFirstRoundDate(
							bot.DB,
							payload.ResponseURL,
							w,
							block.ActionID,
							block.SelectedDate,
						)
					}
				case commands.FerquencyPerMonthBlockID:
					{
						channelID := strings.Split(block.ActionID, commands.NumberActionSeparator)[0]
						setFrequencyPerMonth(
							bot.DB,
							payload.ResponseURL,
							w,
							channelID,
							block.Value,
						)
					}
				case commands.GroupSizeBlockID:
					{
						channelID := strings.Split(block.ActionID, commands.NumberActionSeparator)[0]
						setGroupSize(
							bot.SlackBot,
							bot.DB,
							payload.ResponseURL,
							w,
							channelID,
							block.Value,
						)
					}
				}
			}
		}
	})
	log.Printf("listening on port %s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
