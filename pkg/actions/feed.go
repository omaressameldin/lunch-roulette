package actions

import (
	"fmt"
	"net/http"
	"strconv"
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

func setFirstRoundDate(
	database *db.DB,
	responseURL string,
	w http.ResponseWriter,
	selectedDate string,
) {
	nextRound, err := time.Parse(
		timeLayout,
		fmt.Sprintf("%s %s", selectedDate, commands.RoundTime),
	)
	if err != nil {
		sendCancelResponse(responseURL, w, err.Error())
	}
	sendPendingResponse(responseURL, w)
	err = database.AddNextRoundDate(nextRound)
	if err != nil {
		sendCancelResponse(responseURL, w, err.Error())
	}
	sendReply(responseURL, w, Reply{
		Blocks: commands.FrequencyPerMonth(),
	})
}

func setFrequencyPerMonth(
	database *db.DB,
	responseURL string,
	w http.ResponseWriter,
	selectedFrequency string,
) {
	frequency, err := strconv.Atoi(selectedFrequency)
	if err != nil {
		sendCancelResponse(responseURL, w, err.Error())
	}
	sendPendingResponse(responseURL, w)
	err = database.AddFrequencyPerMonth(frequency)
	if err != nil {
		sendCancelResponse(responseURL, w, err.Error())
	}
}
