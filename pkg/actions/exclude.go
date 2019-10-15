package actions

import (
	"net/http"

	"github.com/nlopes/slack"
	"github.com/omaressameldin/lunch-roulette/internal/db"
	"github.com/omaressameldin/lunch-roulette/internal/utils"
	"github.com/omaressameldin/lunch-roulette/pkg/commands"
	"github.com/shomali11/slacker"
)

func excludeChannel(
	database *db.DB,
	responseURL string,
	w http.ResponseWriter,
	selectedChannel string,
) {
	if err := database.IsChannelLinked(selectedChannel); err != nil {
		sendCancelResponse(responseURL, w, err.Error())
		return
	}

	sendReply(responseURL, w, Reply{
		Blocks: commands.SelectExcludedMember(selectedChannel),
	})
}

func excludeMember(
	database *db.DB,
	bot *slacker.Slacker,
	responseURL string,
	w http.ResponseWriter,
	selectedMember string,
	channelID string,
) {
	sendPendingResponse(responseURL, w)
	info, err := bot.Client().GetChannelInfo(channelID)
	if err != nil {
		sendCancelResponse(responseURL, w, err.Error())
		return
	}
	if !utils.Contains(info.Members, selectedMember) {
		sendCancelResponse(responseURL, w, "Member is not in channel!")
		return
	}

	if err := database.AddMembers(channelID, []string{selectedMember}); err != nil {
		sendCancelResponse(responseURL, w, err.Error())
	}

	sendReply(responseURL, w, Reply{
		Attachments: []slack.Attachment{
			commands.SuccessMessage(commands.ExcludeSuccess),
		},
	})
}
