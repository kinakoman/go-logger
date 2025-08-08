package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// カスタムエンコーダ
func newSimpleEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:    "T",
		LevelKey:   "L",
		MessageKey: "M",
		CallerKey:  "C",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeLevel: func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			switch level {
			case zapcore.InfoLevel:
				enc.AppendString("[INFO]")
			case zapcore.ErrorLevel:
				enc.AppendString("[ERROR]")
			case zapcore.WarnLevel:
				enc.AppendString("[WARN]")
			default:
				enc.AppendString("[" + level.String() + "]")
			}
		},
		EncodeCaller: zapcore.ShortCallerEncoder, // e.g. logger/logger.go:52
	})
}

// path: 相対パス（例: "logs/app.log"）
func NewLogger(paths ...string) (*LoggerManager, error) {

	encoder := newSimpleEncoder()

	consoleSync := zapcore.Lock(os.Stdout)

	var cores []zapcore.Core

	// paths が指定されていればファイル出力を追加
	if len(paths) > 0 && paths[0] != "" {
		path := paths[0]
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		fileSync := zapcore.AddSync(file)
		cores = append(cores, zapcore.NewCore(encoder, fileSync, zapcore.InfoLevel))
	}

	cores = append(cores, zapcore.NewCore(encoder, consoleSync, zapcore.InfoLevel))

	core := zapcore.NewTee(cores...)

	// caller情報を含める（1段上を呼び出し元として）
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &LoggerManager{logger: zapLogger}, nil
}

func (l *LoggerManager) Info(msg string) {
	l.logger.Info(msg)
}

func (l *LoggerManager) Error(msg string) {
	l.logger.Error(msg)
}

func (l *LoggerManager) Warn(msg string) {
	l.logger.Warn(msg)
}
