package env

import (
	"log"
	"strings"

	"github.com/omaressameldin/lunch-roulette/internal/utils"
)

func ValidateEnvKeys() {
	GetActionPort()
	GetDatabaseUrl()
}

func GetActionPort() string {
	port, err := utils.GetEnv(actionsPortKey)
	if err != nil {
		log.Fatal(err)
	}

	return port
}

func DoesHaveCredentials() bool {
	_, err := utils.GetEnv(driveCredentialsKey)

	return err == nil
}

func GetDriveCredentials() string {
	driveCredentials, err := utils.GetEnv(driveCredentialsKey)
	if err != nil {
		log.Fatal(err)
	}

	return driveCredentials
}

func GetToken() string {
	token, err := utils.GetEnv(tokenKey)
	if err != nil {
		log.Fatal(err)
	}

	return token
}

func GetAuthUsers() []string {
	authUsers, err := utils.GetEnv(authUsersKey)
	if err != nil {
		return []string{}
	}

	return strings.Split(authUsers, ",")
}

func GetDatabaseUrl() string {
	databaseUrl, err := utils.GetEnv(databaseUrlKey)
	if err != nil {
		log.Fatal(err)
	}

	return databaseUrl
}
