# go-errors

## "errors" package for try/catch

To work with plain errors

```
import . "github.com/catmorte/go-errors/pkg/tcatch"
...

var plainError = errors.New("plain error")
...

unexpectedErr := Try(func() error {
    return fmt.Errorf("wrapped: %w", plainError)
}).Catch(plainError, func(err error) {
    fmt.Println("catched %w", err)
}).Error()

```

To work with custom errors

```
import . "github.com/catmorte/go-errors/pkg/tcatch"
...

type CustomError struct {
    Err error
}

func (c CustomError) Error() string {
	if c.Err == nil {
		return "not an error"
	}
	return c.Err.Error()
}

func (e CustomError) Is(err error) bool {
	var castedErr *CustomError
	if ok := errors.As(err, &castedErr); ok {
		if e.Err == nil {
			return true
		}
		return errors.Is(e.Err, castedErr.Err)
	} else {
		var castedErr CustomError
		if ok := errors.As(err, &castedErr); ok {
			if e.Err == nil {
				return true
			}
			return errors.Is(e.Err, castedErr.Err)
		}
	}
	return false
}
...

unexpectedErr := Try(func() error {
    return fmt.Errorf("wrapped: %w", CustomError{Err: plainError})
}).Catch(&CustomError{}, func(err error) {
    fmt.Println("catched %w", err)
}).Error()

```

Multiple catch statements are allowed

```
unexpectedErr := Try(...).Catch(...).Catch(...).Error()
```

Instead of returning unexpected error it's possible to Panic

```
Try(...).Catch(...).Catch(...).Panic()
```

Also, it's possible to set custom error handler in case Catch method is not enough

```
Try(...).Handle(func(error)error)
```

In any case it's possible to add Finally statement which wwill be executed in any case

```
Try(...).Finally(func(error)).Handle(func(error)error)

unexpectedErr := Try(...).Catch(...).Catch(...)..Finally(func(error)).Error()

Try(...).Catch(...).Catch(...)..Finally(func(error)).Panic()
```

## "wrap" package - to handle error result and build code as pipeline

The package shows another approach of work with functions that returns an error next to the value. It requires to wrap result with an interface that follow the following format:

```
type Output[T any] interface {
	IsOK() bool
	IsError() bool

	GetOrDefault(defaultValue T) T
	GetOrNil() *T
	ErrorOrNil() error

	IfOK(onOk func(T)) Output[T]
	IfError(onError func(error)) Output[T]
	Flat(onOK func(T), onError func(error)) Output[T]
}
```

There are 2 predefined constructors available:

```
import . "github.com/catmorte/go-errors/pkg/wrap"
...

var output Output[ResultType] = OK[ResultType](valueOfResultType)
var output Output[ResultType] = Err[ResultType](err)
```

In this case your next function can look like this

```
func DoStuff() Output[ResultType] {
    if !ok {
        return Err[ResultType](errors.New("not ok"))
    }
    return OK[ResultType](new(ResultType))
}
```

In case of async function, execution will wait for function to complete

```
var result Output[ResultType] = wrap.Async[ResultType](func() wrap.Output[ResultType] {
	time.Sleep(time.Second * 10)
	return wrap.OK(new(ResultType))
})
```

Functions to help structurize code into pipeline

```
Continue[T any, TT any](r Output[T], f func(v T) Output[TT]) Output[TT]

ContinueAsync[T any, TT any](r Output[T], f func(v T) Output[TT]) Output[TT]
```

### Example

```
r := Async[int](func() Output[int] {
	time.Sleep(time.Second * 10)
	return OK(1)
})

rr := ContinueAsync[int, string](r, func(v int) Output[string] {
	time.Sleep(time.Second * 10)
	return OK[string]("Result" + strconv.Itoa(v))
})
fmt.Println(rr.GetOrDefault(""))
```
