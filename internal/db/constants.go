package db

import "errors"

//buckets
var bucketError = errors.New("This channel is not linked to any schedules!")

const keyNextRound = "nextRound"
const keyBotChannel = "botChannel"
const keyFrequencyPerMonth = "frequencyPerMonth"
const keyGroupSize = "groupSize"
const keyMembers = "members"

// constants
const timeLayout = "02.01.2006 15:04"
