package config

import (
	"fmt"
	"log"
	"os"

	"dario.cat/mergo"

	"gopkg.in/yaml.v3"
)

type FileConfig struct {
	DatabaseConfig struct {
		// Defines how often to retry to connect to the database
		Retry int `yaml:"retry"`

		// Sleep is the amount of time to wait between retries
		Sleep int `yaml:"sleep"`
	} `yaml:"database"`

	KeyCloakConfig struct {
		// Defines how often to retry to connect to keycloak
		Retry int `yaml:"retry"`

		// Sleep is the amount of time to wait between retries
		Sleep int `yaml:"sleep"`
	} `yaml:"keycloak"`

	// } `yaml:"rate_limit"`
	CorsConfig struct {
		// AllowedOrigins is a list of origins a cross-domain request can be executed from
		// If the special "*" value is present in the list, all origins will be allowed.
		// An origin may contain a wildcard (*) to replace 0 or more characters (i.e.: http://*.domain.com).
		AllowedOrigins []string `yaml:"allowed_origins"`
	} `yaml:"cors"`

	// Debug is a flag to enable debug mode
	Debug bool `yaml:"debug"`
}

func NewFileConfig() *FileConfig {
	fileConfig, err := getFileConfig("config.yaml")
	if err != nil {
		log.Printf("failed to get file config: %v.", err)
		panic(err)
	}
	return fileConfig
}

func getFileConfig(configPath string) (*FileConfig, error) {
	// Load default config
	defaultConfig, err := loadConfig("default-config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to load default config: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// File doesn't exist, create it with default values
		defaultData, err := os.ReadFile("default-config.yaml")
		if err != nil {
			return nil, err
		}
		err = os.WriteFile(configPath, defaultData, 0644)
		if err != nil {
			return nil, err
		}
	} else {
		// File exists, load it
		config, err := loadConfig(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %v", err)
		}

		// Merge loaded config into default config
		if err := mergo.Merge(defaultConfig, config, mergo.WithOverride); err != nil {
			return nil, fmt.Errorf("failed to merge configs: %v", err)
		}
	}

	return defaultConfig, nil
}

func loadConfig(path string) (*FileConfig, error) {
	fileConfig := &FileConfig{}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist", path)
	}

	// Open file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&fileConfig); err != nil {
		return nil, err
	}

	return fileConfig, nil
}
