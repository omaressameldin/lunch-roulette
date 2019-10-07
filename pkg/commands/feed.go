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

			utils.ReplyWithError(selectFoodChannel(rtm, channel), feedError, response)
		},
	}
}

func selectFoodChannel(rtm *slack.RTM, channel string) error {
	channels, err := rtm.GetChannels(true)
	if err != nil {
		return err
	}

	options := make([]slack.AttachmentActionOption, len(channels))
	for _, channel := range channels {
		options = append(options, slack.AttachmentActionOption{
			Text:  channel.Name,
			Value: channel.ID,
		})
	}

	attachment := slack.Attachment{
		Text:       selectChannelQuestion,
		Color:      colorPending,
		CallbackID: selectChannelCallbackId,
		Actions: []slack.AttachmentAction{
			Select(
				selectChannelPlaceholder,
				foodChannelKey,
				options,
			),
			CancelButton(),
		},
	}

	rtm.PostMessage(channel, slack.MsgOptionAttachments(attachment))
	return nil
}
