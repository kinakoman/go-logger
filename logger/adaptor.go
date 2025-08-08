package logger

import (
	"go.uber.org/zap"
)

type LoggerManager struct {
	logger *zap.Logger
}
