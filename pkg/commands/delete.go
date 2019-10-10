package commands

import (
	"github.com/nlopes/slack"
	"github.com/omaressameldin/lunch-roulette/internal/utils"
	"github.com/shomali11/slacker"
)

func AddDeleteCmd(bot *slacker.Slacker) {
	bot.Command(deleteCmd, deleteSchedule())
}

func deleteSchedule() *slacker.CommandDefinition {
	return &slacker.CommandDefinition{
		Description: deleteDesc,
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			channel := request.Event().Channel
			rtm := response.RTM()

			utils.ReplyWithError(selectScheduleForDeletion(channel, rtm), deleteError, response)
		},
	}
}

func selectScheduleForDeletion(channel string, rtm *slack.RTM) error {
	rtm.PostMessage(channel, slack.MsgOptionBlocks(
		slack.NewContextBlock(
			"",
			[]slack.MixedElement{
				slack.NewTextBlockObject("mrkdwn", selectDeletedQuestion, false, false),
			}...,
		),
		slack.NewContextBlock(
			"",
			[]slack.MixedElement{
				slack.NewTextBlockObject("mrkdwn", selectDeletedWarning, false, false),
			}...,
		),
		slack.NewActionBlock(
			SelectDeletedBlockID,
			slack.NewOptionsSelectBlockElement(
				"channels_select",
				slack.NewTextBlockObject("plain_text", selectDeletedPlaceholder, false, false),
				deletedKey,
			),
			CancelButton(),
		),
	))
	return nil
}
