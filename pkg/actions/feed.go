package actions

import (
	"fmt"
	"net/http"
	"time"

	"github.com/omaressameldin/lunch-roulette/internal/db"
	"github.com/omaressameldin/lunch-roulette/pkg/commands"
)

func selectChannel(
	database *db.DB,
	responseURL string,
	w http.ResponseWriter,
	selectedChannel string,
) {
	sendPendingResponse(responseURL, w)
	err := commands.AddChannelToDB(
		database,
		selectedChannel,
	)
	if err != nil {
		sendCancelResponse(responseURL, w, err.Error())
	}
	sendReply(responseURL, w, Reply{
		Blocks: commands.FirstRoundDate(),
	})
}
