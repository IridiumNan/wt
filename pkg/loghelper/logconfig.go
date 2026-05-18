package loghelper

import (
    "io"
    "log/slog"
    "os"
    "path/filepath"

    "gitee.com/cai-zixiang_hainan/wt/internal/config"
)

// InitClientLogger initializes slog for client-side logging (text format to stdout)
func InitClientLogger(isDebug bool) {
    level := slog.LevelInfo
    if isDebug {
        level = slog.LevelDebug
    }

    opts := &slog.HandlerOptions{Level: level}
    
    // Use text handler for better human readability in CLI
    handler := slog.NewTextHandler(os.Stdout, opts)
    slog.SetDefault(slog.New(handler))

    if isDebug {
        slog.Debug("Client logger initialized", "mode", "text", "output", "stdout")
    }
}

// InitServerLogger initializes slog to write to a file for the server (JSON format)
func InitServerLogger(debug bool) error {
    // Ensure directory exists
    logPath := config.GetServerLogPath()
    dir := filepath.Dir(logPath)
    if err := os.MkdirAll(dir, 0o755); err != nil {
        return err
    }

    // Open log file in append mode
    f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
    if err != nil {
        return err
    }

    level := slog.LevelInfo
    if debug {
        level = slog.LevelDebug
    }

    // Create a handler that writes to the file in JSON format
    opts := &slog.HandlerOptions{Level: level}
    
    // Use multi-writer to output to both file and console (optional)
    // For production, you might want only file output
    var writer io.Writer = f
    if debug {
        // In debug mode, also print to console for easier debugging
        writer = io.MultiWriter(f, os.Stdout)
    }
    
    handler := slog.NewJSONHandler(writer, opts)
    slog.SetDefault(slog.New(handler))

    slog.Info("Server logger initialized", "path", logPath, "format", "json")
    return nil
}

// JSONLog is kept for backward compatibility but deprecated
// Use InitClientLogger or InitServerLogger instead
func JSONLog(isDEBUG bool) {
    InitClientLogger(isDEBUG)
}