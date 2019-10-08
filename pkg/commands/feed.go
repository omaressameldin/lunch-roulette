package commands

import (
	"github.com/nlopes/slack"
	"github.com/omaressameldin/lunch-roulette/internal/utils"
	"github.com/shomali11/slacker"
)

func AddFeedCmd(bot *slacker.Slacker) {
	bot.Command(feedCmd, feed())
}

func feed() *slacker.CommandDefinition {
	return &slacker.CommandDefinition{
		Description: feedDesc,
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			channel := request.Event().Channel
			rtm := response.RTM()

			utils.ReplyWithError(selectFoodChannel(channel, rtm), feedError, response)
		},
	}
}

func selectFoodChannel(channel string, rtm *slack.RTM) error {
	rtm.PostMessage(channel, slack.MsgOptionBlocks(
		slack.NewContextBlock(
			"",
			[]slack.MixedElement{
				slack.NewTextBlockObject("mrkdwn", selectChannelQuestion, false, false),
			}...,
		),
		slack.NewActionBlock(
			SelectChannelBlockId,
			slack.NewOptionsSelectBlockElement(
				"channels_select",
				slack.NewTextBlockObject("plain_text", selectChannelPlaceholder, false, false),
				FoodChannelKey,
			),
			CancelButton(),
		),
	))
	return nil
	}

func PendingResponse() slack.Attachment {
	return PendingMessage(addingToDatabase)
	}

func AddChannelToDB(d *db.DB, channelID string) error {
	return d.AddBotChannel(channelID)
}

func FirstRoundDate() []slack.Block {

	return []slack.Block{
		slack.NewContextBlock(
			"",
			[]slack.MixedElement{
				slack.NewTextBlockObject("mrkdwn", firstRoundStartText, false, false),
			}...,
			),
		slack.NewActionBlock(
			FirstRoundStartBlockId,
			slack.NewDatePickerBlockElement(firstRoundKey),
			CancelButton(),
		),
	}
}

	rtm.PostMessage(channel, slack.MsgOptionAttachments(attachment))
	return nil
}
