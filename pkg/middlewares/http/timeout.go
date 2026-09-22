package httpmw

import (
	"net/http"
	"time"
)

const (
	// timeoutResponseBody is a string representing the body that is sent when a
	// timeout occurs. It can be a stringified json, stringified html...
	timeoutResponseBody = "request timed out"
)

// Timeout is a middleware that is a proxy to http.TimeoutHandler(). It only
// changes the signature function, and can now be pre loaded as
// `mw := Timeout(1*time.Second)` and then used and reused as `mw(handler)`. It
// also hides timeout message personnalisation to have a coherent message across
// all micro services that uses that middleware.
func Timeout(timeout time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, timeout, timeoutResponseBody)
	}
}
