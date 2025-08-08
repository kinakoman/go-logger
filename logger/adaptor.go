package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.Logger
}

type NewLogger struct {
	logger *zap.Logger
	path string
}
