package testdata

import (
	"errors"
	"fmt"
)

func hello() error { // want `testdata/hello returns an error but does not define an op constant`
	return errors.New("some error message")
}

func anotherFuncInAnotherFile() {
	const op = "testdata/anotherFuncInAnotherFile"

	fmt.Println(op)
}

func itWorks() {
	const op = "itDoesNOTWork" // want `op constant value \(itDoesNOTWork\) does not match function name \(testdata/itWorks\)`
	fmt.Println("it works!")
}
