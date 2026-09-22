package logger

import (
	"context"
	"log/slog"
	"net"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/yoannduc/go-project-template/pkg/traceid"
	"github.com/yoannduc/go-project-template/pkg/userip"
)

func TestNewContextHandler(t *testing.T) {
	testCases := map[string]struct {
		in  slog.Handler
		out *ContextHandler
	}{
		"wraps JSONHandler":            {slog.NewJSONHandler(os.Stderr, nil), &ContextHandler{slog.NewJSONHandler(os.Stderr, nil)}},
		"wraps TextHandler":            {slog.NewTextHandler(os.Stderr, nil), &ContextHandler{slog.NewTextHandler(os.Stderr, nil)}},
		"does not wrap ContextHandler": {NewContextHandler(slog.NewTextHandler(os.Stderr, nil)), &ContextHandler{slog.NewTextHandler(os.Stderr, nil)}},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			v := NewContextHandler(test.in)

			if !reflect.DeepEqual(v, test.out) {
				t.Fatalf(`wrong result for input "%v". Expected %v, got %v`, test.in, test.out, v)
			}
		})
	}
}

func TestEnabled(t *testing.T) {
	testCases := []struct {
		in slog.Level
	}{
		{slog.LevelDebug},
		{slog.LevelInfo},
		{slog.LevelWarn},
		{slog.LevelError},
	}

	for _, test := range testCases {
		t.Run(test.in.String(), func(t *testing.T) {
			ctx := context.Background()
			hdl := slog.NewJSONHandler(os.Stderr, nil)
			v := NewContextHandler(hdl)

			if !reflect.DeepEqual(v.Enabled(ctx, test.in), hdl.Enabled(ctx, test.in)) {
				t.Fatalf(`enabled did not match between parent handler & wrapper`)
			}
		})
	}
}

func TestHandle(t *testing.T) {
	testCases := map[string]struct {
		traceid string
		userip  net.IP
		in      slog.Record
	}{
		"all empty":    {"", net.IP{}, slog.NewRecord(time.Now(), slog.LevelInfo, "", uintptr(0))},
		"traceid only": {"1234", net.IP{}, slog.NewRecord(time.Now(), slog.LevelInfo, "", uintptr(0))},
		"userip only":  {"", net.ParseIP("8.8.8.8"), slog.NewRecord(time.Now(), slog.LevelInfo, "", uintptr(0))},
		"both":         {"1234", net.ParseIP("8.8.8.8"), slog.NewRecord(time.Now(), slog.LevelInfo, "", uintptr(0))},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			if test.traceid != "" {
				ctx = traceid.NewContext(ctx, test.traceid)
			}
			if !reflect.DeepEqual(test.userip, net.IP{}) {
				ctx = userip.NewContext(ctx, test.userip)
			}

			hdl := slog.NewJSONHandler(os.Stderr, nil)
			v := NewContextHandler(hdl)

			if v.Handle(ctx, test.in) != hdl.Handle(ctx, test.in) {
				t.Fatalf(`enabled did not match between parent handler & wrapper`)
			}
		})
	}
}

func TestWithAttrs(t *testing.T) {
	testCases := map[string]struct {
		in []slog.Attr
	}{
		"empty":  {[]slog.Attr{}},
		"any":    {[]slog.Attr{slog.Any("toto", "toto")}},
		"string": {[]slog.Attr{slog.String("toto", "toto")}},
		"bool":   {[]slog.Attr{slog.Bool("toto", true)}},
		"int":    {[]slog.Attr{slog.Int("toto", 123)}},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			hdl := slog.NewJSONHandler(os.Stderr, nil)
			v := NewContextHandler(hdl)

			if !reflect.DeepEqual(v.WithAttrs(test.in), hdl.WithAttrs(test.in)) {
				t.Fatalf(`enabled did not match between parent handler & wrapper`)
			}
		})
	}
}

func TestWithGroup(t *testing.T) {
	testCases := []struct {
		in string
	}{
		{""},
		{"group"},
	}

	for _, test := range testCases {
		t.Run(test.in, func(t *testing.T) {
			hdl := slog.NewJSONHandler(os.Stderr, nil)
			v := NewContextHandler(hdl)

			if !reflect.DeepEqual(v.WithGroup(test.in), hdl.WithGroup(test.in)) {
				t.Fatalf(`enabled did not match between parent handler & wrapper`)
			}
		})
	}
}

func TestHandler(t *testing.T) {
	testCases := map[string]struct {
		in  slog.Handler
		out slog.Handler
	}{
		"wraps JSONHandler":            {slog.NewJSONHandler(os.Stderr, nil), slog.NewJSONHandler(os.Stderr, nil)},
		"wraps TextHandler":            {slog.NewTextHandler(os.Stderr, nil), slog.NewTextHandler(os.Stderr, nil)},
		"does not wrap ContextHandler": {NewContextHandler(slog.NewTextHandler(os.Stderr, nil)), slog.NewTextHandler(os.Stderr, nil)},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			v := NewContextHandler(test.in)

			if !reflect.DeepEqual(v.Handler(), test.out) {
				t.Fatalf(`wrong result for input "%v". Expected %v, got %v`, test.in, test.out, v)
			}
		})
	}
}
