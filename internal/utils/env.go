package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func loadEnv() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("Error loading .env file")
	}
	return nil
}

func GetEnv(envVar string) (string, error) {
	if err := loadEnv(); err != nil {
		return "", err
	}

	v := strings.Trim(os.Getenv(envVar), " ")
	if len(v) == 0 {
		return "", fmt.Errorf("no env variable with the name, %s!", envVar)
	}

	return v, nil
}
