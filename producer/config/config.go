package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Server  `yaml:"server"`
		HTTP    `yaml:"http"`
		Logger  `yaml:"logger"`
		EthNode `yaml:"eth_node"`
		RMQ
	}

	Server struct {
		Name             string `yaml:"name"`
		Version          string `yaml:"version"`
		Stage            string `yaml:"stage"`
		WorkerCount      int    `yaml:"worker_count"`
		BlockStartNumber string `yaml:"block_start_number"`
		RealTimeMode     bool   `yaml:"real_time_mode"`
		CoreServiceURL   string `env-required:"true" env:"CORE_SERVICE_URL"`
	}

	HTTP struct {
		Port string `env-required:"false" yaml:"port"`
	}

	Logger struct {
		File string `env-required:"false" yaml:"file" env:"LOG_FILE"`
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

func NewConfig(configPath, envPath string) (*Config, error) {
	cfg := &Config{}

	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			return nil, err
		}
	}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, err
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	fmt.Println("Resutl config: ", cfg)

	return cfg, nil
}
