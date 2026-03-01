package logger

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// GetLogger returns a logger instance with a package field pre-configured
func GetLogger(packageName string) *logrus.Entry {
	return logrus.WithField("package", packageName)
}

// GetLoggerWithFields returns a logger with multiple fields
func GetLoggerWithFields(fields logrus.Fields) *logrus.Entry {
	return logrus.WithFields(fields)
}

// CallerHook 自動添加 caller 資訊的 hook
type CallerHook struct {
	Field     string // 預設 "caller"
	Skip      int    // 預設 8 (跳過 logrus 和 hook 內部的 stack frames)
	LogLevels []logrus.Level
	Detailed  bool // true: 顯示完整路徑, false: 只顯示檔名
}

// NewCallerHook creates a new caller hook
func NewCallerHook(detailed bool) *CallerHook {
	return &CallerHook{
		Field:     "caller",
		Skip:      7,
		Detailed:  detailed,
		LogLevels: logrus.AllLevels,
	}
}

// Levels returns the log levels this hook applies to
func (hook *CallerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

// Fire adds caller information to the log entry
func (hook *CallerHook) Fire(entry *logrus.Entry) error {
	// 使用 runtime.Caller 取得呼叫者資訊
	if pc, file, line, ok := runtime.Caller(hook.Skip); ok {
		funcName := runtime.FuncForPC(pc).Name()

		// 簡化 function name (去掉 package path)
		if lastSlash := strings.LastIndex(funcName, "/"); lastSlash >= 0 {
			funcName = funcName[lastSlash+1:]
		}

		// 簡化檔案路徑
		if !hook.Detailed {
			if lastSlash := strings.LastIndex(file, "/"); lastSlash >= 0 {
				file = file[lastSlash+1:]
			} else if lastBackslash := strings.LastIndex(file, "\\"); lastBackslash >= 0 {
				// Windows path support
				file = file[lastBackslash+1:]
			}
		}

		// 添加到 entry.Data
		entry.Data["function"] = funcName
		entry.Data["file"] = file
		entry.Data["line"] = line
	}
	return nil
}

// PackageFormatter formats logs with package name prominently displayed
type PackageFormatter struct {
	// Include timestamp
	TimestampFormat string
}

// Format formats a log entry
func (f *PackageFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(f.TimestampFormat)

	// Extract package name if available
	packageName := "unknown"
	if pkg, ok := entry.Data["package"]; ok {
		packageName = fmt.Sprintf("%v", pkg)
	}

	// Extract caller information
	funcName := ""
	fileName := ""
	lineNo := 0
	if fn, ok := entry.Data["function"]; ok {
		funcName = fmt.Sprintf("%v", fn)
	}
	if file, ok := entry.Data["file"]; ok {
		fileName = fmt.Sprintf("%v", file)
	}
	if line, ok := entry.Data["line"]; ok {
		lineNo = line.(int)
	}

	// Combine caller information
	caller := ""
	if fileName != "" && lineNo > 0 {
		caller = fmt.Sprintf("[%s:%d %s] ", fileName, lineNo, funcName)
	}

	// Format: [timestamp] [LEVEL] [package] [caller] message
	logLine := fmt.Sprintf("[%s] [%-5s] [%-10s] %s%s",
		timestamp,
		strings.ToUpper(entry.Level.String()),
		packageName,
		caller,
		entry.Message,
	)

	// Add any additional fields (excluding 'package', 'function', 'file', 'line')
	excludeFields := map[string]bool{
		"package":  true,
		"function": true,
		"file":     true,
		"line":     true,
	}
	for k, v := range entry.Data {
		if !excludeFields[k] {
			logLine += fmt.Sprintf(" %s=%v", k, v)
		}
	}

	logLine += "\n"
	return []byte(logLine), nil
}
