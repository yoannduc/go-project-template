package userip

import (
	"context"
	"net"
	"reflect"
	"testing"
)

func TestNewContext(t *testing.T) {
	tc := []struct {
		userip net.IP
	}{
		{userip: net.IP{}},
		{userip: net.ParseIP("8.8.8.8")},
		{userip: net.ParseIP("127.0.0.1")},
		// IPv6
		{userip: net.ParseIP("2001:db8::68")},
		// IPv4-mapped IPv6
		{userip: net.ParseIP("::ffff:192.0.2.1")},
	}

	for _, test := range tc {
		t.Run(test.userip.String(), func(t *testing.T) {
			ctx := NewContext(context.Background(), test.userip)

			if ctx.Value(userIPKey) == nil {
				t.Fatalf("ctx is nil. Expected %v", test.userip.String())
			}

			userip, ok := ctx.Value(userIPKey).(net.IP)
			if !ok {
				t.Fatal("value held in ctx could not be cast to net.IP")
			}

			if !reflect.DeepEqual(test.userip, userip) {
				t.Fatalf("value from ctx was not equal to test input. Expected %v, got %v", test.userip.String(), userip.String())
			}
		})
	}
}

func TestFromContext(t *testing.T) {
	tc := []struct {
		userip net.IP
	}{
		{userip: net.IP{}},
		{userip: net.ParseIP("8.8.8.8")},
		{userip: net.ParseIP("127.0.0.1")},
		// IPv6
		{userip: net.ParseIP("2001:db8::68")},
		// IPv4-mapped IPv6
		{userip: net.ParseIP("::ffff:192.0.2.1")},
	}

	for _, test := range tc {
		t.Run(test.userip.String(), func(t *testing.T) {
			ctx := context.WithValue(context.Background(), userIPKey, test.userip)

			userip, ok := FromContext(ctx)
			if !ok {
				t.Fatal("value held in ctx could not be cast to net.IP")
			}

			if !reflect.DeepEqual(test.userip, userip) {
				t.Fatalf("value from ctx was not equal to test input. Expected %v, got %v", test.userip.String(), userip.String())
			}
		})
	}
}
