package env

import (
	"log"
	"strconv"
	"time"

	"github.com/omaressameldin/lunch-roulette/internal/utils"
)

func ValidateEnvKeys() {
	GetRoundStart()
	GetRoundFrequencyPerMonth()
	GetFirstRoundDate()
	GetFirstRoundDate()
	GetGroupSize()
}

func GetRoundStart() time.Weekday {
	day, err := utils.GetEnv(roundStart)
	if err != nil {
		log.Fatal(err)
	}
	dayInt, err := strconv.Atoi(day)
	if err != nil {
		log.Fatal(err)
	}
	if dayInt < minDayNumber || dayInt > maxDayNumber {
		log.Fatalf("days need to be between %d and %d!", minDayNumber, maxDayNumber)
	}
	return time.Weekday(dayInt)
}

func GetRoundFrequencyPerMonth() int {
	frequency, err := utils.GetEnv(roundFrequency)
	if err != nil {
		log.Fatal(err)
	}
	frequencyNumber, err := strconv.Atoi(frequency)
	if err != nil {
		log.Fatal(err)
	}

	return frequencyNumber
}

func GetFirstRoundDate() time.Time {
	timeStr, err := utils.GetEnv(firstRoundDate)
	if err != nil {
		log.Fatal(err)
	}

	t, err := time.Parse(timeLayout, timeStr)
	if err != nil {
		log.Fatal(err)
	}

	return t
}

func GetGroupSize() int {
	groupSize, err := utils.GetEnv(groupSize)
	if err != nil {
		log.Fatal(err)
	}

	groupSizeInt, err := strconv.Atoi(groupSize)
	if err != nil {
		log.Fatal(err)
	}
	if groupSizeInt < minGroupSize {
		log.Fatalf("Group size must be greater than %d!", minGroupSize)
	}
	return groupSizeInt
}
