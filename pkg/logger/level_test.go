package logger

import (
	"errors"
	"log/slog"
	"testing"
)

func TestParseLevel(t *testing.T) {
	testCases := []struct {
		in  string
		out slog.Level
		err error
	}{
		{"", LevelInfo, errInvalidLevel},
		{"string", LevelInfo, errInvalidLevel},
		{"STRING", LevelInfo, errInvalidLevel},
		{stringLevelTrace, LevelTrace, nil},
		{stringLevelDebug, LevelDebug, nil},
		{stringLevelInfo, LevelInfo, nil},
		{stringLevelWarn, LevelWarn, nil},
		{stringLevelError, LevelError, nil},
		{stringLevelFatal, LevelFatal, nil},
		{"trace", LevelTrace, nil},
		{"TrAcE", LevelTrace, nil},
	}

	for _, test := range testCases {
		t.Run(test.in, func(t *testing.T) {
			v, err := parseLevel(test.in)
			if test.out != v {
				t.Fatalf(`wrong result for level "%v". Expected %v, got %v`, test.in, test.out, v)
			}
			if err != nil && !errors.Is(err, test.err) {
				t.Fatalf(`unexpected error parsing env. Expected "%v", got "%v"`, test.err, err)
			}
		})
	}
}

func TestStringifyLevel(t *testing.T) {
	testCases := map[string]struct {
		in  slog.Level
		out slog.Value
	}{
		"trace": {LevelTrace, slog.StringValue(stringLevelTrace)},
		"debug": {LevelDebug, slog.StringValue(stringLevelDebug)},
		"info":  {LevelInfo, slog.StringValue(stringLevelInfo)},
		"warn":  {LevelWarn, slog.StringValue(stringLevelWarn)},
		"error": {LevelError, slog.StringValue(stringLevelError)},
		"fatal": {LevelFatal, slog.StringValue(stringLevelFatal)},
		"-1000": {slog.Level(-1000), slog.StringValue(stringLevelTrace)},
		"1000":  {slog.Level(1000), slog.StringValue(stringLevelFatal)},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			v := stringifyLevel(test.in)
			if !v.Equal(test.out) {
				t.Fatalf(`wrong result for input "%v". Expected %v, got %v`, test.in, test.out, v)
			}
		})
	}
}
