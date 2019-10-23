package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/nlopes/slack"
	"github.com/omaressameldin/lunch-roulette/internal/db"
	"github.com/omaressameldin/lunch-roulette/internal/utils"
	"github.com/shomali11/slacker"
)

func AddStatsCmd(d *db.DB, bot *slacker.Slacker) {
	bot.Command(statsCmd, stats(d))
}

func stats(d *db.DB) *slacker.CommandDefinition {
	return &slacker.CommandDefinition{
		Description: statsDesc,
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			channel := request.Event().Channel
			rtm := response.RTM()

			utils.ReplyWithError(getStats(d, channel, rtm), statsError, response)
		},
	}
}

func getStats(d *db.DB, channel string, rtm *slack.RTM) error {
	channels, err := d.GetBotChannels()
	if err != nil {
		return err
	}

	if len(channels) == 0 {
		rtm.PostMessage(channel, slack.MsgOptionBlocks(
			slack.NewContextBlock("", slack.NewTextBlockObject(
				"mrkdwn",
				noChannels,
				false,
				false,
			)),
		))
	}

	message := make([]slack.Block, 0, len(channels))
	for _, channelID := range channels {
		freq, err := d.GetFrequencyPerMonth(channelID)
		if err != nil {
			return err
		}

		nextRound, err := d.GetNextRoundDate(channelID)
		if err != nil {
			return err
		}

		groupSize, err := d.GetGroupSize(channelID)
		if err != nil {
			return err
		}

		channelInfo, err := rtm.GetChannelInfo(channelID)
		if err != nil {
			return err
		}

		members, err := d.AllMembers(channelID)
		if err != nil {
			return err
		}
		memberNames := make([]string, 0, len(members))
		for _, memberID := range members {
			memberInfo, err := rtm.GetUserInfo(memberID)
			if err != nil {
				return err
			}
			memberNames = append(memberNames, memberInfo.Name)
		}

		message = append(message,
			createMessage(
				channelInfo.Name,
				*freq,
				nextRound,
				*groupSize,
				memberNames,
			),
			slack.NewDividerBlock(),
		)
	}

	rtm.PostMessage(channel, slack.MsgOptionBlocks(
		message...,
	))
	return nil
}

func createMessage(
	channelName string,
	freq int,
	nextRoundDate *time.Time,
	groupSize int,
	memberNames []string,
) slack.Block {
	return slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", showChannel(channelName), false, false),
		[]*slack.TextBlockObject{
			slack.NewTextBlockObject("mrkdwn", showNextRoundDate(nextRoundDate), false, false),
			slack.NewTextBlockObject("mrkdwn", showGroupSize(groupSize), false, false),
			slack.NewTextBlockObject("mrkdwn", showFreq(freq), false, false),
			slack.NewTextBlockObject("mrkdwn", showMembers(memberNames), false, false),
		},
		nil,
	)
}

func showChannel(channelName string) string {
	return fmt.Sprintf("*Lunch Info for channel:* %s", channelName)
}

func showFreq(freq int) string {
	return fmt.Sprintf("*🕥Frequency Per Month:* %d", freq)
}

func showNextRoundDate(nextRound *time.Time) string {
	formattedTime := "___"
	if nextRound != nil {
		formattedTime = nextRound.Format(timeLayout)
	}
	return fmt.Sprintf("*📆Next Round Date:*\n %s", formattedTime)
}

func showGroupSize(groupSize int) string {
	return fmt.Sprintf("*🙍👱🙍👱Group Size:* %d", groupSize)
}

func showMembers(memberNames []string) string {
	names := "___"
	if len(memberNames) > 0 {
		names = strings.Join(memberNames, "\n")
	}

	return fmt.Sprintf("*🙍👱🙍👱 Already Selected:*\n %s", names)
}
