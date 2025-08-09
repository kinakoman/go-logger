package logger

import "testing"

func Test_Logger(t *testing.T) {
	// ファイル指定せずてLoggerを作成
	// 出力はコンソールのみ
	logger, err := NewLogger()
	if err != nil {
		t.Fatalf("Loggerの作成に失敗: %v", err)
	}

	logger.Info("This is an info message")
	logger.Error("This is an error message")
	logger.Warn("This is a warning message")

	// ファイルを指定してLoggerを作成
	// 出力はコンソールとファイルの両方
	// ログは全て同じファイルに出力される
	loggerWithFile, err := NewLogger("test_log.log")
	if err != nil {
		t.Fatalf("Loggerの作成に失敗: %v", err)
	}

	loggerWithFile.Info("This is an info message")
	loggerWithFile.Error("This is an error message")
	loggerWithFile.Warn("This is a warning message")

	// ローテーション付きのLoggerを作成
	// ログコンフィグで設定を指定
	// ログファイルの大きさがMaxSizeを超えるとローテーションされる
	loggerWithRoatation, err := NewLogger(LoggerConfig{
		Filename:   "test_rotation.log",
		MaxSize:    1, // MB単位で最大サイズ
		MaxBackups: 1,
		MaxAge:     28,
		Compress:   false, // 古いファイルの圧縮指定
	})
	if err != nil {
		t.Fatalf("Logger with rotation creation failed: %v", err)
	}

	for i := 0; i < 50000; i++ {
		loggerWithRoatation.Info("テストログ出力")
	}

	t.Log("Logger test completed successfully")
}
