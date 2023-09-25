package wrap

func Continue[T any, TT any](r Output[T], f func(v T) Output[TT]) Output[TT] {
	if r.IsOK() {
		var defaultV T
		return f(r.GetOrDefault(defaultV))
	}
	return Err[TT](r.ErrorOrNil())
}
