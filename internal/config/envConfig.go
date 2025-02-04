package config

import (
	"miracummapper/internal/utilities"
	"reflect"
)

type EnvConfig struct {
	Port             string
	DBUser           string
	DBName           string
	DBHost           string
	DBPort           string
	DBPassword       string
	KeycloakUrl      string
	KeycloakRealm    string
	KeycloakClientId string
}

var DefaultConfig = EnvConfig{
	Port:             "8080",
	DBUser:           "miracum_user",
	DBName:           "miracum_db",
	DBHost:           "localhost",
	DBPort:           "5432",
	DBPassword:       "miracum_password",
	KeycloakUrl:      "http://localhost:8081",
	KeycloakRealm:    "master",
	KeycloakClientId: "miracum-mapper",
}

var EnvKeys = EnvConfig{
	Port:             "PORT",
	DBUser:           "DB_USER",
	DBName:           "DB_NAME",
	DBHost:           "DB_HOST",
	DBPassword:       "DB_PASSWORD",
	KeycloakUrl:      "KEYCLOAK_URL",
	KeycloakRealm:    "KEYCLOAK_REALM",
	KeycloakClientId: "KEYCLOAK_CLIENT_ID",
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
