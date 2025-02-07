package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/elmiringos/indexer/producer/config"
	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

	filename := cfg.Logger.File + time.Now().Format("2006-01-02_15-04-05") + ".log"

	writeSyncers = append(writeSyncers, zapcore.AddSync(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    100,
		MaxAge:     1,
		MaxBackups: 1,
		Compress:   true,
	}))

	// Add console output by default if file output is not set
	writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout))

	// Combine the WriteSyncers
	multiWriteSyncer := zapcore.NewMultiWriteSyncer(writeSyncers...)

	return multiWriteSyncer
}
