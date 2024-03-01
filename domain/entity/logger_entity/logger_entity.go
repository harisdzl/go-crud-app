package logger_entity

import (
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Logger struct {
	Honeycomb trace.Tracer
	Zap *zap.Logger
}