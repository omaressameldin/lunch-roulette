package commands

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/nlopes/slack"
	"github.com/omaressameldin/lunch-roulette/internal/db"
	"github.com/omaressameldin/lunch-roulette/internal/utils"
	"github.com/shomali11/slacker"
)

func organzieLunches(d *db.DB, bot *slacker.Slacker) {
	lunchChannels, err := d.GetBotChannels()
	if err != nil {
		log.Fatal(err)
	}

	bot.Client().ConnectRTM()
	for _, channelID := range lunchChannels {
		OrganizeLunch(bot, d, channelID)
	}
}

func OrganizeLunch(bot *slacker.Slacker, d *db.DB, channelID string) {
	go func() {
		for {
			if !canSchedule(d, channelID) {
				return
			}
			if err := waitForRound(d, channelID); err != nil {
				log.Println(err)
				return
			}
			selected, err := selectMembers(bot, d, channelID)
			if err != nil {
				sendError(channelID, bot, err)
			}
			log.Println(selected)
			mentionSelectedMembers(bot, channelID, selected)
		}
	}()
}

func canSchedule(d *db.DB, channelID string) bool {
	freq, err := d.GetFrequencyPerMonth(channelID)
	if err != nil || freq == nil {
		return false
	}

	nextRound, err := d.GetNextRoundDate(channelID)
	if err != nil || nextRound == nil {
		return false
	}

	groupSize, err := d.GetGroupSize(channelID)
	if err != nil || groupSize == nil {
		return false
	}

	return true
}

func waitForRound(d *db.DB, channelID string) error {
	freq, err := d.GetFrequencyPerMonth(channelID)
	if err != nil {
		return err
	}

	nextRound, err := d.GetNextRoundDate(channelID)
	if err != nil {
		return err
	}

	// add freq weeks to next round
	err = utils.SleepTill(*nextRound)
	currentRound, err := d.GetNextRoundDate(channelID)
	if err != nil {
		return err
	}
	// quit if dates changed
	if !currentRound.Equal(*nextRound) {
		return errors.New("something changed and should have been handled by another goroutine")
	}
	newRound := nextRound.AddDate(0, 0, 7*(*freq))
	d.AddNextRoundDate(channelID, newRound)
	if err != nil {
		return err
	}

	return nil
}

func selectMembers(bot *slacker.Slacker, d *db.DB, channelID string) ([]string, error) {
	freq, err := d.GetFrequencyPerMonth(channelID)
	if err != nil {
		return nil, err
	}
	info, err := bot.RTM().GetChannelInfo(channelID)
	if err != nil {
		return nil, err
	}
	members := info.Members
	selected := make([]string, 0, *freq)
	alreadySelectedMembers, err := d.AllMembers(channelID)
	remainingMembers := utils.Difference(members, alreadySelectedMembers)

	for len(selected) < *freq && len(remainingMembers) > 0 {
		rand.Seed(time.Now().Unix())
		selectedIndex := rand.Intn(len(remainingMembers))
		selected = append(selected, remainingMembers[selectedIndex])
		remainingMembers = utils.Remove(remainingMembers, selectedIndex)

		// use already selected members if no remaining members remain
		if len(remainingMembers) == 0 {
			remainingMembers = alreadySelectedMembers
			d.DeleteAllSelectedMembers(channelID)
		}
	}

	err = d.AddMembers(channelID, selected)
	if err != nil {
		return nil, err
	}

	return selected, nil
}

func mentionSelectedMembers(
	bot *slacker.Slacker,
	channelID string,
	selectedMembers []string,
) {
	mentions := bytes.Buffer{}
	for _, member := range selectedMembers {
		mentions.WriteString(fmt.Sprintf("<@%s>, ", member))
	}
	bot.RTM().PostMessage(channelID, slack.MsgOptionAttachments(
		SuccessMessage(fmt.Sprintf("%s%s", mentions.String(), membersSelected)),
	))
}
