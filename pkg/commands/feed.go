package commands

import (
	"fmt"

	"github.com/divan/num2words"
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

func FrequencyPerMonth() []slack.Block {
	return numberSelect(1, 4, FerquencyPerMonthBlockId, frequencyPerMonthText)
}

func GroupSize() []slack.Block {
	return numberSelect(2, 6, GroupSizeBlockId, groupSizeText)
}

func numberSelect(min int, max int, blockId string, title string) []slack.Block {
	elements := make([]slack.BlockElement, 0, max-min+1)
	for i := min; i <= max; i++ {
		elements = append(elements, slack.NewButtonBlockElement(
			"",
			fmt.Sprintf("%d", i),
			slack.NewTextBlockObject("plain_text", num2words.Convert(i), false, false),
		))
	}
	elements = append(elements, CancelButton())

	return []slack.Block{
		slack.NewContextBlock(
			"",
			[]slack.MixedElement{
				slack.NewTextBlockObject("mrkdwn", title, false, false),
			}...,
		),
		slack.NewActionBlock(
			blockId,
			elements...,
		),
	}
}

func DoneText(database *db.DB) (string, error) {
	startDate, err := database.GetNextRoundDate()
	channel, err := database.GetBotChannel()
	frequency, err := database.GetFrequencyPerMonth()
	size, err := database.GetGroupSize()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"🍔 Everything is Done! starting *%s*, *%d* members will be paired *%d* times a month for lunch. Stats will be posted on *%s*",
		startDate.Format(timeLayout),
		*size,
		*frequency,
		*channel,
	), nil
}
