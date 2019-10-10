package commands

import (
	"fmt"

	"github.com/divan/num2words"
	"github.com/nlopes/slack"
	"github.com/shomali11/slacker"
)

func CancelButton() slack.BlockElement {
	return slack.NewButtonBlockElement(
		CancelValue,
		CancelValue,
		slack.NewTextBlockObject("plain_text", cancelText, false, false),
	)
}

func Select(
	text string,
	key string,
	options []slack.AttachmentActionOption,
) slack.AttachmentAction {
	return slack.AttachmentAction{
		Name:    key,
		Text:    text,
		Type:    "select",
		Options: options,
	}
}

func DangerMessage(text string) slack.Attachment {
	return slack.Attachment{
		Text:  text,
		Color: colorDanger,
	}
}

func PendingMessage(text string) slack.Attachment {
	return slack.Attachment{
		Text:  text,
		Color: colorPending,
	}
}

func SuccessMessage(text string) slack.Attachment {
	return slack.Attachment{
		Text:  text,
		Color: colorSuccess,
	}
}

func numberSelect(
	min int,
	max int,
	blockID string,
	actionID string,
	title string,
) []slack.Block {
	elements := make([]slack.BlockElement, 0, max-min+1)
	for i := min; i <= max; i++ {
		elements = append(elements, slack.NewButtonBlockElement(
			fmt.Sprintf("%s%s%d", actionID, NumberActionSeparator, i),
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
			blockID,
			elements...,
		),
	}
}

func sendError(channelID string, bot *slacker.Slacker, err error) {
	bot.RTM().PostMessage(channelID, slack.MsgOptionAttachments(
		DangerMessage(err.Error()),
	))
}
