package logger

import (
	"errors"
	"log/slog"
	"strings"
)

var (
	errInvalidLevel = errors.New("invalid level value")
)

const (
	logLevelEnvVar = "LOG_LEVEL"
)

const (
	LevelTrace = slog.Level(-8)
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
	LevelFatal = slog.Level(12)
)

const (
	stringLevelTrace = "TRACE"
	stringLevelDebug = "DEBUG"
	stringLevelInfo  = "INFO"
	stringLevelWarn  = "WARN"
	stringLevelError = "ERROR"
	stringLevelFatal = "FATAL"
)

func parseLevel(lvl string) (slog.Level, error) {
	switch strings.ToUpper(lvl) {
	case stringLevelTrace:
		return LevelTrace, nil
	case stringLevelDebug:
		return slog.LevelDebug, nil
	case stringLevelInfo:
		return slog.LevelInfo, nil
	case stringLevelWarn:
		return slog.LevelWarn, nil
	case stringLevelError:
		return slog.LevelError, nil
	case stringLevelFatal:
		return LevelFatal, nil
	default:
		return slog.LevelInfo, errInvalidLevel
	}
}
