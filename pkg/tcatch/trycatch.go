package tcatch

import (
	"errors"
)

func empty(e error)            {}
func selfReturn(e error) error { return e }
func selfPanic(e error) error  { panic(e) }

type (
	tryFunc      func() error
	tryErrorFunc func(e error)
	errFunc      struct {
		err error
		f   tryErrorFunc
	}
	tryCatcher struct {
		f        tryFunc
		catchers []errFunc
		finally  tryErrorFunc
	}
)

func (t *tryCatcher) Catch(e error, f tryErrorFunc) *tryCatcher {
	t.catchers = append(t.catchers, errFunc{err: e, f: f})
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
		catchers: make([]errFunc, 0),
	}
	return t
}

func (t tryCatcher) catch(err error) bool {
	for _, ef := range t.catchers {
		if errors.Is(ef.err, err) || errors.Is(err, ef.err) {
			ef.f(err)
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
