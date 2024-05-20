package config

type Config struct {
	Env  *EnvConfig
	File *FileConfig
}

func NewConfig() *Config {
	return &Config{
		Env:  NewEnvConfig(),
		File: NewFileConfig(),
	}
}
