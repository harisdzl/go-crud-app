package db

import (
	"os"
	"strings"

	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/logger"
)
func NewLogger() (logger.LoggerRepo, error) {

	logChannels := os.Getenv("LOGGER_CHANNELS")
	logChannelsList := strings.Split(logChannels, ",")
	logger := logger.NewLoggerRepository(logChannelsList)

	return logger, nil
}

