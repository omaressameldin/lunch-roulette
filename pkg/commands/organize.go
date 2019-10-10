package commands

import (
	"bytes"
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
		quit := make(chan error, 1)
		updateRound := make(chan error, 1)

		for {
			if err := checkOrganizeReadiness(d, channelID); err != nil {
				sendError(channelID, bot, err)
			}

			go waitForRound(d, channelID, quit, updateRound)

			select {
			case err := <-updateRound:
				{

					if addingError := addNextRound(d, channelID); addingError != nil {
						sendError(channelID, bot, addingError)
						return
					}

					// if error is sent to channel, an update should happen but no setting
					// should be done for this round
					if err != nil {
						log.Println(utils.OrganizeError(channelID, err))
						continue
					}

					selected, selectingError := selectMembers(bot, d, channelID)
					if selectingError != nil {
						sendError(channelID, bot, selectingError)
						return
					}
					mentionSelectedMembers(bot, channelID, selected)
				}
			case err := <-quit:
				{
					if err != nil {
						sendError(channelID, bot, err)
					}
					close(updateRound)
					close(quit)
					return
				}
			}
		}
	}()
}

func checkOrganizeReadiness(d *db.DB, channelID string) error {
	if freq, err := d.GetFrequencyPerMonth(channelID); err != nil || freq == nil {
		return utils.OrganizeError(channelID, err)
	}

	if nextRound, err := d.GetNextRoundDate(channelID); err != nil || nextRound == nil {
		return utils.OrganizeError(channelID, err)
	}

	if groupSize, err := d.GetGroupSize(channelID); err != nil || groupSize == nil {
		return utils.OrganizeError(channelID, err)
	}

	return nil
}

func waitForRound(
	d *db.DB, channelID string,
	quit chan<- error,
	updateRound chan<- error,
) {
	nextRound, err := d.GetNextRoundDate(channelID)
	if err != nil {
		quit <- err //quit if can't get roundDate
	}

	log.Println(organizeLogMessage(channelID, nextRound))

	if err = utils.SleepTill(*nextRound); err != nil {
		updateRound <- fmt.Errorf("error sleeping: %s", err.Error())
		return
	}

	// quit if dates changed
	currentRound, err := d.GetNextRoundDate(channelID)
	if err != nil {
		quit <- err
	}
	if !currentRound.Equal(*nextRound) {
		log.Printf(dateChangeMessage(channelID))
		quit <- nil
	}

	updateRound <- nil
}

func addNextRound(d *db.DB, channelID string) error {
	currentRound, err := d.GetNextRoundDate(channelID)
	if err != nil {
		return err
	}

	freq, err := d.GetFrequencyPerMonth(channelID)
	if err != nil {
		return err
	}

	nextRound := currentRound.AddDate(0, 0, 30/(*freq)+1)
	err = d.AddNextRoundDate(channelID, nextRound)
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
