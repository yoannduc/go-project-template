package httpmw

import (
	"net/http"
	"strings"
)

const (
	accessControlAllowOriginHeader  = "Access-Control-Allow-Origin"
	accessControlAllowMethodsHeader = "Access-Control-Allow-Methods"
	accessControlAllowHeadersHeader = "Access-Control-Allow-Headers"
	authorizationHeader             = "Authorization"
	accessControlAllowOriginValue   = "*"
)

// DevCORS is a middleware that sets dev appropriate cors headers
// then calls next.ServeHTTP(). This middleware is not production
// apporpriate and should be used only on local environments.
func DevCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(accessControlAllowOriginHeader, accessControlAllowOriginValue)
		w.Header().Set(accessControlAllowMethodsHeader, strings.Join([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions}, ", "))
		w.Header().Set(accessControlAllowHeadersHeader, strings.Join([]string{contentTypeHeader, authorizationHeader}, ", "))

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
