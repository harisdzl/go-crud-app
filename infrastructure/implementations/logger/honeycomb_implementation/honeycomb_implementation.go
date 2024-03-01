package honeycomb

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// HoneycombRepo is the Honeycomb logger implementation
type HoneycombRepo struct {
	p *base.Persistence
	Context *context.Context
	Span trace.Span
}

// NewHoneycombRepository creates a new Honeycomb logger repository
func NewHoneycombRepository(p *base.Persistence, c *context.Context, info string) *HoneycombRepo {
	// Implement info logging with Honeycomb
	tracer := p.Logger.Honeycomb
	
	ctx, span := tracer.Start(*c, info)

	return &HoneycombRepo{p, &ctx, span}
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
