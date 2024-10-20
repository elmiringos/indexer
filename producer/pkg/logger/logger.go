package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/elmiringos/indexer/producer/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	dirPerm  = 0o755 // Owner: read, write, execute; Group/Others: read, execute
	filePerm = 0o666 // Read and write permissions for owner, group, and others
)

var log *zap.Logger

func GetLogger() *zap.Logger {
	return log
}

func New(cfg *config.Config) *zap.Logger {
	var core zapcore.Core

	encoder, level := createEncoder(cfg.Server.Stage)
	writeSyncer := createWriteSyncer(cfg)
	core = zapcore.NewCore(
		encoder,
		writeSyncer,
		level,
	)

	log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return log
}

func createEncoder(stage string) (zapcore.Encoder, zapcore.Level) {
	switch stage {
	case "dev":
		encoderConfig := zap.NewDevelopmentEncoderConfig()

		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder

		return zapcore.NewConsoleEncoder(encoderConfig), zap.DebugLevel
	case "prod":
		encoderConfig := zap.NewProductionEncoderConfig()

		return zapcore.NewJSONEncoder(encoderConfig), zap.InfoLevel
	default:
		fmt.Print("Error in parsing stage, using default encoder for logger")

		encoderConfig := zap.NewProductionEncoderConfig()

		return zapcore.NewJSONEncoder(encoderConfig), zap.InfoLevel
	}
}

func createWriteSyncer(cfg *config.Config) zapcore.WriteSyncer {
	var writeSyncers []zapcore.WriteSyncer

	if cfg.Logger.File != "" {
		dir := filepath.Dir(cfg.Logger.File)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creating log directory: %v\n", err)
		}

		file, err := os.OpenFile(cfg.Logger.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening log file: %v\n", err)
			// If opening the file fails, just use stdout
			return zapcore.AddSync(os.Stdout)
		}
		writeSyncers = append(writeSyncers, zapcore.AddSync(file))
	}

	// Add console output by default if file output is not set
	writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout))

	// Combine the WriteSyncers
	multiWriteSyncer := zapcore.NewMultiWriteSyncer(writeSyncers...)

	return multiWriteSyncer
}
