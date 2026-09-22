package traceid

import (
	"net/http"
	"testing"
	"uuid"
)

func TestFromRequest(t *testing.T) {
	tc := []struct {
		traceid string
	}{
		{traceid: ""},
		// Random valid uuid to put in header & check we get the same one out
		{traceid: "7d03a148-844d-4e87-b9cf-1c348c4e9b50"},
		{traceid: "c9d482c1-81f9-466f-9a98-2299c646d35c"},
		{traceid: "4abadaaf-b9e4-4669-98db-12b43c3afecf"},
	}

	for _, test := range tc {
		t.Run(test.traceid, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "", nil)
			req.Header.Add(headerName, test.traceid)

			traceid := FromRequest(req)

			if test.traceid != "" && test.traceid != traceid {
				t.Fatalf("traceid is not same from headers. Expected %v, got %v", test.traceid, traceid)
			}

			_, err := uuid.Parse(traceid)
			if err != nil {
				t.Fatalf("error parsing resulting uuid. Err: %v", err)
			}
		})
	}
}

func TestAddHeader(t *testing.T) {
	tc := []struct {
		traceid string
	}{
		{traceid: ""},
		// Random valid uuid to put in header & check we get the same one out
		{traceid: "7d03a148-844d-4e87-b9cf-1c348c4e9b50"},
	}

	for _, test := range tc {
		t.Run(test.traceid, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "", nil)

			AddHeader(req, test.traceid)

			if h := req.Header.Get(headerName); h != test.traceid {
				t.Fatalf(`could not retrieved traceid from header "%v". Expected "%v", got "%v"`, headerName, test.traceid, h)
			}

			if traceid := FromRequest(req); test.traceid != "" && traceid != test.traceid {
				t.Fatalf(`retrieved traceId from request was not same as set with AddHeader. Expected "%v", got "%v"`, test.traceid, traceid)
			}
		})
	}
}
