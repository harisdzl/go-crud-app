package logger_entity

import (
	"go.uber.org/zap"
)

type Logger struct {
	Zap *zap.Logger
}