package wrap

type Result[T any] interface {
	IsOK() bool
	IsError() bool

	GetOrDefault(defaultValue T) T
	GetOrNil() *T
	ErrorOrNil() error

	IfOK(onOk func(T)) Result[T]
	IfError(onError func(error)) Result[T]
	Flat(onOK func(T), onError func(error)) Result[T]
}
