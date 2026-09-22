package singleton_test

import (
	"fmt"
	"math/rand/v2"
	"sync"

	"github.com/yoannduc/go-project-template/pkg/singleton"
)

type CustomType int

var (
	once  sync.Once
	singl singleton.Singleton[CustomType]
)

func GetSingleValue() CustomType {
	if singl == nil {
		once.Do(func() {
			v := CustomType(rand.Int())

			singl = singleton.New(v)
		})
	}

	return singl.GetValue()
}

func ExampleSingleton() {
	v1 := GetSingleValue()
	v2 := GetSingleValue()
	v3 := GetSingleValue()

	fmt.Println(v1)
	fmt.Println(v2)
	fmt.Println(v3)
}
