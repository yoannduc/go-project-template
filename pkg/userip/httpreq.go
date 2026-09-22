package userip

import (
	"errors"
	"net"
	"net/http"
	"strings"
)

// headerName is the name of the header for user ip.
const headerName = "X-Forwarded-For"

// func

// FromRequest tries to extract the user IP address from req headers, if
// present. If it was not present, it extracts it from req if present.
// It errors if not present from either headers or request.
func FromRequest(req *http.Request) (net.IP, error) {
	ip, _, _ := strings.Cut(req.Header.Get(headerName), ",")
	if v := net.ParseIP(strings.Trim(ip, " ")); v != nil {
		return v, nil
	} else if ip != "" {
		return nil, errors.New("userip: \"" + ip + "\" comming from headers is not valid")
	}

	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, errors.New("userip: \"" + req.RemoteAddr + "\" is not IP:port")
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return nil, errors.New("userip: \"" + ip + "\" is not valid")
	}
	return userIP, nil
}

// AddHeader adds the user IP in the headers of the http.Request.
func AddHeader(req *http.Request, ip net.IP) {
	req.Header.Add(headerName, ip.String())
}
