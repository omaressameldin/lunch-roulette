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

func AddStatsCmd(bot *slacker.Slacker) {
	bot.Command(statsCmd, stats())
}

func stats() *slacker.CommandDefinition {
	return &slacker.CommandDefinition{
		Description: statsDesc,
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			channel := request.Event().Channel
			rtm := response.RTM()

			utils.ReplyWithError(getStats(channel, rtm), statsError, response)
		},
	}
}

func getStats(channel string, rtm *slack.RTM) error {
	lunches, err := db.GetLunchChannels()
	if err != nil {
		return err
	}

	if len(lunches) == 0 {
		rtm.PostMessage(channel, slack.MsgOptionBlocks(
			slack.NewContextBlock("", slack.NewTextBlockObject(
				"mrkdwn",
				noChannels,
				false,
				false,
			)),
		))
	}

	message := make([]slack.Block, 0, len(lunches))
	for _, lunch := range lunches {
		freq := lunch.FrequencyPerMonth
		nextRound := lunch.NextRoundDate
		groupSize := lunch.GroupSize

		channelInfo, err := rtm.GetChannelInfo(lunch.ChannelID)
		if err != nil {
			return err
		}

		members, err := db.AllMembers(lunch.ChannelID)
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
				freq,
				nextRound,
				groupSize,
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
