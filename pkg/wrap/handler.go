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

func WrapVoid(err error) Output[Void] {
	if err != nil {
		return Err[Void](err)
	}
	return OK[Void](Void{})
}

func ErrProof[T interface {
	ErrorOrNil() error
}](r ...T) Output[Void] {
	for _, v := range r {
		if err := v.ErrorOrNil(); err != nil {
			return Err[Void](err)
		}
	}
	return OK(Void{})
}
