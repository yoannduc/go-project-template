package singleton

import "testing"

func TestSingleton(t *testing.T) {
	testCases := map[string]struct {
		in any
	}{
		"1234":    {1234},
		"string":  {"string"},
		"pointer": {new(1234)},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			v := New(test.in)
			if v.GetValue() != test.in {
				t.Fatalf(`wrong result for singleton value "%v". Expected %v, got %v`, test.in, test.in, v.GetValue())
			}
		})
	}
}
