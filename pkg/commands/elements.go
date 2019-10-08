package commands

import "github.com/nlopes/slack"

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
