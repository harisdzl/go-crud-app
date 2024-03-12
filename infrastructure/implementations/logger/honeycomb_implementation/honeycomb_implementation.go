package honeycomb

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// HoneycombRepo is the Honeycomb logger implementation
type HoneycombRepo struct {
    c *gin.Context
	span trace.Span
}



func NewHoneycombRepository() *HoneycombRepo {
    return &HoneycombRepo{nil, nil}
}

// Start Honeycomb
func (h *HoneycombRepo) Start(c *gin.Context, info string) trace.Span {
    h.c = c
    _, file, line, _ := runtime.Caller(2)
    // Retrieve otel_context from Gin context if it exists
    ctxValue, exists := c.Get("otel_context")
    // Implement info logging with Honeycomb using the retrieved context
    var ctx context.Context
    if exists {
        log.Println("Otel context exists: Getting this from " + info)
        ctx, _ = ctxValue.(context.Context)
    } else {
        log.Println("Getting this from" + info)
        ctx = c.Request.Context()
    }
    tracer := otel.Tracer("")
    _, span := tracer.Start(ctx, info, trace.WithAttributes(
        attribute.String("file", file),
        attribute.String("client_ip", h.c.ClientIP()),
        attribute.Int("line", line),
    ))
    h.span = span

    return span
}

// Debug logs a debug message
func (h *HoneycombRepo) Debug(msg string, fields map[string]interface{}) {
	jsonData, jsonDataErr := json.Marshal(fields)
    _, file, line, _ := runtime.Caller(2)

	if jsonDataErr != nil {
		return
	}

	h.span.AddEvent(msg, trace.WithAttributes(
		attribute.String("level", "Debug"),
		attribute.String("data", string(jsonData)),
        attribute.String("file", file),
        attribute.String("client_ip", h.c.ClientIP()),
        attribute.Int("line", line)))
	
}

// Info logs an info message
func (h *HoneycombRepo) Info(msg string, fields map[string]interface{}) {
	jsonData, jsonDataErr := json.Marshal(fields)
    _, file, line, _ := runtime.Caller(2)

	if jsonDataErr != nil {
        log.Println("jsonData marshal error")
		return
	}
    
	h.span.AddEvent(msg, trace.WithAttributes(
		attribute.String("level", "Info"),
		attribute.String("data", string(jsonData)),
        attribute.String("file", file),
        attribute.String("client_ip", h.c.ClientIP()),
        attribute.Int("line", line)))
}

func (h *HoneycombRepo) Error(msg string, fields map[string]interface{}) {
    jsonData, jsonDataErr := json.Marshal(fields)
    _, file, line, _ := runtime.Caller(2)
    if jsonDataErr != nil {
        return
    }

    h.span.RecordError(errors.New(msg), trace.WithAttributes(
		attribute.String("level", "Error"),
        attribute.String("data", string(jsonData)),
        attribute.String("file", file),
        attribute.String("client_ip", h.c.ClientIP()),
        attribute.Int("line", line)))
}


// Warn logs a warning message
func (h *HoneycombRepo) Warn(msg string, fields map[string]interface{}) {
    // Implement warn logging with Honeycomb
    // This function can be implemented similar to the Info function
    jsonData, jsonDataErr := json.Marshal(fields)
    _, file, line, _ := runtime.Caller(2)
    if jsonDataErr != nil {
        return
    }

    h.span.AddEvent(msg, trace.WithAttributes(
		attribute.String("level", "Warn"),
        attribute.String("data", string(jsonData)),
        attribute.String("file", file),
        attribute.String("client_ip", h.c.ClientIP()),
        attribute.Int("line", line)))
}

// Fatal logs a fatal message
func (h *HoneycombRepo) Fatal(msg string, fields map[string]interface{}) {
    // Implement fatal logging with Honeycomb
    // This function can be implemented similar to the Info function
    _, file, line, _ := runtime.Caller(2)

    jsonData, jsonDataErr := json.Marshal(fields)
    if jsonDataErr != nil {
        return
    }

    h.span.AddEvent(msg, trace.WithAttributes(
		attribute.String("level", "Fatal"),
        attribute.String("data", string(jsonData)),
        attribute.String("file", file),
        attribute.String("client_ip", h.c.ClientIP()),
        attribute.Int("line", line)))

	os.Exit(1)
}

func (h *HoneycombRepo) GetSpan() trace.Span {
    return h.span
}

func (h *HoneycombRepo) GetContext() *gin.Context {
    return h.c
}