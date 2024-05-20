package utilities

import (
	"log"
	"os"
)

func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
		log.Printf("Environment variable '%s' is not set. Using default value: '%s'", key, fallback)
	}
	return value
}

func GetOrDefault[T any](value *T, defaultValue T) T {
	if value != nil {
		return *value
	}
	return defaultValue
}
