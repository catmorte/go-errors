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
    fmt.Println("catch %w", err)
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
