package commands

import (
	"github.com/nlopes/slack"
	"github.com/omaressameldin/lunch-roulette/internal/utils"
	"github.com/shomali11/slacker"
)

func AddExcludeCmd(bot *slacker.Slacker) {
	bot.Command(excludeCmd, exclude())
}

func exclude() *slacker.CommandDefinition {
	return &slacker.CommandDefinition{
		Description: excludeDesc,
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			channel := request.Event().Channel
			rtm := response.RTM()

			utils.ReplyWithError(selectScheduleForExclusion(channel, rtm), excludeError, response)
		},
	}
}

func selectScheduleForExclusion(channel string, rtm *slack.RTM) error {
	rtm.PostMessage(channel, slack.MsgOptionBlocks(
		slack.NewContextBlock(
			"",
			[]slack.MixedElement{
				slack.NewTextBlockObject("mrkdwn", excludeQuestion, false, false),
			}...,
		),
		slack.NewContextBlock(
			"",
			[]slack.MixedElement{
				slack.NewTextBlockObject("mrkdwn", excludeWarning, false, false),
			}...,
		),
		slack.NewActionBlock(
			ExcludeChannelBlockID,
			slack.NewOptionsSelectBlockElement(
				"channels_select",
				slack.NewTextBlockObject("plain_text", excludeChannelPlaceholder, false, false),
				excludedChannelKey,
			),
			CancelButton(),
		),
	))

	return nil
}
