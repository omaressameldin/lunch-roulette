package commands

import (
	"fmt"

	"github.com/nlopes/slack"
	"github.com/omaressameldin/lunch-roulette/internal/db"
	"github.com/omaressameldin/lunch-roulette/internal/utils"
	"github.com/shomali11/slacker"
)

func AddFeedCmd(bot *slacker.Slacker) {
	bot.Command(feedCmd, feed())
}

func feed() *slacker.CommandDefinition {
	return &slacker.CommandDefinition{
		Description:       feedDesc,
		AuthorizationFunc: authFunction,
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
		slack.NewContextBlock(
			"",
			[]slack.MixedElement{
				slack.NewTextBlockObject("mrkdwn", selectChannelWarning, false, false),
			}...,
		),
		slack.NewActionBlock(
			SelectChannelBlockID,
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

func FirstRoundDate(channelID string) []slack.Block {
	return []slack.Block{
		slack.NewContextBlock(
			"",
			[]slack.MixedElement{
				slack.NewTextBlockObject("mrkdwn", firstRoundStartText, false, false),
			}...,
		),
		slack.NewActionBlock(
			FirstRoundStartBlockID,
			slack.NewDatePickerBlockElement(channelID),
			CancelButton(),
		),
	}
}

func FrequencyPerMonth(channelID string) []slack.Block {
	return numberSelect(1, 4, FerquencyPerMonthBlockID, channelID, frequencyPerMonthText)
}

func GroupSize(channelID string) []slack.Block {
	return numberSelect(2, 6, GroupSizeBlockID, channelID, groupSizeText)
}

func DoneText(channelID string, bot *slacker.Slacker) (string, error) {
	lunchInfo, err := db.GetLunchInfo(channelID)
	if err != nil {
		return "", err
	}

	channelInfo, err := bot.RTM().GetChannelInfo(channelID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"🍔 Everything is Done! starting *%s*, *%d* members will be paired *%d* times a month for lunch. Stats will be posted on *%s*",
		lunchInfo.NextRoundDate.Format(timeLayout),
		lunchInfo.GroupSize,
		lunchInfo.FrequencyPerMonth,
		channelInfo.Name,
	), nil
}
