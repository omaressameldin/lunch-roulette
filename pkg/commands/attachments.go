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

