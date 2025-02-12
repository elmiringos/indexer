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
		Redis
		JWT
		RMQ `yaml:"rabbitmq"`
	}

	Server struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
		Stage   string `env-required:"true" yaml:"stage"   env:"APP_STAGE"`
		Worker  int    `env-required:"true" yaml:"worker"  env:"APP_WORKER"`
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

	Redis struct {
		URL string `env-required:"true" env:"REDIS_URL"`
	}

	JWT struct {
		SECRET string `env-required:"true" env:"JWT_SECRET"`
	}

	RMQ struct {
		ServerExchange string `env-required:"true" yaml:"rpc_server_exchange" env:"RMQ_RPC_SERVER"`
		ClientExchange string `env-required:"true" yaml:"rpc_client_exchange" env:"RMQ_RPC_CLIENT"`
		URL            string `env-required:"true"                            env:"RMQ_URL"`
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
