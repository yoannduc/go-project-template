package traceid

import (
	"net/http"
	"uuid"
)

// headerName is the name of the header for trace id.
const headerName = "X-Request-Id"

// FromRequest extracts the trace id from req headers, if present.
// It generates a new trace id if not present in req headers.
func FromRequest(req *http.Request) string {
	if v := req.Header.Get(headerName); v != "" {
		return v
	}

	return uuid.New().String()
}

// AddHeader adds the traceid in the headers of the http.Request.
func AddHeader(req *http.Request, tID string) {
	req.Header.Add(headerName, tID)
}
