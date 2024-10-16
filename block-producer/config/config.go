package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		App    `yaml:"app"`
		HTTP   `yaml:"http"`
		Logger `yaml:"logger"`
		Redis
		EthNode `yaml:"eth_node"`
		RMQ
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
		Stage   string `env-required:"true" yaml:"stage"   env:"APP_STAGE"`
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
		ServerExchange string `env-required:"true" yaml:"rpc_server_exchange" env:"RMQ_RPC_SERVER"`
		ClientExchange string `env-required:"true" yaml:"rpc_client_exchange" env:"RMQ_RPC_CLIENT"`
		URL            string `env-required:"true"                            env:"RMQ_URL"`
	}

	EthNode struct {
		HttpURL string `env-required:"true" env:"ETH_HTTP_NODE_RPC"`
		WsURL   string `env-required:"true" env:"ETH_WS_NODE_RPC"`
		ApiKey  string `env-required:"false" env:"ETHN_RPC_KEY"`
		Network string `yaml:"network_type"`
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
