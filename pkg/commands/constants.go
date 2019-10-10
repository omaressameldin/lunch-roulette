package commands

import (
	"fmt"
	"time"
)

// colors
const colorPending = "#f9a41b"
const colorDanger = "#dc3545"
const colorSuccess = "#28a745"

// cancel button
const cancelText = "💣Cancel"
const CancelValue = "cancel-request"

// messages
const addingToDatabase = "⏳ Your response is being saved..."
const timeLayout = "Mon Jan 2 15:04 MST"

// number select
const NumberActionSeparator = "____"

// feed command
const feedError = "Can not start feeding people"
const feedCmd = "feed"
const feedDesc = "Pair group of people together for company paid lunch"

// ---- select channel -----
const selectChannelQuestion = "_*📰 Which channel do you wanna link the bot to?*_"
const SelectChannelBlockID = "select-food-channel"
const selectChannelPlaceholder = "pick a channel"
const selectChannelWarning = "_*⏰Note:*_  `If you choose a channel that is already linked previous data for this channel will be overridden`"
const FoodChannelKey = "food-channel"

// ---- set first round date -----
const FirstRoundStartBlockID = "set-first-round-start"
const firstRoundStartText = "_*📅 When should the first round start?*_"
const RoundTime = "12:20"

// ---- set frequency per month -----
const FerquencyPerMonthBlockID = "set-frequency-per-month"
const frequencyPerMonthText = "_*🕥 how many times per month do you wanna schedule lunches?*_"
const frequencyPerMonyhPlaceholder = "pick frequency"

// ---- set group size ----
const GroupSizeBlockID = "set-group-size"
const groupSizeText = "_* 🙍👱🙍👱 How many people should be paired in one group?*_"

// organize
const membersSelected = "*congratulations🥳* You have been selected for a *free* lunch for this month!"

func organizeLogMessage(channelID string, nextRound *time.Time) string {
	return fmt.Sprintf("organize lunch for %s, on %s", channelID, nextRound)
}

func dateChangeMessage(channelID string) string {
	return fmt.Sprintf(
		"something changed for %s and is being handled by another goroutine",
		channelID,
	)
}
