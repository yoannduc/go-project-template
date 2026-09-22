package httpmw

import (
	"net/http"

	"github.com/yoannduc/go-project-template/pkg/traceid"
)

// LoadTraceID is a middleware that uses context/traceid to retrieve from
// headers or create new trace id. It then loads it to the context to be used
// anywhere this context is available.
func LoadTraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(traceid.NewContext(r.Context(), traceid.FromRequest(r))))
	})
}
