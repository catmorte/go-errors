package wrap

type ok[T any] struct {
	v T
}

func (ok[T]) ErrorOrNil() error {
	return nil
}

func (s ok[T]) Flat(onOK func(T), onError func(error)) Result[T] {
	onOK(s.v)
	return s
}

func (s ok[T]) GetOrDefault(defaultValue T) T {
	return s.v
}

func (s ok[T]) GetOrNil() *T {
	return &s.v
}

func (ok[T]) IsError() bool {
	return false
}

func (ok[T]) IsOK() bool {
	return true
}

func (s ok[T]) IfError(onError func(error)) Result[T] {
	return s
}

func (s ok[T]) IfOK(onOK func(T)) Result[T] {
	onOK(s.v)
	return s
}

func OK[T any](value T) Result[T] {
	output := new(ok[T])
	output.v = value

	return *output
}
