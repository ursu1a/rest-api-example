package utils

import (
	"fmt"
	"os"
)

func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Check that required environment variables are being set
func CheckEnvs(requiredEnvVars []string) error {
	for _, key := range requiredEnvVars {
		if value := GetEnv(key, ""); value == "" {
			return fmt.Errorf("missing required environment variable: %s", key)
		}
	}

	return nil
}
