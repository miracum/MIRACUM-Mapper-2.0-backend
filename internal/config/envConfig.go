package config

import (
	"miracummapper/internal/utilities"
	"reflect"
)

type EnvConfig struct {
	Port       string
	DBUser     string
	DBName     string
	DBHost     string
	DBPort     string
	DBPassword string
}

var DefaultConfig = EnvConfig{
	Port:       "8080",
	DBUser:     "postgres",
	DBName:     "postgres",
	DBHost:     "localhost",
	DBPort:     "5432",
	DBPassword: "postgres",
}

var EnvKeys = EnvConfig{
	DBUser:     "DB_USER",
	DBName:     "DB_NAME",
	DBHost:     "DB_HOST",
	DBPassword: "DB_PASSWORD",
}

func NewEnvConfig() *EnvConfig {
	cfg := &EnvConfig{}
	keys := reflect.ValueOf(EnvKeys)
	defaults := reflect.ValueOf(DefaultConfig)
	cfgValue := reflect.ValueOf(cfg).Elem()

	for i := 0; i < keys.NumField(); i++ {
		key := keys.Field(i).String()
		defaultValue := defaults.Field(i).String()
		cfgValue.Field(i).SetString(utilities.GetEnv(key, defaultValue))
	}

	return cfg
}
