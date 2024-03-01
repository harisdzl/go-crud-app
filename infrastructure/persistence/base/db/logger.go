package db

import (
	"github.com/harisquqo/quqo-challenge-1/domain/entity/logger_entity"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)
func NewLogger() (*logger_entity.Logger, error) {
	var logger *logger_entity.Logger

	// Initialize Zap
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	// Initialize Honeycomb Tracer
	otelTracer := otel.Tracer("")
	// ctx, span := otelTracer.Start("")

	logger = &logger_entity.Logger{
		Honeycomb: otelTracer,
		Zap: zapLogger,
	}

	

	return logger, nil
}

