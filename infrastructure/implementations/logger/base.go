package logger

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/logger_repository"
	honeycomb "github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/logger/honeycomb_implementation"
	zap "github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/logger/zap_implementation"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	"go.opentelemetry.io/otel/trace"
)

// LoggerRepo is a logger repository that can use multiple channels
type LoggerRepo struct {
	loggers []logger_repository.LoggerRepository
	c *gin.Context
	p *base.Persistence
	Span trace.Span
}

const (
	Zap = "Zap"
	Honeycomb = "Honeycomb"
)

// NewLoggerRepository creates a new logger repository based on the specified channels
func NewLoggerRepository(channels []string, p *base.Persistence, c *gin.Context, info string) (LoggerRepo, error) {
	var loggers []logger_repository.LoggerRepository
	var honeycombRepo *honeycomb.HoneycombRepo
	for _, channel := range channels {
		switch channel {
		case Zap:
			loggers = append(loggers, zap.NewZapRepository(p, c))
		case Honeycomb:
			honeycombRepo = honeycomb.NewHoneycombRepository(p, c, info)
			loggers = append(loggers, honeycombRepo)
		default:
			// You might want to log or handle unsupported channels
			continue
		}
	}

	// if len(loggers) == 0 {
	// 	return nil, errors.New("no supported logger type found in the provided channels")
	// }

	return LoggerRepo{loggers: loggers, c: c, Span: honeycombRepo.Span}, nil
}

// Debug logs a debug message
func (l *LoggerRepo) Debug(msg string, fields map[string]interface{}) {
	for _, logger := range l.loggers {
		logger.Debug(msg, fields)
	}
}

// Info logs an info message
func (l *LoggerRepo) Info(msg string, fields map[string]interface{}) {
	for _, logger := range l.loggers {
		logger.Info(msg, fields)
	}
}

// Warn logs a warning message
func (l *LoggerRepo) Warn(msg string, fields map[string]interface{}) {
	for _, logger := range l.loggers {
		logger.Warn(msg, fields)
	}
}

// Error logs an error message
func (l *LoggerRepo) Error(msg string, fields map[string]interface{}) {
	for _, logger := range l.loggers {
		logger.Error(msg, fields)
	}
}

// Fatal logs a fatal message
func (l *LoggerRepo) Fatal(msg string, fields map[string]interface{}) {
	for _, logger := range l.loggers {
		logger.Fatal(msg, fields)
	}
}

// End function

func (l *LoggerRepo) End() {
	ctxValue, exists := l.c.Get("otel_context")
    var ctx context.Context
    if exists {
        ctx, _ = ctxValue.(context.Context)
    } else {
        ctx = l.c.Request.Context()
    }
	log.Println(ctx)
	span := trace.SpanFromContext(ctx)
	span.End()
}