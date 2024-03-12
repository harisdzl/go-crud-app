package zap

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ZapRepo is the Zap logger implementation
type ZapRepo struct {
	zap *zap.Logger
}

// NewZapRepository creates a new Zap logger repository
func NewZapRepository() *ZapRepo {
	// Initialize Zap
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	return &ZapRepo{zapLogger}
}

// Debug logs a debug message
func (z *ZapRepo) Debug(msg string, fields map[string]interface{}) {
	z.zap.Debug(msg, z.convertFields(fields)...)
}

// Info logs an info message
func (z *ZapRepo) Info(msg string, fields map[string]interface{}) {
	z.zap.Info(msg, z.convertFields(fields)...)
}

// Warn logs a warning message
func (z *ZapRepo) Warn(msg string, fields map[string]interface{}) {
	z.zap.Warn(msg, z.convertFields(fields)...)
}

// Error logs an error message
func (z *ZapRepo) Error(msg string, fields map[string]interface{}) {
	z.zap.Error(msg, z.convertFields(fields)...)
}

// Fatal logs a fatal message
func (z *ZapRepo) Fatal(msg string, fields map[string]interface{}) {
	z.zap.Fatal(msg, z.convertFields(fields)...)
}

// convertFields converts fields into Zap-compatible fields
func (z *ZapRepo) convertFields(fields map[string]interface{}) []zap.Field {
	// Convert fields to Zap fields
	var zapFields []zap.Field
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	return zapFields
} 

func (z *ZapRepo) Start(c *gin.Context, info string) trace.Span {return nil}
