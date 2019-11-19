package commands

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/nlopes/slack"
	"github.com/omaressameldin/lunch-roulette/internal/db"
	"github.com/omaressameldin/lunch-roulette/internal/utils"
	"github.com/shomali11/slacker"
)

func organzieLunches(bot *slacker.Slacker) {
	lunchChannels, err := db.GetLunchChannels()
	if err != nil {
		log.Fatal(err)
	}

	bot.Client().ConnectRTM()
	for _, lunch := range lunchChannels {
		OrganizeLunch(bot, lunch.ChannelID)
	}
}

func OrganizeLunch(bot *slacker.Slacker, channelID string) {
	go func() {
		quit := make(chan error, 1)
		updateRound := make(chan error, 1)

		for {
			if err := checkOrganizeReadiness(channelID); err != nil {
				log.Println(err)
				return
			}

			go waitForRound(channelID, quit, updateRound)

			select {
			case err := <-updateRound:
				{

					if addingError := addNextRound(channelID); addingError != nil {
						sendError(channelID, bot, addingError)
						return
					}

					// if error is sent to channel, an update should happen but no setting
					// should be done for this round
					if err != nil {
						log.Println(utils.OrganizeError(channelID, err))
						continue
					}

					selected, selectingError := selectMembers(bot, channelID)
					if selectingError != nil {
						sendError(channelID, bot, selectingError)
						return
					}
					mentionSelectedMembers(bot, channelID, selected)
				}
			case err := <-quit:
				{
					if err != nil {
						log.Println(err)
					}
					close(updateRound)
					close(quit)
					return
				}
			}
		}
	}()
}

func checkOrganizeReadiness(channelID string) error {
	lunchInfo, err := db.GetLunchInfo(channelID)
	if err != nil {
		return utils.OrganizeError(channelID, err)
	}

	if lunchInfo.FrequencyPerMonth == 0 {
		err = fmt.Errorf("freq can't be empty")
		return utils.OrganizeError(channelID, err)
	}

	if lunchInfo.NextRoundDate == nil {
		err = fmt.Errorf("nextRound can't be empty")
		return utils.OrganizeError(channelID, err)
	}

	if lunchInfo.GroupSize == 0 {
		err = fmt.Errorf("groupSize can't be empty")
		return utils.OrganizeError(channelID, err)
	}

	return nil
}

func waitForRound(
	channelID string,
	quit chan<- error,
	updateRound chan<- error,
) {
	lunchInfo, err := db.GetLunchInfo(channelID)
	if err != nil {
		quit <- err
		return
	}
	nextRound := lunchInfo.NextRoundDate
	if nextRound == nil {
		quit <- err //quit if can't get roundDate
	}

	log.Println(organizeLogMessage(channelID, nextRound))

	if err = utils.SleepTill(*nextRound); err != nil {
		updateRound <- fmt.Errorf("error sleeping: %s", err.Error())
		return
	}

	// quit if dates changed
	lunchInfo, err = db.GetLunchInfo(channelID)
	if err != nil {
		quit <- err
		return
	}
	currentRound := lunchInfo.NextRoundDate
	if !currentRound.Equal(*nextRound) {
		log.Printf(dateChangeMessage(channelID))
		quit <- nil
		return
	}

	updateRound <- nil
}

func addNextRound(channelID string) error {
	lunchInfo, err := db.GetLunchInfo(channelID)
	if err != nil {
		return err
	}
	currentRound := lunchInfo.NextRoundDate
	if err != nil {
		return err
	}

	freq := lunchInfo.FrequencyPerMonth

	nextRound := currentRound.AddDate(0, 0, 30/freq+1)
	err = db.AddNextRoundDate(channelID, nextRound)
	if err != nil {
		return err
	}

	return nil
}

func selectMembers(bot *slacker.Slacker, channelID string) ([]string, error) {
	lunchInfo, err := db.GetLunchInfo(channelID)
	if err != nil {
		return nil, err
	}

	groupSize := lunchInfo.GroupSize
	info, err := bot.RTM().GetChannelInfo(channelID)
	if err != nil {
		return nil, err
	}
	members := info.Members

	selectionLimit := int(math.Min(float64(len(members)), float64(groupSize)))
	selected := make([]string, 0, groupSize)
	alreadySelectedMembers, err := db.AllMembers(channelID)
	remainingMembers := utils.Difference(members, alreadySelectedMembers)
	for len(selected) < selectionLimit {
		// use already selected members if no remaining members remain
		if len(remainingMembers) == 0 {
			remainingMembers = alreadySelectedMembers
			if err = db.DeleteAllSelectedMembers(channelID); err != nil {
				return nil, err
			}
		}

		rand.Seed(time.Now().Unix())
		selectedIndex := rand.Intn(len(remainingMembers))
		selected = append(selected, remainingMembers[selectedIndex])
		remainingMembers = utils.Remove(remainingMembers, selectedIndex)
	}

	if err = db.AddMembers(channelID, selected); err != nil {
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
