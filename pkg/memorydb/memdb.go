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

// Compare is the signature for the compare function
// performed on binary searches.
type Compare[T any] func(T, T) int

// Filter is the signature for the filter function applied to element list.
type Filter[T any] func(T) bool

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
	Update(context.Context, uint64, T, Compare[T]) (T, error)
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
		m: make(map[uint64]T),
	}
	db.index.Add(1)

	return db
}

// FindAll returns a slice of all element in internal map.
func (db *memDB[T]) FindAll(_ context.Context) ([]T, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	return slices.Collect(maps.Values(db.m)), nil
}

// FindFilterFunc returns a slice of filtered element in internal map.
// It takes in a closure to let caller define how results should be filtered.
func (db *memDB[T]) FindFilterFunc(_ context.Context, fltr Filter[T]) ([]T, error) {
	filter := func(in iter.Seq[T]) iter.Seq[T] {
		return func(yield func(T) bool) {
			for v := range in {
				if fltr(v) && !yield(v) {
					return
				}
			}
		}
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	return slices.Collect(filter(maps.Values(db.m))), nil
}

// FindByID returns the element at id index.
// It returns an error if no element found.
func (db *memDB[T]) FindByID(_ context.Context, id uint64) (T, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	v, ok := db.m[id]
	if ok {
		return v, nil
	}

	return v, errNotFound
}

// Create inserts in in the map. It performs a binary search ordered
// by cmp to make sure all elements in db are unique. It returns an
// error if duplicate element was found or if it could not update ID
// field in element (see [MemoryDB] for more info on element restrictions).
// It returns the created element.
func (db *memDB[T]) Create(_ context.Context, in T, cmp Compare[T]) (T, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := slices.BinarySearchFunc(
		slices.SortedFunc(maps.Values(db.m), cmp),
		in,
		cmp,
	); ok {
		return in, errAlreadyExists
	}

	idx := db.index.Load()
	out, err := updateIDField(in, idx)
	if err != nil {
		return in, err
	}

	db.m[idx] = out
	db.index.Add(1)
	return out, nil
}

// Update updates in at id index. It performs a binary search ordered
// by cmp to make sure all elements in db are unique. It returns an
// error if duplicate element was found or if it could not update ID
// field in element (see [MemoryDB] for more info on element restrictions).
// It returns the updated element.
func (db *memDB[T]) Update(_ context.Context, id uint64, in T, cmp Compare[T]) (T, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := slices.BinarySearchFunc(
		slices.SortedFunc(maps.Values(db.m), cmp),
		in,
		cmp,
	); ok {
		return in, errAlreadyExists
	}

	if v, ok := db.m[id]; !ok {
		return v, errNotFound
	}

	out, err := updateIDField(in, id)
	if err != nil {
		return in, err
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
