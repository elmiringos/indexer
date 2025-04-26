package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultConfig_FileNotFound(t *testing.T) {
	// Rename/move config file temporarily so it doesn't exist
	os.Rename("./config/config.yml", "./config/config.yml.bak")
	defer os.Rename("./config/config.yml.bak", "./config/config.yml")

	cfg, err := NewDefaultConfig()

	assert.Nil(t, cfg)
	assert.Error(t, err)
}

func TestNewConfig_Success(t *testing.T) {
	// Arrange: Create temporary .env and config files
	_ = os.WriteFile("./test.env", []byte(`CORE_SERVICE_URL=http://localhost
RMQ_URL=amqp://localhost
ETH_HTTP_NODE_RPC=http://localhost
ETH_WS_NODE_RPC=ws://localhost
`), 0644)

	_ = os.WriteFile("./test_config.yml", []byte(`
server:
  name: "TestService"
  version: "1.0"
  stage: "dev"
  worker_count: 5
  block_start_number: "1000"
  real_time_mode: true
  core_service_url: "http://core"
http:
  port: "8080"
logger:
  file: "app.log"
eth_node:
  network_type: "mainnet"
  trace_enabled: true
`), 0644)

	defer os.Remove("./test.env")
	defer os.Remove("./test_config.yml")

	// Act
	cfg, err := NewConfig("./test_config.yml", "./test.env")

	// Assert
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "TestService", cfg.Server.Name)
	assert.Equal(t, "8080", cfg.HTTP.Port)
	assert.Equal(t, "mainnet", cfg.EthNode.Network)
}

func TestNewConfig_EnvMissing(t *testing.T) {
	// Arrange
	_ = os.WriteFile("./test_config.yml", []byte(`
server:
  name: "Service"
  version: "1.0"
  stage: "dev"
  worker_count: 3
  block_start_number: "1"
  real_time_mode: false
  core_service_url: "http://service"
http:
  port: "9090"
logger:
  file: "service.log"
eth_node:
  network_type: "rinkeby"
  trace_enabled: false
`), 0644)

	defer os.Remove("./test_config.yml")

	cfg, err := NewConfig("./test_config.yml", "./non_existent.env")

	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "9090", cfg.HTTP.Port)
	assert.Equal(t, "rinkeby", cfg.EthNode.Network)
}
