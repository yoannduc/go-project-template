package singleton

type Singleton[T any] struct {
	Instance T
}
