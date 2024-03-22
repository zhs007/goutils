package goutils

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"log/slog"

	"gopkg.in/natefinch/lumberjack.v2"
)

func parseLevel(str string) slog.Level {
	str = strings.ToLower(str)

	if str == "debug" {
		return slog.LevelDebug
	} else if str == "warn" {
		return slog.LevelWarn
	} else if str == "error" {
		return slog.LevelError
	}

	return slog.LevelInfo
}

func InitLogger2(appName string, appVersion string, strLevel string, isConsole bool, logpath string) {
	lv := parseLevel(strLevel)

	var lvl = &slog.LevelVar{}
	lvl.Set(lv)
	opts := &slog.HandlerOptions{
		Level: lvl,
	}

	if isConsole {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, opts)))
	} else {
		logger := &lumberjack.Logger{
			Filename:   path.Join(logpath, fmt.Sprintf("%v.log", appName)),
			MaxBackups: 99,
			MaxAge:     30,   //days
			Compress:   true, // disabled by default
		}

		slog.SetDefault(slog.New(slog.NewJSONHandler(logger, opts)))
	}
}

// Debug logs a debug message with the given fields
func Debug(message string, fields ...slog.Attr) {
	slog.LogAttrs(context.Background(), slog.LevelDebug, message, fields...)
}

// Info logs a debug message with the given fields
func Info(message string, fields ...slog.Attr) {
	slog.LogAttrs(context.Background(), slog.LevelInfo, message, fields...)
}

// Warn logs a debug message with the given fields
func Warn(message string, fields ...slog.Attr) {
	slog.LogAttrs(context.Background(), slog.LevelWarn, message, fields...)
}

// Error logs a debug message with the given fields
func Error(message string, fields ...slog.Attr) {
	slog.LogAttrs(context.Background(), slog.LevelError, message, fields...)
}

func Err(err error) slog.Attr {
	return slog.Any("err", err)
}
