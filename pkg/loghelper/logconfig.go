package loghelper

import (
	"log/slog"
	"os"
)

func JSONLog(isDEBUG bool) {
	// 强制使用 JSON 格式
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	if isDEBUG {
		opts = &slog.HandlerOptions{Level: slog.LevelDebug}
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))

	// 立即打印一条测试日志
	slog.Debug("System initialized", "format_check", "this_should_be_json")
}
