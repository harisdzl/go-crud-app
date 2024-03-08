package db

import (
	"github.com/harisquqo/quqo-challenge-1/domain/entity/logger_entity"
	"go.uber.org/zap"
)
func NewLogger() (*logger_entity.Logger, error) {
	var logger *logger_entity.Logger

	// Initialize Zap
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	logger = &logger_entity.Logger{
		Zap: zapLogger,
	}

	

	return logger, nil
}

