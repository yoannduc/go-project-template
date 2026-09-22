package logger

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/yoannduc/go-project-template/pkg/env"
)

const (
	// logLevelEnvVar is the env var name to check log level against.
	logLevelEnvVar = "LOG_LEVEL"
	// handlerTypeEnvVar is the env var name to check log handler against.
	handlerTypeEnvVar = "LOG_HANDLER"
)

var (
	// envV is a global var to be loaded only once because
	// its value is needed for EACH attr of EACH log line.
	// It would amount to too many calls otherwise.
	envV = env.Get()
)

// handlerOptionReplaceAttr is the function which transforms certain
// log attributes format or value. It removes time entirely for
// test consistency purposes and replace level key value from int
// level representation to its string representation.
func handlerOptionReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	// Remove time from the output for predictable test output.
	if envV.IsTest() && a.Key == slog.TimeKey {
		return slog.Attr{}
	}

	if a.Key == slog.LevelKey {
		// Handle custom level values.
		level := a.Value.Any().(slog.Level)

		a.Value = stringifyLevel(level)
	}

	return a
}

var (
	errInvalidHandler = errors.New("invalid log handler")
)

// Strings for handled handlers
const (
	jsonHandler = "JSON"
	textHandler = "TEXT"
)

var (
	handlerOptions = &slog.HandlerOptions{
		AddSource: true,
		Level: func() slog.Leveler {
			lvl, _ := parseLevel(os.Getenv(logLevelEnvVar))
			return lvl
		}(),
		ReplaceAttr: handlerOptionReplaceAttr,
	}
)

// parseHandler returns the slog.Handler based on handler type string.
// It returns json handler and an error if string did not match
// any handled handler or is empty. It is case-insensitive.
func parseHandler(hdl string) (slog.Handler, error) {
	switch strings.ToUpper(hdl) {
	case jsonHandler:
		return slog.NewJSONHandler(os.Stderr, handlerOptions), nil
	case textHandler:
		return slog.NewTextHandler(os.Stderr, handlerOptions), nil
	default:
		return slog.NewJSONHandler(os.Stderr, handlerOptions), errInvalidHandler
	}
}

// Get returns the instanciated logger with config selected params.
// It uses singleton design pattern to configure and instanciate only once.
func Get() *slog.Logger {
	h, _ := parseHandler(os.Getenv(handlerTypeEnvVar))

	return slog.New(NewContextHandler(h))
}
