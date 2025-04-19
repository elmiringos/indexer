package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Server `yaml:"server"`
		HTTP   `yaml:"http"`
		Logger `yaml:"logger"`
		PG
	}

	Server struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
		Stage   string `env-required:"true" yaml:"stage"   env:"APP_STAGE"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Logger struct {
		File string `env-required:"false" yaml:"file" env:"LOG_FILE"`
	}

	PG struct {
		URL string `env-required:"true" env:"PG_URL"`
	}
)

func NewDefaultConfig() (*Config, error) {
	configPath := "./config/config.yml"
	envPath := "./.env"

	return NewConfig(configPath, envPath)
}

// NewConfig returns app config.
func NewConfig(configPath, envPath string) (*Config, error) {
	cfg := &Config{}

	// Check if the .env file exists at the provided path
	if _, err := os.Stat(envPath); err == nil {
		err := godotenv.Load(envPath)
		if err != nil {
			return nil, err
		}
	}

	err := godotenv.Load(envPath)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
