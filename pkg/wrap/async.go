package wrap

import "sync"

type runResult[T any] struct {
	sync.Mutex
	value T
	err   error

	valueCh chan T
	errCh   chan error
}

func (r *runResult[T]) waitResult() (T, error) {
	select {
	case v, ok := <-r.valueCh:
		if ok {
			r.value = v
			r.closeChannels()
		}
	case err, ok := <-r.errCh:
		if ok {
			r.err = err
			r.closeChannels()
		}
	}
	return r.value, r.err
}

func (r *runResult[T]) closeChannels() {
	close(r.valueCh)
	close(r.errCh)
}

// ErrorOrNil implements Result.
func (r *runResult[T]) ErrorOrNil() error {
	_, err := r.waitResult()
	return err
}

// Flat implements Result.
func (r *runResult[T]) Flat(onOK func(T), onError func(error)) Result[T] {
	v, err := r.waitResult()
	if err != nil {
		onError(err)
		return r
	}
	onOK(v)
	return r
}

// GetOrDefault implements Result.
func (r *runResult[T]) GetOrDefault(defaultValue T) T {
	v, err := r.waitResult()
	if err != nil {
		return defaultValue
	}
	return v
}

// GetOrNil implements Result.
func (r *runResult[T]) GetOrNil() *T {
	v, _ := r.waitResult()
	return &v
}

// IfError implements Result.
func (r *runResult[T]) IfError(onError func(error)) Result[T] {
	_, err := r.waitResult()
	if err != nil {
		onError(err)
	}
	return r
}

// IfOK implements Result.
func (r *runResult[T]) IfOK(onOk func(T)) Result[T] {
	v, err := r.waitResult()
	if err == nil {
		onOk(v)
	}
	return r
}

// IsError implements Result.
func (r *runResult[T]) IsError() bool {
	_, err := r.waitResult()
	return err != nil
}

// IsOK implements Result.
func (r *runResult[T]) IsOK() bool {
	_, err := r.waitResult()
	return err == nil
}

func Async[T any](fn func() (T, error)) Result[T] {
	valueCh := make(chan T)
	errCh := make(chan error)
	go func() {
		val, err := fn()
		if err != nil {
			errCh <- err
			return
		}
		valueCh <- val
	}()
	return &runResult[T]{valueCh: valueCh, errCh: errCh}
}
