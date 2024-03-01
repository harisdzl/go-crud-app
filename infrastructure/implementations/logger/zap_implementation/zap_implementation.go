package zap

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	"go.uber.org/zap"
)

// ZapRepo is the Zap logger implementation
type ZapRepo struct {
	p *base.Persistence
	c *gin.Context
}

// NewZapRepository creates a new Zap logger repository
func NewZapRepository(p *base.Persistence, c *gin.Context) *ZapRepo {
	return &ZapRepo{p, c}
}

// Debug logs a debug message
func (z *ZapRepo) Debug(msg string, fields map[string]interface{}) {
	z.p.Logger.Zap.Debug(msg, z.convertFields(fields)...)
}

// Info logs an info message
func (z *ZapRepo) Info(msg string, fields map[string]interface{}) {
	z.p.Logger.Zap.Info(msg, z.convertFields(fields)...)
}

// Warn logs a warning message
func (z *ZapRepo) Warn(msg string, fields map[string]interface{}) {
	z.p.Logger.Zap.Warn(msg, z.convertFields(fields)...)
}

// Error logs an error message
func (z *ZapRepo) Error(msg string, fields map[string]interface{}) {
	z.p.Logger.Zap.Error(msg, z.convertFields(fields)...)
}

// Fatal logs a fatal message
func (z *ZapRepo) Fatal(msg string, fields map[string]interface{}) {
	z.p.Logger.Zap.Fatal(msg, z.convertFields(fields)...)
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
