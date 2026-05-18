package loghelper

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"gitee.com/cai-zixiang_hainan/wt/internal/config"
)

// LogEntry represents a single log entry in JSON format
type LogEntry struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Message string `json:"msg"`
	// ExtraFields captures any other key-value pairs in the log line
	ExtraFields map[string]interface{} `json:"-"`
}

// UnmarshalJSON implements custom unmarshaling to capture extra fields
func (e *LogEntry) UnmarshalJSON(data []byte) error {
	// First, unmarshal into a generic map to find all keys
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Extract known fields
	if t, ok := raw["time"]; ok {
		e.Time = fmt.Sprintf("%v", t)
		delete(raw, "time")
	}
	if l, ok := raw["level"]; ok {
		e.Level = fmt.Sprintf("%v", l)
		delete(raw, "level")
	}
	if m, ok := raw["msg"]; ok {
		e.Message = fmt.Sprintf("%v", m)
		delete(raw, "msg")
	}

	// The remaining keys are extra fields
	e.ExtraFields = raw
	return nil
}

// ReadServerLog reads and displays the server log file
func ReadServerLog() error {
	logPath := config.GetServerLogPath()

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return fmt.Errorf("log file not found: %s", logPath)
	}

	file, err := os.Open(logPath)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	fmt.Println("=== Server Log ===")
	fmt.Printf("File: %s\n\n", logPath)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err == nil {
			fmt.Println(formatLogEntry(entry))
		} else {
			fmt.Println(line)
		}
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading log file: %w", err)
	}

	fmt.Printf("\n=== Total: %d log entries ===\n", lineCount)
	return nil
}

// ReadServerLogTail reads and displays the last N lines of the server log
func ReadServerLogTail(n int) error {
	logPath := config.GetServerLogPath()

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return fmt.Errorf("log file not found: %s", logPath)
	}

	file, err := os.Open(logPath)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading log file: %w", err)
	}

	start := 0
	if len(lines) > n {
		start = len(lines) - n
	}

	fmt.Printf("=== Last %d Log Entries ===\n\n", n)
	for i := start; i < len(lines); i++ {
		var entry LogEntry
		if err := json.Unmarshal([]byte(lines[i]), &entry); err == nil {
			fmt.Println(formatLogEntry(entry))
		} else {
			fmt.Println(lines[i])
		}
	}

	return nil
}

// SearchServerLog searches for entries containing the specified keyword
func SearchServerLog(keyword string) error {
	logPath := config.GetServerLogPath()

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return fmt.Errorf("log file not found: %s", logPath)
	}

	file, err := os.Open(logPath)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	matchCount := 0
	keyword = strings.ToLower(keyword)

	fmt.Printf("=== Searching for '%s' ===\n\n", keyword)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), keyword) {
			var entry LogEntry
			if err := json.Unmarshal([]byte(line), &entry); err == nil {
				fmt.Println(formatLogEntry(entry))
			} else {
				fmt.Println(line)
			}
			matchCount++
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading log file: %w", err)
	}

	fmt.Printf("\n=== Found %d matching entries ===\n", matchCount)
	return nil
}

// formatLogEntry formats a LogEntry into a human-readable string
func formatLogEntry(entry LogEntry) string {
	// Parse time for better formatting
	timeStr := entry.Time
	if t, err := time.Parse(time.RFC3339, entry.Time); err == nil {
		timeStr = t.Format("2006-01-02 15:04:05")
	}

	// Build formatted output
	result := fmt.Sprintf("[%s] %-5s | %s", timeStr, entry.Level, entry.Message)

	// ✅ 核心修复：遍历并打印所有额外字段
	if len(entry.ExtraFields) > 0 {
		var extras []string
		for k, v := range entry.ExtraFields {
			extras = append(extras, fmt.Sprintf("%s=%v", k, v))
		}
		result += " | " + strings.Join(extras, ", ")
	}

	return result
}

