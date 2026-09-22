package logger

import (
	"errors"
	"log/slog"
	"strings"
)

var (
	errInvalidLevel = errors.New("invalid level value")
)

// Name for handled levels.
const (
	LevelTrace = slog.Level(-8)
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
	LevelFatal = slog.Level(12)
)

// Strings for handled log levels.
const (
	stringLevelTrace = "TRACE"
	stringLevelDebug = "DEBUG"
	stringLevelInfo  = "INFO"
	stringLevelWarn  = "WARN"
	stringLevelError = "ERROR"
	stringLevelFatal = "FATAL"
)

// parseLevel returns the log level based on lvl string.
// It returns info level and an error if string did not match
// any level handled or was empty. It is case-insensitive.
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

// stringifyLevel returns the string representation of the
// log level as a slog.Value.
func stringifyLevel(lvl slog.Level) slog.Value {
	switch {
	case lvl < slog.LevelDebug:
		return slog.StringValue(stringLevelTrace)
	case lvl < slog.LevelInfo:
		return slog.StringValue(stringLevelDebug)
	case lvl < slog.LevelWarn:
		return slog.StringValue(stringLevelInfo)
	case lvl < slog.LevelError:
		return slog.StringValue(stringLevelWarn)
	case lvl < LevelFatal:
		return slog.StringValue(stringLevelError)
	default:
		return slog.StringValue(stringLevelFatal)
	}
}
