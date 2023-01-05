# oplint

Linter to ensure `op` is correct for every function

## Preamble

To understand where an error originated, define a constant named `op` to prepend an error message with the function name, for example:

```go
package opdemo

import "fmt"

// IsEven returns an error if the number given is not even
func IsEven(n int) error {
    const op = "opdemo/IsEven"

    if n%2 != 0 {
        return fmt.Errorf("%s: %d is not even", op, n)
    }
    return nil
}
```

The error returned is formatted with a helpful locator:
> `opdemo/IsEven: 3 is not even`

The `op` format for a typical function is:

`package name` + `"/"` + `function name`
> e.g. `const op = "opdemo/IsEven"`

If the function has a value or pointer receiver:

```go
package opdemo

import "fmt"

type Number int

// IsEven returns an error if the number is not even
func (n Number) IsEven() error {
    const op = "opdemo/Number.IsEven"

    if n%2 != 0 {
        return fmt.Errorf("%s: %d is not even", op, n)
    }
    return nil
}
```

then the format should be:

`package name` + `"/"` + `type name` + `"."` + `function name`
> e.g. `const op = "opdemo/number.IsEven"`

Adding the `op` to errors adds context which can be critical to understanding your errors, particularly when wrapping errors and sending through a long chain of function calls.

## oplint command

The `oplint` command will scan the body of every function (including those with function receivers) in a given file or directory/package. If a constant is defined as either `op` or `Op`, it will check to see if the value given for said constant matches the function name which envelops it. If the value does not match, `oplint` will report the mismatch and the position.

```sh
$ cat foo.go
package opdemo

import (
    "fmt"
)

type yoda string

// Do does nothing really but returns an error
func (s *yoda) Do() error {
    const op = "opdemo/yoda.Try"

    return fmt.Errorf("%s: There is no try", op)
}

$ oplint foo.go
/Users/gilcrest/opdemo/foo.go:11:8: op constant value (opdemo/yoda.Try) does not match function name (opdemo/yoda.Do)
```

Click on the diagnostic from `oplint` to be taken directly to the spot of the value mismatch and correct it using the diagnostic message.

### oplint -missing flag

Optionally, oplint can report diagnostics on any functions that return an error, but do not have an op constant defined.

```sh
$ cat foo.go
package testdata

import (
    "errors"
    "fmt"
)

func hello() error {
    return errors.New("some error message")
}

$ oplint -missing foo.go
/Users/gilcrest/opdemo/foo.go:8:1: testdata/hello returns an error but does not define an op constant
```

> The default value is false for the `-missing` flag. If it is not present, these checks will not be run.

## Acknowledgements

I started using this `op` pattern after reading Rob Pike's [article on Error handling in Upspin](https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html) back in 2017. Every so often, I go back and read this post, it's one of my favorites.
