package wrap

type Void struct{}

type Output[T any] interface {
	IsOK() bool
	IsError() bool

	GetOrDefault(defaultValue T) T
	GetOrNil() *T
	ErrorOrNil() error

	IfOK(onOk func(T)) Output[T]
	IfError(onError func(error)) Output[T]
	Flat(onOK func(T), onError func(error)) Output[T]

	Unwrap() (T, error)
}
