package userip

import (
	"net"
	"net/http"
	"testing"
)

func TestFromRequest(t *testing.T) {
	tc := map[string]struct {
		header     string
		remoteAddr string
		err        string
	}{
		// Header
		"header valid IPv4": {
			header:     "127.0.0.1",
			remoteAddr: "",
			err:        "",
		},
		"header valid IPv4 list": {
			header:     "127.0.0.1, 127.0.0.2, 127.0.0.3,",
			remoteAddr: "",
			err:        "",
		},
		"header valid IPv4 list with space": {
			header:     "127.0.0.1  , 127.0.0.2 , 127.0.0.3,",
			remoteAddr: "",
			err:        "",
		},
		"header valid IPv6": {
			header:     "2001:db8::68",
			remoteAddr: "",
			err:        "",
		},
		"header valid ip IPv4-mapped IPv6": {
			header:     "::ffff:192.0.2.1",
			remoteAddr: "",
			err:        "",
		},
		"header invalid IPv4": {
			header:     "1270.0.0.1",
			remoteAddr: "",
			err:        `userip: "1270.0.0.1" comming from headers is not valid`,
		},
		"header invalid IPv6": {
			header:     "20010:db8::68",
			remoteAddr: "",
			err:        `userip: "20010:db8::68" comming from headers is not valid`,
		},
		"header invalid ip IPv4-mapped IPv6": {
			header:     "::ffff:1920.0.2.1",
			remoteAddr: "",
			err:        `userip: "::ffff:1920.0.2.1" comming from headers is not valid`,
		},
		// Remote addr
		"valid remote addr IPv4": {
			header:     "",
			remoteAddr: "127.0.0.1:80",
			err:        "",
		},
		"valid remote addr IPv6": {
			header:     "",
			remoteAddr: "[2001:db8::68]:80",
			err:        "",
		},
		"valid remote addr IPv4-mapped IPv6": {
			header:     "",
			remoteAddr: "[::ffff:192.0.2.1]:80",
			err:        "",
		},
		"invalid remote addr IPv4 - no port": {
			header:     "",
			remoteAddr: "127.0.0.1",
			err:        `userip: "127.0.0.1" is not IP:port`,
		},
		"invalid remote addr IPv6 - no port": {
			header:     "",
			remoteAddr: "2001:db8::68",
			err:        `userip: "2001:db8::68" is not IP:port`,
		},
		"invalid remote addr IPv4-mapped IPv6 - no port": {
			header:     "",
			remoteAddr: "::ffff:192.0.2.1",
			err:        `userip: "::ffff:192.0.2.1" is not IP:port`,
		},
		"invalid remote addr IPv4 - invalid ip": {
			header:     "",
			remoteAddr: "1270.0.0.1:80",
			err:        `userip: "1270.0.0.1" is not valid`,
		},
		"invalid remote addr IPv6 - invalid ip": {
			header:     "",
			remoteAddr: "[20010:db8::68]:80",
			err:        `userip: "20010:db8::68" is not valid`,
		},
		"invalid remote addr IPv4-mapped IPv6 - invalid ip": {
			header:     "",
			remoteAddr: "[::ffff:1920.0.2.1]:80",
			err:        `userip: "::ffff:1920.0.2.1" is not valid`,
		},
	}

	for name, test := range tc {
		t.Run(name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "", nil)
			req.Header.Add(headerName, test.header)
			req.RemoteAddr = test.remoteAddr

			_, err := FromRequest(req)
			if test.err != "" && (err == nil || err.Error() != test.err) {
				t.Fatalf(`unexpected error getting ip from request. Expected "%v", got "%v"`, test.err, err)
			}
		})
	}
}

func TestAddHeader(t *testing.T) {
	tc := map[string]struct {
		ip net.IP
	}{
		"valid IPv4":                {net.ParseIP("127.0.0.1")},
		"valid IPv6":                {net.ParseIP("2001:db8::68")},
		"valid ip IPv4-mapped IPv6": {net.ParseIP("::ffff:192.0.2.1")},
	}

	for name, test := range tc {
		t.Run(name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "", nil)

			AddHeader(req, test.ip)

			if h := req.Header.Get(headerName); h != test.ip.String() {
				t.Fatalf(`could not retrieved ip from header "%v". Expected "%v", got "%v"`, headerName, test.ip.String(), h)
			}

			ip, err := FromRequest(req)
			if err != nil {
				t.Fatalf(`unexpected error getting ip from request. Error: "%v"`, err)
			} else if !ip.Equal(test.ip) {
				t.Fatalf(`retrieved ip from request was not same as set with AddHeader. Expected "%v", got "%v"`, test.ip, ip)
			}
		})
	}
}
