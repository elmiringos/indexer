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
		Redis
		EthNode `yaml:"eth_node"`
		RMQ
	}

	Server struct {
		Name             string `yaml:"name"`
		Version          string `yaml:"version"`
		Stage            string `yaml:"stage"`
		WorkerCount      int    `yaml:"worker_count"`
		BlockStartNumber string `yaml:"block_start_number"`
		CoreServiceUrl   string `yaml:"core_service_url"`
	}

	HTTP struct {
		Port string `env-required:"false" yaml:"port"`
	}

	Logger struct {
		File string `env-required:"false" yaml:"file" env:"LOG_FILE"`
	}

	Redis struct {
		URL string `env-required:"true" env:"REDIS_URL"`
	}

	RMQ struct {
		URL string `env-required:"true" env:"RMQ_URL"`
	}

	EthNode struct {
		HttpURL string `env-required:"true" env:"ETH_HTTP_NODE_RPC"`
		WsURL   string `env-required:"true" env:"ETH_WS_NODE_RPC"`
		ApiKey  string `env:"ETH_RPC_KEY"`
		Network string `yaml:"network_type"`
		Trace   bool   `yaml:"trace_enabled"`
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
