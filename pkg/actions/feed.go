package actions

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/nlopes/slack"
	"github.com/shomali11/slacker"

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
	err := database.AddBotChannel(selectedChannel)
	if err != nil {
		sendCancelResponse(responseURL, w, err.Error())
		return
	}
	sendReply(responseURL, w, Reply{
		Blocks: commands.FirstRoundDate(selectedChannel),
	})
}

func setFirstRoundDate(
	database *db.DB,
	responseURL string,
	w http.ResponseWriter,
	channelID string,
	selectedDate string,
) {
	nextRound, err := time.Parse(
		timeLayout,
		fmt.Sprintf("%s %s", selectedDate, commands.RoundTime),
	)
	if err != nil {
		sendCancelResponse(responseURL, w, err.Error())
		return
	}
	sendPendingResponse(responseURL, w)
	err = database.AddNextRoundDate(channelID, nextRound)
	if err != nil {
		sendCancelResponse(responseURL, w, err.Error())
		return
	}
	sendReply(responseURL, w, Reply{
		Blocks: commands.FrequencyPerMonth(channelID),
	})
}

func setFrequencyPerMonth(
	database *db.DB,
	responseURL string,
	w http.ResponseWriter,
	channelID string,
	selectedFrequency string,
) {
	frequency, err := strconv.Atoi(selectedFrequency)
	if err != nil {
		sendCancelResponse(responseURL, w, err.Error())
		return
	}
	sendPendingResponse(responseURL, w)
	err = database.AddFrequencyPerMonth(channelID, frequency)
	if err != nil {
		sendCancelResponse(responseURL, w, err.Error())
		return
	}
	sendReply(responseURL, w, Reply{
		Blocks: commands.GroupSize(channelID),
	})
}

func setGroupSize(
	bot *slacker.Slacker,
	database *db.DB,
	responseURL string,
	w http.ResponseWriter,
	channelID string,
	selectedSize string,
) {
	size, err := strconv.Atoi(selectedSize)
	if err != nil {
		sendCancelResponse(responseURL, w, err.Error())
		return
	}
	sendPendingResponse(responseURL, w)
	err = database.AddGroupSize(channelID, size)
	if err != nil {
		sendCancelResponse(responseURL, w, err.Error())
		return
	}

	successMessage, err := commands.DoneText(database, channelID, bot)
	if err != nil {
		sendCancelResponse(responseURL, w, err.Error())
		return
	}
	sendReply(responseURL, w, Reply{
		Attachments: []slack.Attachment{
			commands.SuccessMessage(successMessage),
		},
	})
}
