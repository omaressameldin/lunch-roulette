package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func GetEnv(envVar string) (string, error) {
	godotenv.Load()

	v := strings.Trim(os.Getenv(envVar), " ")
	if len(v) == 0 {
		return "", fmt.Errorf("no env variable with the name, %s!", envVar)
	}

	return v, nil
}
