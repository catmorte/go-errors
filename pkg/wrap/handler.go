package wrap

func Try[T any, TT any](r Result[T], f func(v T) Result[TT]) Result[TT] {
	if r.IsOK() {
		var defaultV T
		return f(r.GetOrDefault(defaultV))
	}
	return Err[TT](r.ErrorOrNil())
}
