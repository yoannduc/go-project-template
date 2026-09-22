package memorydb

import (
	"context"
	"errors"
	"iter"
	"maps"
	"reflect"
	"slices"
)

var (
	errNotFound      = errors.New("not found")
	errAlreadyExists = errors.New("entry already exists")
	errNoIDField     = errors.New("object has no ID field, could not use in DB")
)

const (
	idFieldName = "ID"
)

func updateIDField[T any](in T, id int) (T, error) {
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

type Cmp[T any] func(T, T) int

type MemoryDB[T any] interface {
	FindAll(context.Context) ([]T, error)
	FindFilterFunc(context.Context, func(T) bool) ([]T, error)
	FindByID(context.Context, int) (T, error)
	Create(context.Context, T, Cmp[T]) (T, error)
	Update(context.Context, int, T, Cmp[T]) (T, error)
	Delete(context.Context, int) (T, error)
}

type memDB[T any] struct {
	m     map[int]T
	index int
}

func NewMemoryDB[T any]() MemoryDB[T] {
	return &memDB[T]{
		m:     make(map[int]T, 100),
		index: 1,
	}
}

func (db *memDB[T]) FindAll(_ context.Context) ([]T, error) {
	return slices.Collect(maps.Values(db.m)), nil
}

func (db *memDB[T]) FindFilterFunc(_ context.Context, fltr func(T) bool) ([]T, error) {
	filter := func(in iter.Seq[T]) iter.Seq[T] {
		return func(yield func(T) bool) {
			for v := range in {
				if fltr(v) && !yield(v) {
					return
				}
			}
		}
	}

	return slices.Collect(filter(maps.Values(db.m))), nil
}

func (db *memDB[T]) FindByID(_ context.Context, id int) (T, error) {
	if v, ok := db.m[id]; ok {
		return v, nil
	}

	var v T
	return v, errNotFound
}

func (db *memDB[T]) Create(_ context.Context, in T, cmp Cmp[T]) (T, error) {
	if _, ok := slices.BinarySearchFunc(
		slices.SortedFunc(maps.Values(db.m), cmp),
		in,
		cmp,
	); ok {
		return in, errAlreadyExists
	}

	out, err := updateIDField(in, db.index)
	if err != nil {
		return in, err
	}

	db.m[db.index] = out
	db.index += 1
	return out, nil
}

func (db *memDB[T]) Update(_ context.Context, id int, in T, cmp Cmp[T]) (T, error) {
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

func (db *memDB[T]) Delete(_ context.Context, id int) (T, error) {
	v, ok := db.m[id]
	if !ok {
		return v, errNotFound
	}

	delete(db.m, id)
	return v, nil
}
