package ports

import "net/http"

// Handler is the api contract for a handler.
type Handler interface {
	LoadRoutes(mux *http.ServeMux)
}
