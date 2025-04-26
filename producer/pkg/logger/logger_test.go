package logger

import (
	"bytes"
	"testing"

	"github.com/elmiringos/indexer/producer/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestCreateEncoder(t *testing.T) {
	tests := []struct {
		stage        string
		expectedType zapcore.Encoder
		expectedLvl  zapcore.Level
	}{
		{
			stage:        "dev",
			expectedType: zapcore.NewConsoleEncoder(zapcore.EncoderConfig{}),
			expectedLvl:  zap.DebugLevel,
		},
		{
			stage:        "prod",
			expectedType: zapcore.NewJSONEncoder(zapcore.EncoderConfig{}),
			expectedLvl:  zap.InfoLevel,
		},
		{
			stage:        "unknown",
			expectedType: zapcore.NewJSONEncoder(zapcore.EncoderConfig{}),
			expectedLvl:  zap.InfoLevel,
		},
	}

	for _, tt := range tests {
		encoder, lvl := createEncoder(tt.stage)

		assert.IsType(t, tt.expectedType, encoder, "encoder type should match stage")
		assert.Equal(t, tt.expectedLvl, lvl, "log level should match stage")
	}
}

func TestNewLogger(t *testing.T) {
	cfg := &config.Config{
		Server: config.Server{
			Stage: "dev",
		},
		Logger: config.Logger{
			File: "testlog",
		},
	}

	logger := New(cfg)
	assert.NotNil(t, logger, "Logger should not be nil")
	assert.IsType(t, &zap.Logger{}, logger, "Logger should be of type *zap.Logger")
}

func TestGetLogger(t *testing.T) {
	cfg := &config.Config{
		Server: config.Server{
			Stage: "dev",
		},
		Logger: config.Logger{
			File: "testlog",
		},
	}

	New(cfg)
	retrievedLogger := GetLogger()

	assert.NotNil(t, retrievedLogger, "Retrieved logger should not be nil")
	assert.IsType(t, &zap.Logger{}, retrievedLogger, "Retrieved logger should be of type *zap.Logger")
}

func TestCreateWriteSyncer(t *testing.T) {
	cfg := &config.Config{
		Logger: config.Logger{
			File: "testlog",
		},
	}

	var buf bytes.Buffer
	lumberjackLoggerFactory = func(filename string) zapcore.WriteSyncer {
		return zapcore.AddSync(&buf)
	}

	defer func() {
		lumberjackLoggerFactory = func(filename string) zapcore.WriteSyncer {
			return zapcore.AddSync(&lumberjack.Logger{
				Filename:   filename,
				MaxSize:    100,
				MaxAge:     1,
				MaxBackups: 1,
				Compress:   true,
			})
		}
	}()

	syncer := createWriteSyncer(cfg)

	assert.NotNil(t, syncer, "WriteSyncer should not be nil")

	_, err := syncer.Write([]byte("test log"))
	assert.NoError(t, err, "Write to syncer should not error")
}
