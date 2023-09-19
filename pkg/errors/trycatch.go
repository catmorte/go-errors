package errors

import "errors"

func empty(e error)            {}
func selfReturn(e error) error { return e }
func selfPanic(e error) error  { panic(e) }

type (
	tryFunc      func() error
	tryErrorFunc func(e error)
	tryCatcher   struct {
		f        tryFunc
		catchers map[error]tryErrorFunc
		finally  tryErrorFunc
	}
)

func (t *tryCatcher) Catch(e error, f tryErrorFunc) *tryCatcher {
	t.catchers[e] = f
	return t
}

func (t *tryCatcher) Finally(f tryErrorFunc) *tryCatcher {
	t.finally = f
	return t
}

func (t tryCatcher) Handle(handle func(e error) error) error {
	err := t.f()
	if t.finally != nil {
		defer t.finally(err)
	}
	if err != nil {
		if !t.catch(err) {
			return handle(err)
		}
	}
	return nil
}

func Try(f tryFunc) *tryCatcher {
	t := &tryCatcher{
		f:        f,
		catchers: make(map[error]tryErrorFunc),
	}
	return t
}

func (t tryCatcher) catch(err error) bool {
	for k := range t.catchers {
		if errors.Is(err, k) {
			return true
		}
	}
	return false
}

func (t tryCatcher) Panic() {
	t.Handle(selfPanic)
}

func (t tryCatcher) Error() error {
	return t.Handle(selfReturn)
}
