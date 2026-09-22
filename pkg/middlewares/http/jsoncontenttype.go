package httpmw

import "net/http"

const (
	contentTypeHeader    = "Content-Type"
	applicationJSONValue = "application/json"
)

// writerWrapper encapsulates a http.ResponseWriter and
// wraps some methods to record informations that can not be retrieved on the
// writer otherwise. This includes total bytes written to the response body and
// status sent in the response.
type writerWrapper struct {
	http.ResponseWriter
}

// WriteHeader uses underlying http.ResponseWriter WriteHeader() method. It
// wraps it in its own method to add application/json header by default on
// success.
func (rec *writerWrapper) WriteHeader(code int) {
	if rec.ResponseWriter.Header().Get(contentTypeHeader) == "" && code >= http.StatusOK && code < http.StatusMultipleChoices {
		rec.ResponseWriter.Header().Set(contentTypeHeader, applicationJSONValue)
	}
	rec.ResponseWriter.WriteHeader(code)
}

// SuccessJSONContentType is a middleware that calls next.ServeHTTP() to let
// handler work and uses writerWrapper to add application/json content type
// header on success, if no content type was specified.
func SuccessJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrap := writerWrapper{w}
		next.ServeHTTP(&wrap, r)
	})
}
