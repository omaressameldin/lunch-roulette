package commands

import "github.com/nlopes/slack"

func CancelButton() slack.AttachmentAction {
	return slack.AttachmentAction{
		Name:  CancelCallback,
		Text:  cancelText,
		Type:  "button",
		Style: "danger",
	}
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
