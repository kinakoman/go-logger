package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.Logger
}

type LoggerConfig struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}
