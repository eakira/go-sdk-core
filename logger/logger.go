package logger

import (
	"os"
	"strings"

	"github.com/eakira/go-sdk-core/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log            *zap.Logger
	LOG_OUTPUT     = env.GetEnv("LOG_OUTPUT")
	LOG_PATH       = env.GetEnv("LOG_INFO_PATH")
	LOG_ERROR_PATH = env.GetEnv("LOG_ERROR_PATH")
	LOG_LEVEL      = env.GetEnv("LOG_LEVEL")
	AMBIENTE_DEV   = env.GetAmbienteDev()
)

func init() {
	// Abra ou crie o arquivo de log de erro
	errFile, err := os.OpenFile(LOG_ERROR_PATH, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil && AMBIENTE_DEV {
		panic(err)
	}
	defer errFile.Close() // Feche o arquivo depois de usar

	logConfig := zap.Config{
		OutputPaths:      []string{"stdout", LOG_PATH},
		ErrorOutputPaths: []string{LOG_ERROR_PATH},
		Level:            zap.NewAtomicLevelAt(getLevelLogs()),
		Encoding:         "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "message",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	log, _ = logConfig.Build()
}

func Info(message string, journey string, tags ...zap.Field) {
	tags = append(tags, zap.String("journey", journey))

	log.Info(message, tags...)
	errSync := log.Sync()
	if errSync != nil {
		return
	}
}

func Error(message string, err error, journey string, tags ...zap.Field) {
	tags = append(tags, zap.String("journey", journey))
	tags = append(tags, zap.NamedError("error", err))

	log.Error(message, tags...)
	log.Sync()
	if AMBIENTE_DEV {
		panic(err)
	}
}

func getLevelLogs() zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(LOG_LEVEL)) {
	case "info":
		return zapcore.InfoLevel
	case "error":
		return zapcore.ErrorLevel
	case "debug":
		return zapcore.DebugLevel
	default:
		return zapcore.InfoLevel
	}
}

// InitTestLogger inicializa um logger simples para uso nos testes
func InitTestLogger() {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{}

	var err error
	log, err = cfg.Build()
	if err != nil {
		panic(err)
	}
}

func GetTestLogger() *zap.Logger {
	if log == nil {
		InitTestLogger()
	}
	return log
}
