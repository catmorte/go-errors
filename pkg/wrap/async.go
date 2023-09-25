package wrap

import "sync"

type runResult[T any] struct {
	result Result[T]

	sync.Mutex
	value T
	err   error

	valueCh chan T
	errCh   chan error

	resultCh chan Result[T]
}

func (r *runResult[T]) waitResult() Result[T] {
	res, ok := <-r.resultCh
	if ok {
		r.result = res
		close(r.resultCh)
	}
	return r.result
}

// ErrorOrNil implements Result.
func (r *runResult[T]) ErrorOrNil() error {
	return r.waitResult().ErrorOrNil()
}

// Flat implements Result.
func (r *runResult[T]) Flat(onOK func(T), onError func(error)) Result[T] {
	return r.waitResult().Flat(onOK, onError)
}

// GetOrDefault implements Result.
func (r *runResult[T]) GetOrDefault(defaultValue T) T {
	return r.waitResult().GetOrDefault(defaultValue)
}

// GetOrNil implements Result.
func (r *runResult[T]) GetOrNil() *T {
	return r.waitResult().GetOrNil()
}

// IfError implements Result.
func (r *runResult[T]) IfError(onError func(error)) Result[T] {
	return r.waitResult().IfError(onError)
}

// IfOK implements Result.
func (r *runResult[T]) IfOK(onOK func(T)) Result[T] {
	return r.waitResult().IfOK(onOK)
}

// IsError implements Result.
func (r *runResult[T]) IsError() bool {
	return r.waitResult().IsError()
}

// IsOK implements Result.
func (r *runResult[T]) IsOK() bool {
	return r.waitResult().IsOK()
}

func Async[T any](fn func() Result[T]) Result[T] {
	resultCh := make(chan Result[T])
	go func() {
		resultCh <- fn()
	}()
	return &runResult[T]{resultCh: resultCh}
}
