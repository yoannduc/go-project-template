package httpmw

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/yoannduc/go-project-template/pkg/logger"
)

// httpResponseWriterMetadataRecorder encapsulates a http.ResponseWriter and
// wraps some methods to record informations that can not be retrieved on the
// writer otherwise. This includes total bytes written to the response body and
// status sent in the response.
type httpResponseWriterMetadataRecorder struct {
	http.ResponseWriter
	statusCode int
	bodySize   int
}

// WriteHeader uses underlying http.ResponseWriter WriteHeader() method. It
// wraps it in its own method to record the status written to log it.
func (rec *httpResponseWriterMetadataRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

// Write uses underlying http.ResponseWriter Write() method. It wraps it in its
// own method to count the total bytes written to log it.
func (rec *httpResponseWriterMetadataRecorder) Write(b []byte) (int, error) {
	n, err := rec.ResponseWriter.Write(b)
	rec.bodySize += n
	return n, err
}

// LogResponse is a function that takes in a *slog.Logger to have it
// preloaded and not regenerated each middleware call and returns
// a middleware that calls next.ServeHTTP() to let handler work and uses
// inputed logger after the handler returns to log infos on the request.
func LogResponse(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rec := httpResponseWriterMetadataRecorder{w, http.StatusOK, 0}
			next.ServeHTTP(&rec, r)

			log.LogAttrs(r.Context(), logger.LevelInfo, "request",
				slog.String("handler", "http"),
				slog.String("latency", time.Since(start).String()),
				slog.Int("statusCode", rec.statusCode),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("bodySize", rec.bodySize),
			)
		})
	}
}
