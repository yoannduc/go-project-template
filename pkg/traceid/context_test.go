package traceid

import (
	"context"
	"reflect"
	"testing"
)

func TestNewContext(t *testing.T) {
	tc := []struct {
		traceid string
	}{
		{traceid: ""},
		{traceid: "1234"},
		{traceid: "abcd"},
		{traceid: "ABCD"},
		{traceid: "#&!%$*,"},
	}

	for _, test := range tc {
		test := test
		t.Run(test.traceid, func(t *testing.T) {
			ctx := NewContext(context.Background(), test.traceid)

			if ctx.Value(traceIDKey) == nil {
				t.Fatalf("ctx is nil. Expected %v", test.traceid)
			}

			traceid, ok := ctx.Value(traceIDKey).(string)
			if !ok {
				t.Fatal("value held in ctx could not be cast to string")
			}

			if !reflect.DeepEqual(test.traceid, traceid) {
				t.Fatalf("session from ctx was not equal to test input. Expected %v, got %v", test.traceid, traceid)
			}
		})
	}
}

func TestFromContext(t *testing.T) {
	tc := []struct {
		traceid string
	}{
		{traceid: ""},
		{traceid: "1234"},
		{traceid: "abcd"},
		{traceid: "ABCD"},
		{traceid: "#&!%$*,"},
	}

	for _, test := range tc {
		test := test
		t.Run(test.traceid, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), traceIDKey, test.traceid)

			traceid, ok := FromContext(ctx)
			if !ok {
				t.Fatal("value held in ctx could not be cast to string")
			}

			if !reflect.DeepEqual(test.traceid, traceid) {
				t.Fatalf("session from ctx was not equal to test input. Expected %v, got %v", test.traceid, traceid)
			}
		})
	}
}
