package memorydb

import (
	"context"
	"errors"
	"iter"
	"maps"
	"reflect"
	"slices"
	"sync"
	"sync/atomic"
)

var (
	errNotFound      = errors.New("not found")
	errAlreadyExists = errors.New("entry already exists")
	errNoIDField     = errors.New("object has no ID field, could not use in DB")
)

const (
	idFieldName = "ID"
)

// updateIDField updates the ID field of in. It uses reflect,
// as in is of type any to accomodate any struct with an ID field.
func updateIDField[T any](in T, id uint64) (T, error) {
	str := reflect.ValueOf(&in).Elem()
	if str.Kind() != reflect.Struct {
		return in, errNoIDField
	}
	f := str.FieldByName(idFieldName)
	if !f.IsValid() ||
		(f.Kind() != reflect.Int && f.Kind() != reflect.Int64 && f.Kind() != reflect.Int32) {
		return in, errNoIDField
	}
	f.SetInt(int64(id))

	return in, nil
}

// Compare is the signature for the compare function applied on
// searches. First argument is current element in loop, second
// is compare compare element.
type Compare[T any] func(T, T) bool

// existsFunc searches whether any element in iter.Seq matches el
// when compared using cmp.
func existsFunc[T any](it iter.Seq[T], el T, cmp Compare[T]) bool {
	for v := range it {
		if cmp(v, el) {
			return true
		}
	}

	return false
}

// Filter is the signature for the filter function applied to element list.
type Filter[T any] func(T) bool

// filterFunc filters iter.Seq via f.
func filterFunc[T any](it iter.Seq[T], f Filter[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range it {
			if f(v) && !yield(v) {
				return
			}
		}
	}
}

// A MemoryDB is a concurrent safe very minimal in memory db.
// It needs the object to be stored to have an ID field of type int,
// such as:
//
//	type MyType struct {
//		ID: int, // <- here
//		Key: any,
//		...,
//	}
//
// MemoryDB is very much not intended to be used in production,
// it is intended to be used for MVPs or small personal projects.
type MemoryDB[T any] interface {
	FindAll(context.Context) ([]T, error)
	FindFilterFunc(context.Context, Filter[T]) ([]T, error)
	FindByID(context.Context, uint64) (T, error)
	Create(context.Context, T, Compare[T]) (T, error)
	Update(_ context.Context, id uint64, in T, fnd Compare[T]) (T, error)
	Delete(context.Context, uint64) (T, error)
}

// memDB is the concrete type implementing MemoryDB.
type memDB[T any] struct {
	mu    sync.Mutex
	m     map[uint64]T
	index atomic.Uint64
}

// New returns a MemoryDB.
func New[T any]() MemoryDB[T] {
	db := &memDB[T]{
		m: make(map[uint64]T, 100),
	}
	db.index.Add(1)

	return db
}

// FindAll returns a slice of all element in internal map.
// Order is not guarenteed.
func (db *memDB[T]) FindAll(_ context.Context) ([]T, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	return slices.Collect(maps.Values(db.m)), nil
}

// FindFilterFunc returns a slice of filtered element in internal map.
// It takes in a closure to let caller define how results should be filtered.
// Order is not guarenteed.
func (db *memDB[T]) FindFilterFunc(_ context.Context, fltr Filter[T]) ([]T, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	return slices.Collect(filterFunc(maps.Values(db.m), fltr)), nil
}

// FindByID returns the element at id index.
// It returns an error if no element found.
func (db *memDB[T]) FindByID(_ context.Context, id uint64) (T, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	v, ok := db.m[id]
	if !ok {
		return v, errNotFound
	}

	return v, nil
}

// Create inserts in in the map. It performs a search based on cmp
// to let caller decide how uniqueness works for T. It returns an
// error if duplicate element was found or if it could not update
// ID field in element (see [MemoryDB] for more info on element
// restrictions). It returns the created element.
func (db *memDB[T]) Create(_ context.Context, in T, cmp Compare[T]) (T, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	idx := db.index.Load()
	// Update ID before search to let caller use ID as valid field.
	out, err := updateIDField(in, idx)
	if err != nil {
		return in, err
	}
	if existsFunc(maps.Values(db.m), out, cmp) {
		return in, errAlreadyExists
	}

	db.m[idx] = out
	db.index.Add(1)
	return out, nil
}

// Update updates in at id index. It updates the element in place,
// not updating each individual field, so the full updated object
// should be passed (ID is ignored). It performs a search based on
// cmp to let caller decide how uniqueness works for T. It returns
// an error if duplicate element was found or if it could not update
// ID field in element (see [MemoryDB] for more info on element
// restrictions). It returns the updated element.
func (db *memDB[T]) Update(_ context.Context, id uint64, in T, cmp Compare[T]) (T, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Update ID before search to let caller use ID as valid field.
	out, err := updateIDField(in, id)
	if err != nil {
		return in, err
	}
	if existsFunc(maps.Values(db.m), out, cmp) {
		return in, errAlreadyExists
	}

	if v, ok := db.m[id]; !ok {
		return v, errNotFound
	}

	db.m[id] = out
	return out, nil
}

// Delete deletes from map the element at id index.
// It returns an error if no element found.
// It returns the deleted element.
func (db *memDB[T]) Delete(_ context.Context, id uint64) (T, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	v, ok := db.m[id]
	if !ok {
		return v, errNotFound
	}

	delete(db.m, id)
	return v, nil
}
