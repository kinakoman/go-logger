package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
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

func NewLogger(args ...interface{}) (*Logger, error) {
	var s string
	var config LoggerConfig
	var hasS, hasConfig bool

	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			s = v
			hasS = true
		case LoggerConfig:
			config = v
			hasConfig = true
		default:
			fmt.Printf("未知の型: %T\n", v)
		}
	}

	switch {
	case hasS:
		return newLogger(s)
	case !hasS && hasConfig:
		return newLoggerWithLumberjack(config)
	case hasS && hasConfig:
		// configが指定されている場合は、configを使用してログファイルを設定
		return newLoggerWithLumberjack(config)
	}

	return newLogger()
}

func newLoggerWithLumberjack(config LoggerConfig) (*Logger, error) {
	encoder := newSimpleEncoder()

	if config.Filename == "" {
		return nil, fmt.Errorf("filename must be specified in LoggerConfig")
	}

	if config.MaxSize <= 0 {
		config.MaxSize = 1 // デフォルト値: 100MB
	}

	if config.MaxBackups <= 0 {
		config.MaxBackups = 1 // デフォルト値: 無制限
	}

	// lumberjack ローテーション設定
	rotator := &lumberjack.Logger{
		Filename:   config.Filename,   // 例: "logs/app.log"
		MaxSize:    config.MaxSize,    // MB単位で最大サイズ
		MaxBackups: config.MaxBackups, // 最大バックアップ数
		MaxAge:     config.MaxAge,     // 最大保持日数
		Compress:   config.Compress,   // 圧縮する場合
	}

	// コンソール出力
	consoleSync := zapcore.Lock(os.Stdout)
	consoleCore := zapcore.NewCore(encoder, consoleSync, zapcore.InfoLevel)

	// ファイル出力（ローテーション付き）
	fileSync := zapcore.AddSync(rotator)
	fileCore := zapcore.NewCore(encoder, fileSync, zapcore.InfoLevel)

	// Teeで両方に出力
	core := zapcore.NewTee(consoleCore, fileCore)

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
	return &Logger{logger: zapLogger}, nil
}

// path: 相対パス（例: "logs/app.log"）
func newLogger(paths ...string) (*Logger, error) {

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
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))

	return &Logger{logger: zapLogger}, nil
}

func (l *Logger) Info(msg string) {
	l.logger.Info(msg)
}

func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}

func (l *Logger) Warn(msg string) {
	l.logger.Warn(msg)
}
