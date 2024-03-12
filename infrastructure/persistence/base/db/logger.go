package db

import (
	"os"

	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/logger"
)
func NewLogger() (logger.LoggerRepo, error) {

	logChannels := []string{os.Getenv("LOG_CHANNEL_ZAP"), os.Getenv("LOG_CHANNEL_HONEYCOMB")}
	logger := logger.NewLoggerRepository(logChannels)

	return logger, nil
}

