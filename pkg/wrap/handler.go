package wrap

func Continue[T any, TT any](r Output[T], f func(v T) Output[TT]) Output[TT] {
	if r.IsOK() {
		var defaultV T
		return f(r.GetOrDefault(defaultV))
	}
	return Err[TT](r.ErrorOrNil())
}

func ContinueAsync[T any, TT any](r Output[T], f func(v T) Output[TT]) Output[TT] {
	return Async[TT](func() Output[TT] {
		if r.IsOK() {
			var defaultV T
			return f(r.GetOrDefault(defaultV))
		}
		return Err[TT](r.ErrorOrNil())
	})
}

func Wrap[T any](val T, err error) Output[T] {
	if err != nil {
		return Err[T](err)
	}
	return OK[T](val)
}

func WrapVoid(err error) Output[struct{}] {
	if err != nil {
		return Err[struct{}](err)
	}
	return OK[struct{}](struct{}{})
}
