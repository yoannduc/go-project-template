package logger

import (
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/yoannduc/go-project-template/pkg/singleton"
)

const (
	handlerTypeEnvVar = "LOG_HANDLER"
	textHandler       = "TEXT"
)

var (
	once  sync.Once
	singl *singleton.Singleton[*slog.Logger]

	handlerOptions = &slog.HandlerOptions{
		AddSource:   true,
		Level:       handlerOptionLevel(),
		ReplaceAttr: handlerOptionReplaceAttr,
	}
)

func Get() *slog.Logger {
	if singl == nil {
		once.Do(func() {
			var h slog.Handler = slog.NewJSONHandler(os.Stderr, handlerOptions)
			if strings.ToUpper(os.Getenv(handlerTypeEnvVar)) == textHandler {
				h = slog.NewTextHandler(os.Stderr, handlerOptions)
			}

			singl = &singleton.Singleton[*slog.Logger]{Instance: slog.New(NewContextHandler(h))}
		})
	}

	return singl.Instance
}

func handlerOptionLevel() slog.Leveler {
	// Configure log level, defaulting to info level
	if lvl, err := parseLevel(os.Getenv(logLevelEnvVar)); os.Getenv(logLevelEnvVar) != "" && err == nil {
		return lvl
	} else {
		return slog.LevelInfo
	}
}

func handlerOptionReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	// Remove time from the output for predictable test output.
	if a.Key == slog.TimeKey {
		return slog.Attr{}
	}

	if a.Key == slog.LevelKey {
		// Handle custom level values.
		level := a.Value.Any().(slog.Level)

		switch {
		case level < slog.LevelDebug:
			a.Value = slog.StringValue(stringLevelTrace)
		case level < slog.LevelInfo:
			a.Value = slog.StringValue(stringLevelDebug)
		case level < slog.LevelWarn:
			a.Value = slog.StringValue(stringLevelInfo)
		case level < slog.LevelError:
			a.Value = slog.StringValue(stringLevelWarn)
		case level < LevelFatal:
			a.Value = slog.StringValue(stringLevelError)
		default:
			a.Value = slog.StringValue(stringLevelFatal)
		}
	}

	return a
}
