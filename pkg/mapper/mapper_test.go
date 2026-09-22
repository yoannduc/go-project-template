package mapper

import (
	"net"
	"reflect"
	"testing"
)

type domain struct {
	Int     int
	Str     string
	Ptr     *string
	Any     any
	Private string
}

func (d domain) IsDomain() struct{} { return struct{}{} }

type dto struct {
	RenamedInt int
	RenamedStr string
	RenamedPtr *string
	RenamedAny any
}

func (d dto) FromDomain(dom domain) dto {
	d.RenamedInt = dom.Int
	d.RenamedStr = dom.Str
	d.RenamedPtr = dom.Ptr
	d.RenamedAny = dom.Any

	return d
}

func (d dto) ToDomain() domain {
	return domain{
		Int: d.RenamedInt,
		Str: d.RenamedStr,
		Ptr: d.RenamedPtr,
		Any: d.RenamedAny,
	}
}

type wrongDomain struct{}

func (d wrongDomain) IsDomain() struct{} { return struct{}{} }

type wrongDto struct{}

func (d wrongDto) FromDomain(dom wrongDomain) wrongDto {
	return d
}

func (d wrongDto) ToDomain() wrongDomain {
	return wrongDomain{}
}

func TestMapper(t *testing.T) {
	if !reflect.DeepEqual(New[domain, dto](), mapper[domain, dto]{}) {
		t.Fatalf(`mapper with correct types did not match between new and raw`)
	}

	if reflect.DeepEqual(New[wrongDomain, wrongDto](), mapper[domain, dto]{}) {
		t.Fatalf(`mapper with incorrect types should not match between new and raw`)
	}
}

func TestFromDomain(t *testing.T) {
	tmp := "ptr string"
	ptrValue := &tmp

	testCases := map[string]struct {
		in  domain
		out dto
	}{
		"empty":        {domain{}, dto{}},
		"private only": {domain{Private: "private"}, dto{}},
		"all fields": {
			domain{
				Int:     123,
				Str:     "string",
				Ptr:     ptrValue,
				Any:     net.ParseIP("8.8.8.8"),
				Private: "private",
			},
			dto{
				RenamedInt: 123,
				RenamedStr: "string",
				RenamedPtr: ptrValue,
				RenamedAny: net.ParseIP("8.8.8.8"),
			},
		},
		"only one exported field": {
			domain{
				Any: net.ParseIP("8.8.8.8"),
			},
			dto{
				RenamedAny: net.ParseIP("8.8.8.8"),
			},
		},
		"half fields": {
			domain{
				Ptr: ptrValue,
				Any: net.ParseIP("8.8.8.8"),
			},
			dto{
				RenamedPtr: ptrValue,
				RenamedAny: net.ParseIP("8.8.8.8"),
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			mpr := New[domain, dto]()
			v := mpr.FromDomain(test.in)

			if !reflect.DeepEqual(v, test.out) {
				t.Fatalf(`wrong type equality for input "%v". Expected %v, got %v`, test.in, test.out, v)
			}
		})
	}
}

func TestToDomain(t *testing.T) {
	tmp := "ptr string"
	ptrValue := &tmp

	testCases := map[string]struct {
		in  dto
		out domain
	}{
		"empty": {dto{}, domain{}},
		"all fields": {
			dto{
				RenamedInt: 123,
				RenamedStr: "string",
				RenamedPtr: ptrValue,
				RenamedAny: net.ParseIP("8.8.8.8"),
			},
			domain{
				Int: 123,
				Str: "string",
				Ptr: ptrValue,
				Any: net.ParseIP("8.8.8.8"),
			},
		},
		"only one exported field": {
			dto{
				RenamedAny: net.ParseIP("8.8.8.8"),
			},
			domain{
				Any: net.ParseIP("8.8.8.8"),
			},
		},
		"half fields": {
			dto{
				RenamedPtr: ptrValue,
				RenamedAny: net.ParseIP("8.8.8.8"),
			},
			domain{
				Ptr: ptrValue,
				Any: net.ParseIP("8.8.8.8"),
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			mpr := New[domain, dto]()
			v := mpr.ToDomain(test.in)

			if !reflect.DeepEqual(v, test.out) {
				t.Fatalf(`wrong type equality for input "%v". Expected %v, got %v`, test.in, test.out, v)
			}
		})
	}
}

func TestFromDomainList(t *testing.T) {
	tmp := "ptr string"
	ptrValue := &tmp

	testCases := map[string]struct {
		in  []domain
		out []dto
	}{
		"empty": {[]domain{}, []dto{}},
		"only one with only private": {
			[]domain{
				{Private: "private"},
			},
			[]dto{
				{},
			},
		},
		"different cases": {
			[]domain{
				{
					Int:     123,
					Str:     "string",
					Ptr:     ptrValue,
					Any:     net.ParseIP("8.8.8.8"),
					Private: "private",
				},
				{
					Any: net.ParseIP("8.8.8.8"),
				},
				{
					Ptr: ptrValue,
					Any: net.ParseIP("8.8.8.8"),
				},
			},
			[]dto{
				{
					RenamedInt: 123,
					RenamedStr: "string",
					RenamedPtr: ptrValue,
					RenamedAny: net.ParseIP("8.8.8.8"),
				},
				{
					RenamedAny: net.ParseIP("8.8.8.8"),
				},
				{
					RenamedPtr: ptrValue,
					RenamedAny: net.ParseIP("8.8.8.8"),
				},
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			mpr := New[domain, dto]()
			v := mpr.FromDomainList(test.in)

			if !reflect.DeepEqual(v, test.out) {
				t.Fatalf(`wrong type equality for input "%v". Expected %v, got %v`, test.in, test.out, v)
			}
		})
	}
}

func TestToDomainList(t *testing.T) {
	tmp := "ptr string"
	ptrValue := &tmp

	testCases := map[string]struct {
		in  []dto
		out []domain
	}{
		"empty": {[]dto{}, []domain{}},
		"different cases": {
			[]dto{
				{
					RenamedInt: 123,
					RenamedStr: "string",
					RenamedPtr: ptrValue,
					RenamedAny: net.ParseIP("8.8.8.8"),
				},
				{
					RenamedAny: net.ParseIP("8.8.8.8"),
				},
				{
					RenamedPtr: ptrValue,
					RenamedAny: net.ParseIP("8.8.8.8"),
				},
			},
			[]domain{
				{
					Int: 123,
					Str: "string",
					Ptr: ptrValue,
					Any: net.ParseIP("8.8.8.8"),
				},
				{
					Any: net.ParseIP("8.8.8.8"),
				},
				{
					Ptr: ptrValue,
					Any: net.ParseIP("8.8.8.8"),
				},
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			mpr := New[domain, dto]()
			v := mpr.ToDomainList(test.in)

			if !reflect.DeepEqual(v, test.out) {
				t.Fatalf(`wrong type equality for input "%v". Expected %v, got %v`, test.in, test.out, v)
			}
		})
	}
}
