package env

import (
	"log"
	"strconv"
	"time"

	"github.com/omaressameldin/lunch-roulette/internal/utils"
)

func ValidateEnvKeys() {
	GetActionPort()
	GetRoundFrequencyPerMonth()
	GetFirstRoundDate()
	GetFirstRoundDate()
	GetGroupSize()
	GetDBName()
	GetDBBucket()
}

func GetActionPort() string {
	port, err := utils.GetEnv(actionsPortKey)
	if err != nil {
		log.Fatal(err)
	}

	return port
}

func GetRoundFrequencyPerMonth() int {
	frequency, err := utils.GetEnv(roundFrequencyKey)
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
	timeStr, err := utils.GetEnv(firstRoundDateKey)
	if err != nil {
		log.Fatal(err)
	}

	t, err := time.Parse(TimeLayout, timeStr)
	if err != nil {
		log.Fatal(err)
	}

	return t
}

func GetGroupSize() int {
	groupSize, err := utils.GetEnv(groupSizeKey)
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

func GetDBName() string {
	dbName, err := utils.GetEnv(dbNameKey)
	if err != nil {
		log.Fatal(err)
	}

	return dbName
}

func GetDBBucket() string {
	dbBucket, err := utils.GetEnv(dbBucketKey)
	if err != nil {
		log.Fatal(err)
	}

	return dbBucket
}

func GetToken() string {
	token, err := utils.GetEnv(tokenKey)
	if err != nil {
		log.Fatal(err)
	}

	return token
}
