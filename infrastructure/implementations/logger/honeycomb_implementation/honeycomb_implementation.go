package honeycomb

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// HoneycombRepo is the Honeycomb logger implementation
type HoneycombRepo struct {
	p *base.Persistence
	c *gin.Context
	Span trace.Span
}

func NewHoneycombRepository(p *base.Persistence, c *gin.Context, info string) *HoneycombRepo {
    // Retrieve otel_context from Gin context if it exists
    ctxValue, exists := c.Get("otel_context")
    // Implement info logging with Honeycomb using the retrieved context
    var ctx context.Context
    if exists {
        ctx, _ = ctxValue.(context.Context)
    } else {
        ctx = c.Request.Context()
    }
    tracer := p.Logger.Honeycomb
    _, span := tracer.Start(ctx, info)
    return &HoneycombRepo{p, c, span}
}
// Debug logs a debug message
func (h *HoneycombRepo) Debug(msg string, fields map[string]interface{}) {
	jsonData, jsonDataErr := json.Marshal(fields)

	if jsonDataErr != nil {
		return
	}

	h.Span.AddEvent(msg, trace.WithAttributes(
		attribute.String("level", "Debug"),
		attribute.String("data", string(jsonData))))
	
}

// Info logs an info message
func (h *HoneycombRepo) Info(msg string, fields map[string]interface{}) {

	jsonData, jsonDataErr := json.Marshal(fields)

	if jsonDataErr != nil {
		return
	}

	h.Span.AddEvent(msg, trace.WithAttributes(
		attribute.String("level", "Info"),
		attribute.String("data", string(jsonData))))
	
}

func (h *HoneycombRepo) Error(msg string, fields map[string]interface{}) {
    jsonData, jsonDataErr := json.Marshal(fields)
    if jsonDataErr != nil {
        return
    }

    h.Span.RecordError(errors.New(msg), trace.WithAttributes(
		attribute.String("level", "Error"),
        attribute.String("data", string(jsonData))))
}


// Warn logs a warning message
func (h *HoneycombRepo) Warn(msg string, fields map[string]interface{}) {
    // Implement warn logging with Honeycomb
    // This function can be implemented similar to the Info function

    jsonData, jsonDataErr := json.Marshal(fields)
    if jsonDataErr != nil {
        return
    }

    h.Span.AddEvent(msg, trace.WithAttributes(
		attribute.String("level", "Warn"),
        attribute.String("data", string(jsonData))))
}

// Fatal logs a fatal message
func (h *HoneycombRepo) Fatal(msg string, fields map[string]interface{}) {
    // Implement fatal logging with Honeycomb
    // This function can be implemented similar to the Info function

    jsonData, jsonDataErr := json.Marshal(fields)
    if jsonDataErr != nil {
        return
    }

    h.Span.AddEvent(msg, trace.WithAttributes(
		attribute.String("level", "Fatal"),
        attribute.String("data", string(jsonData))))

	os.Exit(1)
}
