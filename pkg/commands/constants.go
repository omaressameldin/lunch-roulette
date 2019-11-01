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
const RoundTime = "14:58"

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

// delete command
const deleteError = "Can not delete schedule"
const deleteCmd = "delete"
const deleteDesc = "deletes a schedule linked to a channel"
const DeleteSuccess = "💣 channel is Successfully unlinked!"

// ---- select channel ----
const selectDeletedQuestion = "_*💣 Which channel do you wanna unlink?*_"
const SelectDeletedBlockID = "select-deleted-channel"
const selectDeletedPlaceholder = "pick a channel"
const selectDeletedWarning = "_*⏰Warning:*_ `this will remove the schedule linked to that channel`"
const deletedKey = "deleted-channel"

// stats command
const statsError = "Can not show stats"
const statsCmd = "stats"
const statsDesc = "get stats of a all scheduled lunches"
const noChannels = "*😕No Channels are linked!*"

// exclude command
const excludeError = "Can not exclude members"
const excludeCmd = "exclude"
const excludeDesc = "excludes a member for one full round"
const ExcludeSuccess = "😑 Member is Successfully Excluded for this round!"
const excludeWarning = "_*⏰Warning:*_ `this will only exclude a member for a full round. If you want to remove a member permanently you should remove him/her from the channel!`"

// ---- select channel ----
const excludeQuestion = "_*💣 Which schedule do you wanna exclude from?*_"
const ExcludeChannelBlockID = "select-exclude-channel"
const excludeChannelPlaceholder = "pick a channel"
const excludedChannelKey = "exclude-channel"

// ---- select Member ----
const memberExcludeQuestion = "_*💣 Which Member do you wanna exclude from channel?*_"
const ExcludeMemberBlockID = "select-exclude-member"
const excludeMemberPlaceholder = "pick a member"
