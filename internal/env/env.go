package env

import (
	"log"

	"github.com/omaressameldin/lunch-roulette/internal/utils"
)

func ValidateEnvKeys() {
	GetActionPort()
	GetDBName()
}

func GetActionPort() string {
	port, err := utils.GetEnv(actionsPortKey)
	if err != nil {
		log.Fatal(err)
	}

	return port
}

func GetDBName() string {
	dbName, err := utils.GetEnv(dbNameKey)
	if err != nil {
		log.Fatal(err)
	}

	return dbName
}

func GetToken() string {
	token, err := utils.GetEnv(tokenKey)
	if err != nil {
		log.Fatal(err)
	}

	return token
}
