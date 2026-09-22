package memorydb_test

import (
	"context"
	"fmt"
	"strings"

	"github.com/yoannduc/go-project-template/pkg/memorydb"
)

type CustomType struct {
	ID    int
	Field string
}

func cmp(a, b CustomType) bool {
	if a.ID == b.ID {
		return false
	}

	return strings.EqualFold(a.Field, b.Field)
}

func fltr(el CustomType) bool {
	return strings.Contains(el.Field, "yes")
}

func ExampleMemoryDB() {
	db := memorydb.New[CustomType]()
	ctx := context.Background()

	created, err := db.Create(ctx, CustomType{Field: "label"}, cmp)
	if err != nil {
		panic(err)
	}
	fmt.Println(created)

	found, err := db.FindByID(ctx, uint64(1))
	if err != nil {
		panic(err)
	}
	fmt.Println(found)

	updated, err := db.Update(ctx, uint64(1), CustomType{Field: "oh yes it contains"}, cmp)
	if err != nil {
		panic(err)
	}
	fmt.Println(updated)

	filtered, err := db.FindFilterFunc(ctx, fltr)
	if err != nil {
		panic(err)
	}
	fmt.Println(filtered)

	deleted, err := db.Delete(ctx, uint64(1))
	if err != nil {
		panic(err)
	}
	fmt.Println(deleted)

	all, err := db.FindAll(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(all)
}
