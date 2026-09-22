package memorydb

import (
	"context"
	"errors"
	"net"
	"reflect"
	"slices"
	"strings"
	"testing"
)

type tone struct {
	ID  int
	Str string
	Any any
	Ptr *string
}

var cmpNoCompareFunc Compare[tone] = func(a, b tone) bool {
	return false
}

var cmpStrFunc Compare[tone] = func(a, b tone) bool {
	if a.ID == b.ID {
		return false
	}

	return strings.EqualFold(a.Str, b.Str)
}

type ttwo struct{}

var cmpNoCompareFuncTwo Compare[ttwo] = func(a, b ttwo) bool {
	return false
}

func TestUpdateIDField(t *testing.T) {
	t.Run("has ID field", func(t *testing.T) {
		testCases := map[string]struct {
			el  tone
			id  uint64
			out tone
			err error
		}{
			"empty": {
				tone{},
				1,
				tone{ID: 1},
				nil,
			},
			"overrides id": {
				tone{ID: 10},
				1,
				tone{ID: 1},
				nil,
			},
		}

		for name, test := range testCases {
			t.Run(name, func(t *testing.T) {
				v, err := updateIDField(test.el, test.id)
				if !reflect.DeepEqual(v, test.out) {
					t.Fatalf(`wrong type equality. Expected %v, got %v`, test.out, v)
				}
				if !errors.Is(err, test.err) {
					t.Fatalf(`unexpected error. Expected "%v", got "%v"`, test.err, err)
				}
			})
		}
	})

	t.Run("has no ID field", func(t *testing.T) {
		testCases := map[string]struct {
			el  ttwo
			id  uint64
			out ttwo
			err error
		}{
			"err": {
				ttwo{},
				1,
				ttwo{},
				errNoIDField,
			},
		}

		for name, test := range testCases {
			t.Run(name, func(t *testing.T) {
				v, err := updateIDField(test.el, test.id)
				if !reflect.DeepEqual(v, test.out) {
					t.Fatalf(`wrong type equality. Expected %v, got %v`, test.out, v)
				}
				if !errors.Is(err, test.err) {
					t.Fatalf(`unexpected error. Expected "%v", got "%v"`, test.err, err)
				}
			})
		}
	})

	t.Run("is not struct", func(t *testing.T) {
		testCases := map[string]struct {
			el  int
			id  uint64
			out int
			err error
		}{
			"err": {
				1,
				1,
				1,
				errNoIDField,
			},
		}

		for name, test := range testCases {
			t.Run(name, func(t *testing.T) {
				v, err := updateIDField(test.el, test.id)
				if !reflect.DeepEqual(v, test.out) {
					t.Fatalf(`wrong type equality. Expected %v, got %v`, test.out, v)
				}
				if !errors.Is(err, test.err) {
					t.Fatalf(`unexpected error. Expected "%v", got "%v"`, test.err, err)
				}
			})
		}
	})
}

func TestExistsFunc(t *testing.T) {
	testCases := map[string]struct {
		slice []tone
		el    tone
		cmp   Compare[tone]
		out   bool
	}{
		"empty": {
			[]tone{},
			tone{},
			func(a, b tone) bool { return true },
			false,
		},
		"does not find": {
			[]tone{
				{
					ID:  1,
					Str: "toto",
				},
				{
					ID:  2,
					Str: "titi",
				},
				{
					ID:  3,
					Str: "tutu",
				},
				{
					ID:  4,
					Str: "tata",
				},
			},
			tone{
				ID:  4,
				Str: "toto",
			},
			func(a, b tone) bool { return a.ID == b.ID && strings.EqualFold(a.Str, b.Str) },
			false,
		},
		"found id": {
			[]tone{
				{
					ID:  1,
					Str: "toto",
				},
				{
					ID:  2,
					Str: "titi",
				},
				{
					ID:  3,
					Str: "tutu",
				},
				{
					ID:  4,
					Str: "tata",
				},
			},
			tone{
				ID:  4,
				Str: "toto",
			},
			func(a, b tone) bool { return a.ID == b.ID },
			true,
		},
		"found string": {
			[]tone{
				{
					ID:  1,
					Str: "toto",
				},
				{
					ID:  2,
					Str: "titi",
				},
				{
					ID:  3,
					Str: "tutu",
				},
				{
					ID:  4,
					Str: "tata",
				},
			},
			tone{
				ID:  4,
				Str: "toto",
			},
			func(a, b tone) bool { return strings.EqualFold(a.Str, b.Str) },
			true,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			v := existsFunc(slices.Values(test.slice), test.el, test.cmp)
			if !reflect.DeepEqual(v, test.out) {
				t.Fatalf(`wrong type equality. Expected %v, got %v`, test.out, v)
			}
		})
	}
}

func TestFilterFunc(t *testing.T) {
	testCases := map[string]struct {
		slice []tone
		fltr  Filter[tone]
		out   []tone
	}{
		"empty": {
			[]tone{},
			func(t tone) bool { return true },
			nil,
		},
		"filter one field": {
			[]tone{
				{
					ID:  1,
					Str: "toto",
				},
				{
					ID:  2,
					Str: "titi",
				},
				{
					ID:  3,
					Str: "tutu",
				},
				{
					ID:  4,
					Str: "tata",
				},
			},
			func(t tone) bool { return t.ID < 3 },
			[]tone{
				{
					ID:  1,
					Str: "toto",
				},
				{
					ID:  2,
					Str: "titi",
				},
			},
		},
		"filter two field": {
			[]tone{
				{
					ID:  1,
					Str: "that",
				},
				{
					ID:  2,
					Str: "this",
				},
				{
					ID:  3,
					Str: "that",
				},
				{
					ID:  4,
					Str: "this",
				},
				{
					ID:  5,
					Str: "this",
				},
				{
					ID:  6,
					Str: "that",
				},
			},
			func(t tone) bool { return t.ID%2 == 0 && strings.EqualFold(t.Str, "this") },
			[]tone{
				{
					ID:  2,
					Str: "this",
				},
				{
					ID:  4,
					Str: "this",
				},
			},
		},
		"no filter": {
			[]tone{
				{
					ID:  1,
					Str: "toto",
				},
				{
					ID:  2,
					Str: "titi",
				},
				{
					ID:  3,
					Str: "tutu",
				},
				{
					ID:  4,
					Str: "tata",
				},
			},
			func(t tone) bool { return true },
			[]tone{
				{
					ID:  1,
					Str: "toto",
				},
				{
					ID:  2,
					Str: "titi",
				},
				{
					ID:  3,
					Str: "tutu",
				},
				{
					ID:  4,
					Str: "tata",
				},
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			it := filterFunc(slices.Values(test.slice), test.fltr)
			v := slices.Collect(it)
			if !reflect.DeepEqual(v, test.out) {
				t.Fatalf(`wrong type equality. Expected %v, got %v`, test.out, v)
			}
		})
	}
}

func TestNew(t *testing.T) {
	out := &memDB[tone]{
		m: make(map[uint64]tone),
	}
	out.index.Add(1)

	if !reflect.DeepEqual(New[tone](), out) {
		t.Fatalf(`memDB with correct types did not match between new and raw`)
	}

	if reflect.DeepEqual(New[ttwo](), out) {
		t.Fatalf(`memDB with incorrect types should not match between new and raw`)
	}
}

func TestFindAll(t *testing.T) {
	tmp := "ptr string"
	ptrValue := &tmp

	testCases := map[string]struct {
		out []tone
	}{
		"empty": {nil},
		"one": {
			[]tone{
				{
					ID:  1,
					Str: "string",
				},
			},
		},
		"some": {
			[]tone{
				{
					ID:  1,
					Str: "string",
				},
				{
					ID:  2,
					Ptr: ptrValue,
				},
				{
					ID:  3,
					Any: net.ParseIP("8.8.8.8"),
				},
				{
					ID:  4,
					Str: "string",
					Ptr: ptrValue,
					Any: net.ParseIP("8.8.8.8"),
				},
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			db := New[tone]()

			ctx := context.Background()
			for _, v := range test.out {
				db.Create(ctx, v, cmpNoCompareFunc)
			}

			v, _ := db.FindAll(ctx)
			// Order return for deepequal compare. Order is not actually
			// guarenteed with findAll, it is caller's role to order.
			slices.SortFunc(v, func(a, b tone) int { return a.ID - b.ID })
			if !reflect.DeepEqual(v, test.out) {
				t.Fatalf(`wrong type equality for input "%v". Expected %v, got %v`, test.out, test.out, v)
			}
		})
	}
}

func TestFindFilterFunc(t *testing.T) {
	type findFilterFunc struct {
		filter Filter[tone]
		out    []tone
	}

	tmp := "ptr string"
	ptrValue := &tmp

	testCases := map[string]struct {
		toCreate []tone
		find     []findFilterFunc
	}{
		// empty not equals, need to investigate
		"empty": {
			[]tone{},
			[]findFilterFunc{
				{
					func(tone tone) bool {
						return true
					},
					nil,
				},
				{
					func(tone tone) bool {
						return strings.Contains(tone.Str, "ok")
					},
					nil,
				},
			},
		},
		"some": {
			[]tone{
				{
					Str: "one",
				},
				{
					Str: "two",
					Ptr: ptrValue,
				},
				{
					Str: "three",
					Any: net.ParseIP("8.8.8.8"),
				},
				{
					Str: "four",
					Ptr: ptrValue,
					Any: net.ParseIP("8.8.8.8"),
				},
			},
			[]findFilterFunc{
				{
					func(tone tone) bool {
						return true
					},
					[]tone{
						{
							ID:  1,
							Str: "one",
						},
						{
							ID:  2,
							Str: "two",
							Ptr: ptrValue,
						},
						{
							ID:  3,
							Str: "three",
							Any: net.ParseIP("8.8.8.8"),
						},
						{
							ID:  4,
							Str: "four",
							Ptr: ptrValue,
							Any: net.ParseIP("8.8.8.8"),
						},
					},
				},
				{
					func(tone tone) bool {
						return strings.Contains(tone.Str, "t")
					},
					[]tone{
						{
							ID:  2,
							Str: "two",
							Ptr: ptrValue,
						},
						{
							ID:  3,
							Str: "three",
							Any: net.ParseIP("8.8.8.8"),
						},
					},
				},
				{
					func(tone tone) bool {
						return tone.Ptr == ptrValue
					},
					[]tone{
						{
							ID:  2,
							Str: "two",
							Ptr: ptrValue,
						},
						{
							ID:  4,
							Str: "four",
							Ptr: ptrValue,
							Any: net.ParseIP("8.8.8.8"),
						},
					},
				},
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			db := New[tone]()
			ctx := context.Background()

			for _, v := range test.toCreate {
				db.Create(ctx, v, cmpNoCompareFunc)
			}

			for _, tst := range test.find {
				v, _ := db.FindFilterFunc(ctx, tst.filter)
				slices.SortFunc(v, func(a, b tone) int { return a.ID - b.ID })
				if !reflect.DeepEqual(v, tst.out) {
					t.Fatalf(`wrong type equality for input "%v". Expected %v, got %v`, tst.filter, tst.out, v)
				}
			}
		})
	}
}

func TestFindByID(t *testing.T) {
	type find struct {
		ID  uint64
		out tone
		err error
	}

	tmp := "ptr string"
	ptrValue := &tmp

	testCases := map[string]struct {
		toCreate []tone
		find     []find
	}{
		"empty": {[]tone{}, nil},
		"did not exist": {
			[]tone{},
			[]find{
				{
					1,
					tone{},
					errNotFound,
				},
			},
		},
		"some": {
			[]tone{
				{
					Str: "one",
				},
				{
					Str: "two",
					Ptr: ptrValue,
				},
				{
					Str: "three",
					Any: net.ParseIP("8.8.8.8"),
				},
				{
					Str: "four",
					Ptr: ptrValue,
					Any: net.ParseIP("8.8.8.8"),
				},
			},
			[]find{
				{
					1,
					tone{
						ID:  1,
						Str: "one",
					},
					nil,
				},
				{
					4,
					tone{
						ID:  4,
						Str: "four",
						Ptr: ptrValue,
						Any: net.ParseIP("8.8.8.8"),
					},
					nil,
				},
				{
					2,
					tone{
						ID:  2,
						Str: "two",
						Ptr: ptrValue,
					},
					nil,
				},
				{
					3,
					tone{
						ID:  3,
						Str: "three",
						Any: net.ParseIP("8.8.8.8"),
					},
					nil,
				},
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			db := New[tone]()
			ctx := context.Background()

			for _, v := range test.toCreate {
				db.Create(ctx, v, cmpStrFunc)
			}

			for _, tst := range test.find {
				v, err := db.FindByID(ctx, tst.ID)
				if err == nil && !reflect.DeepEqual(v, tst.out) {
					t.Fatalf(`wrong type equality for input "%v". Expected %v, got %v`, tst.ID, tst.out, v)
				}
				if err != nil && !errors.Is(err, tst.err) {
					t.Fatalf(`unexpected error. Expected "%v", got "%v"`, tst.err, err)
				}
			}
		})
	}
}

func TestCreate(t *testing.T) {
	type inout[T tone | ttwo] struct {
		in  T
		out T
		err error
	}

	tmp := "ptr string"
	ptrValue := &tmp

	t.Run("no errors, many types", func(t *testing.T) {
		testCases := []inout[tone]{
			{
				tone{
					Str: "string",
				},
				tone{
					ID:  1,
					Str: "string",
				},
				nil,
			},
			{
				tone{
					Ptr: ptrValue,
				},
				tone{
					ID:  2,
					Ptr: ptrValue,
				},
				nil,
			},
			{
				tone{
					Any: net.ParseIP("8.8.8.8"),
				},
				tone{
					ID:  3,
					Any: net.ParseIP("8.8.8.8"),
				},
				nil,
			},
			{
				tone{
					Str: "string",
					Ptr: ptrValue,
					Any: net.ParseIP("8.8.8.8"),
				},
				tone{
					ID:  4,
					Str: "string",
					Ptr: ptrValue,
					Any: net.ParseIP("8.8.8.8"),
				},
				nil,
			},
		}

		db := New[tone]()
		ctx := context.Background()

		for _, test := range testCases {
			v, err := db.Create(ctx, test.in, cmpNoCompareFunc)
			if err == nil && !reflect.DeepEqual(v, test.out) {
				t.Fatalf(`wrong type equality for input "%v". Expected %v, got %v`, test.in, test.out, v)
			}
			if err != nil && !errors.Is(err, test.err) {
				t.Fatalf(`unexpected error. Expected "%v", got "%v"`, test.err, err)
			}
		}
	})

	t.Run("no error, does not take in inputed ID", func(t *testing.T) {
		testCases := []inout[tone]{
			{
				tone{
					ID:  155,
					Str: "first",
				},
				tone{
					ID:  1,
					Str: "first",
				},
				nil,
			},
			{
				tone{
					ID:  -5000,
					Str: "second",
				},
				tone{
					ID:  2,
					Str: "second",
				},
				nil,
			},
			{
				tone{
					ID:  1,
					Str: "third",
				},
				tone{
					ID:  3,
					Str: "third",
				},
				nil,
			},
		}

		db := New[tone]()
		ctx := context.Background()

		for _, test := range testCases {
			v, err := db.Create(ctx, test.in, cmpStrFunc)
			if err != nil && !errors.Is(err, test.err) {
				t.Fatalf(`unexpected error. Expected "%v", got "%v"`, test.err, err)
			}
			if !reflect.DeepEqual(v, test.out) {
				t.Fatalf(`wrong type equality for input "%v". Expected %v, got %v`, test.in, test.out, v)
			}
		}
	})

	t.Run("err already exists", func(t *testing.T) {
		testCases := []inout[tone]{
			{
				tone{
					Str: "string",
				},
				tone{
					ID:  1,
					Str: "string",
				},
				nil,
			},
			{
				tone{
					Str: "string",
				},
				tone{
					Str: "string",
				},
				errAlreadyExists,
			},
			{
				tone{
					Ptr: ptrValue,
				},
				tone{
					ID:  2,
					Ptr: ptrValue,
				},
				nil,
			},
			{
				tone{
					Ptr: ptrValue,
				},
				tone{
					Ptr: ptrValue,
				},
				errAlreadyExists,
			},
			{
				tone{
					Any: net.ParseIP("8.8.8.8"),
				},
				tone{
					Any: net.ParseIP("8.8.8.8"),
				},
				errAlreadyExists,
			},
			{
				tone{
					Str: "string",
					Ptr: ptrValue,
					Any: net.ParseIP("8.8.8.8"),
				},
				tone{
					Str: "string",
					Ptr: ptrValue,
					Any: net.ParseIP("8.8.8.8"),
				},
				errAlreadyExists,
			},
		}

		db := New[tone]()
		ctx := context.Background()

		for _, test := range testCases {
			v, err := db.Create(ctx, test.in, cmpStrFunc)
			if err != nil && !errors.Is(err, test.err) {
				t.Fatalf(`unexpected error. Expected "%v", got "%v"`, test.err, err)
			}
			if !reflect.DeepEqual(v, test.out) {
				t.Fatalf(`wrong type equality for input "%v". Expected %v, got %v`, test.in, test.out, v)
			}
		}
	})

	t.Run("err no ID field", func(t *testing.T) {
		testCases := []inout[ttwo]{
			{
				ttwo{},
				ttwo{},
				errNoIDField,
			},
		}

		db := New[ttwo]()
		ctx := context.Background()

		for _, test := range testCases {
			v, err := db.Create(ctx, test.in, cmpNoCompareFuncTwo)
			if err == nil && !reflect.DeepEqual(v, test.out) {
				t.Fatalf(`wrong type equality for input "%v". Expected %v, got %v`, test.in, test.out, v)
			}
			if !errors.Is(err, test.err) {
				t.Fatalf(`unexpected error. Expected "%v", got "%v"`, test.err, err)
			}
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("case with ID field", func(t *testing.T) {
		type update struct {
			ID  uint64
			in  tone
			out tone
			err error
		}

		tmp := "ptr string"
		ptrValue := &tmp

		testCases := map[string]struct {
			toCreate []tone
			update   []update
		}{
			"err already exists but not for itself (depending on compare func)": {
				[]tone{
					{
						Str: "first",
					},
					{
						Str: "second",
					},
				},
				[]update{
					{
						1,
						tone{
							Str: "second",
							Any: "new",
						},
						tone{
							Str: "second",
							Any: "new",
						},
						errAlreadyExists,
					},
					{
						2,
						tone{
							Str: "second",
							Any: "new",
						},
						tone{
							ID:  2,
							Str: "second",
							Any: "new",
						},
						nil,
					},
				},
			},
			"err not found": {
				[]tone{},
				[]update{
					{
						1,
						tone{},
						tone{},
						errNotFound,
					},
				},
			},
			"does not update ID": {
				[]tone{
					{
						Str: "string",
					},
				},
				[]update{
					{
						1,
						tone{
							ID:  1000,
							Str: "str",
						},
						tone{
							ID:  1,
							Str: "str",
						},
						nil,
					},
				},
			},
			"some": {
				[]tone{
					{
						Str: "one",
					},
					{
						Str: "two",
						Ptr: ptrValue,
					},
					{
						Str: "three",
						Any: net.ParseIP("8.8.8.8"),
					},
					{
						Str: "four",
						Ptr: ptrValue,
						Any: net.ParseIP("8.8.8.8"),
					},
				},
				[]update{
					{
						1,
						tone{
							Str: "one",
							Any: net.ParseIP("8.8.8.8"),
						},
						tone{
							ID:  1,
							Str: "one",
							Any: net.ParseIP("8.8.8.8"),
						},
						nil,
					},
					{
						1,
						tone{
							Str: "first",
							Any: net.ParseIP("8.8.8.8"),
						},
						tone{
							ID:  1,
							Str: "first",
							Any: net.ParseIP("8.8.8.8"),
						},
						nil,
					},

					{
						3,
						tone{
							Str: "three",
							Any: "~",
						},
						tone{
							ID:  3,
							Str: "three",
							Any: "~",
						},
						nil,
					},
					{
						2,
						tone{
							Str: "four",
						},
						tone{
							Str: "four",
						},
						errAlreadyExists,
					},
				},
			},
		}

		for name, test := range testCases {
			t.Run(name, func(t *testing.T) {
				db := New[tone]()
				ctx := context.Background()

				for _, v := range test.toCreate {
					db.Create(ctx, v, cmpNoCompareFunc)
				}

				for _, tst := range test.update {
					v, err := db.Update(ctx, tst.ID, tst.in, cmpStrFunc)
					if err != nil && !errors.Is(err, tst.err) {
						t.Fatalf(`unexpected error. Expected "%v", got "%v"`, tst.err, err)
					}
					if !reflect.DeepEqual(v, tst.out) {
						t.Fatalf(`wrong type equality for input "%v". Expected %v, got %v`, tst.ID, tst.out, v)
					}
				}
			})
		}
	})

	t.Run("", func(t *testing.T) {
		type update struct {
			ID  uint64
			in  ttwo
			out ttwo
			err error
		}

		testCases := map[string]struct {
			toCreate []ttwo
			update   []update
		}{
			"err not found": {
				[]ttwo{},
				[]update{
					{
						1,
						ttwo{},
						ttwo{},
						errNoIDField,
					},
				},
			},
		}

		for name, test := range testCases {
			t.Run(name, func(t *testing.T) {
				db := New[ttwo]()
				ctx := context.Background()

				for _, v := range test.toCreate {
					db.Create(ctx, v, cmpNoCompareFuncTwo)
				}

				for _, tst := range test.update {
					v, err := db.Update(ctx, tst.ID, tst.in, cmpNoCompareFuncTwo)
					if err != nil && !errors.Is(err, tst.err) {
						t.Fatalf(`unexpected error. Expected "%v", got "%v"`, tst.err, err)
					}
					if !reflect.DeepEqual(v, tst.out) {
						t.Fatalf(`wrong type equality for input "%v". Expected %v, got %v`, tst.ID, tst.out, v)
					}
				}
			})
		}
	})
}

func TestDelete(t *testing.T) {
	type delete struct {
		ID  uint64
		out tone
		err error
	}

	tmp := "ptr string"
	ptrValue := &tmp

	testCases := map[string]struct {
		toCreate []tone
		delete   []delete
	}{
		"err not found": {
			[]tone{},
			[]delete{
				{
					1,
					tone{},
					errNotFound,
				},
			},
		},
		"err not found after delete": {
			[]tone{
				{
					Str: "string",
				},
			},
			[]delete{
				{
					1,
					tone{
						ID:  1,
						Str: "string",
					},
					nil,
				},
				{
					1,
					tone{},
					errNotFound,
				},
			},
		},
		"some": {
			[]tone{
				{
					Str: "string",
				},
				{
					Ptr: ptrValue,
				},
				{
					Any: net.ParseIP("8.8.8.8"),
				},
				{
					Str: "string",
					Ptr: ptrValue,
					Any: net.ParseIP("8.8.8.8"),
				},
			},
			[]delete{
				{
					3,
					tone{
						ID:  3,
						Any: net.ParseIP("8.8.8.8"),
					},
					nil,
				},
				{
					1,
					tone{
						ID:  1,
						Str: "string",
					},
					nil,
				},
				{
					4,
					tone{
						ID:  4,
						Str: "string",
						Ptr: ptrValue,
						Any: net.ParseIP("8.8.8.8"),
					},
					nil,
				},
				{
					2,
					tone{
						ID:  2,
						Ptr: ptrValue,
					},
					nil,
				},
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			db := New[tone]()
			ctx := context.Background()

			for _, v := range test.toCreate {
				db.Create(ctx, v, cmpNoCompareFunc)
			}

			for _, tst := range test.delete {
				v, err := db.Delete(ctx, tst.ID)
				if err != nil && !errors.Is(err, tst.err) {
					t.Fatalf(`unexpected error. Expected "%v", got "%v"`, tst.err, err)
				}
				if !reflect.DeepEqual(v, tst.out) {
					t.Fatalf(`wrong type equality for input "%v". Expected %v, got %v`, tst.ID, tst.out, v)
				}
			}
		})
	}
}
