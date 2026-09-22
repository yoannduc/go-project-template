package singleton

// Singleton is a wrapper around a value to use more easily
// like `singl == nil` whereas with basic types like int, 0 can
// sometimes be a value and sometimes a zero value (uninstanciated).
type Singleton[T any] interface {
	GetValue() T
}

// singleton is the concrete type that implements Singleton.
type singleton[T any] struct {
	instance T
}

// New returns a Singleton.
func New[T any](v T) Singleton[T] {
	return &singleton[T]{
		instance: v,
	}
}

// GetValue implements Singleton.GetValue.
func (s singleton[T]) GetValue() T {
	return s.instance
}
