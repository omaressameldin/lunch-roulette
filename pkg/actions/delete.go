package actions

import (
	"net/http"

	"github.com/nlopes/slack"
	"github.com/omaressameldin/lunch-roulette/internal/db"
	"github.com/omaressameldin/lunch-roulette/pkg/commands"
)

func deleteChannel(
	database *db.DB,
	responseURL string,
	w http.ResponseWriter,
	selectedChannel string,
) {
	sendPendingResponse(responseURL, w)
	err := database.DeleteBotChannel(selectedChannel)
	if err != nil {
		sendCancelResponse(responseURL, w, err.Error())
		return
	}

	sendReply(responseURL, w, Reply{
		Attachments: []slack.Attachment{
			commands.SuccessMessage(commands.DeleteSuccess),
		},
	})
}
