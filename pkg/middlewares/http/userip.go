package httpmw

import (
	"net/http"

	"github.com/yoannduc/go-project-template/pkg/userip"
)

// LoadUserIP is a middleware that uses context/userip to retrieve from
// headers (from proxy) or request caller ip. It then loads it to the context to
// be used anywhere this context is available.
func LoadUserIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ip, err := userip.FromRequest(r); err == nil {
			r = r.WithContext(userip.NewContext(r.Context(), ip))
		}

		next.ServeHTTP(w, r)
	})
}
