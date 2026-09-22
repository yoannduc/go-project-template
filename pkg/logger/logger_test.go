package logger

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestHandlerOptionReplaceAttr(t *testing.T) {
	now := time.Now()

	testCases := map[string]struct {
		in  slog.Attr
		out slog.Attr
	}{
		"empty, do nothing": {slog.Attr{}, slog.Attr{}},
		"do nothing": {
			slog.Attr{Key: "Key", Value: slog.StringValue("Value")},
			slog.Attr{Key: "Key", Value: slog.StringValue("Value")},
		},
		"time key, no effect": {
			slog.Time(slog.TimeKey, now),
			slog.Time(slog.TimeKey, now),
		},
		"level key, replace by string level (info)": {
			slog.Attr{Key: slog.LevelKey, Value: slog.AnyValue(LevelInfo)},
			slog.Attr{Key: slog.LevelKey, Value: slog.StringValue(stringLevelInfo)},
		},
		"level key, replace by string level (error)": {
			slog.Attr{Key: slog.LevelKey, Value: slog.AnyValue(LevelError)},
			slog.Attr{Key: slog.LevelKey, Value: slog.StringValue(stringLevelError)},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			v := handlerOptionReplaceAttr([]string{}, test.in)
			if !v.Equal(test.out) {
				t.Fatalf(`wrong result for input "%v". Expected %v, got %v`, test.in, test.out, v)
			}
		})
	}
}

func TestParseHandler(t *testing.T) {
	testCases := []struct {
		in  string
		out slog.Handler
		err error
	}{
		{"", slog.NewJSONHandler(os.Stderr, handlerOptions), errInvalidHandler},
		{"TOTO", slog.NewJSONHandler(os.Stderr, handlerOptions), errInvalidHandler},
		{jsonHandler, slog.NewJSONHandler(os.Stderr, handlerOptions), nil},
		{"json", slog.NewJSONHandler(os.Stderr, handlerOptions), nil},
		{"JsOn", slog.NewJSONHandler(os.Stderr, handlerOptions), nil},
		{textHandler, slog.NewTextHandler(os.Stderr, handlerOptions), nil},
	}

	for _, test := range testCases {
		t.Run(test.in, func(t *testing.T) {
			fmt.Println()
			v, err := parseHandler(test.in)

			outType := reflect.TypeOf(test.out)
			if outType.Kind() == reflect.Ptr {
				outType = outType.Elem()
			}

			vType := reflect.TypeOf(v)
			if vType.Kind() == reflect.Ptr {
				vType = vType.Elem()
			}

			// reflect.DeepEqual returns false because of ReplaceAttr,
			// so we check only type equality for now.
			// if !reflect.DeepEqual(v, test.out) {
			if vType != outType {
				t.Fatalf(`wrong type equality for input "%v". Expected %v, got %v`, test.in, outType, vType)
			}
			if err != nil && !errors.Is(err, test.err) {
				t.Fatalf(`unexpected error for input "%v". Expected "%v", got "%v"`, test.in, test.err, err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	testCases := []struct {
		out *slog.Logger
		env string
	}{
		{slog.New(slog.NewJSONHandler(os.Stderr, handlerOptions)), "JSON"},
		{slog.New(slog.NewTextHandler(os.Stderr, handlerOptions)), "TEXT"},
		{slog.New(slog.NewJSONHandler(os.Stderr, handlerOptions)), "UNHANDLED"},
		{slog.New(slog.NewJSONHandler(os.Stderr, handlerOptions)), ""},
	}

	for _, test := range testCases {
		t.Run(test.env, func(t *testing.T) {
			t.Setenv(handlerTypeEnvVar, test.env)

			v := Get()
			outType := reflect.TypeOf(test.out)
			if outType.Kind() == reflect.Ptr {
				outType = outType.Elem()
			}

			vType := reflect.TypeOf(v)
			if vType.Kind() == reflect.Ptr {
				vType = vType.Elem()
			}

			// reflect.DeepEqual returns false because of ReplaceAttr,
			// so we check only type equality for now.
			// if !reflect.DeepEqual(v, test.out) {
			if vType != outType {
				t.Fatalf(`wrong type equality for env "%v". Expected %v, got %v`, test.env, outType, vType)
			}
		})
	}
}
