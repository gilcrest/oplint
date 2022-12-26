package testdata

import (
	"fmt"
)

type helloWorld string

// Do does nothing really but returns an error
func (s *helloWorld) Do() error {
	const op = "testdata/helloWorld.Doh" // want `op constant value \(testdata/helloWorld.Doh\) does not match function name \(testdata/helloWorld.Do\)`

	return fmt.Errorf("%s: There is no try", op)
}

func (s helloWorld) DoNot() error {
	const op = "testdata/helloWorld.Do" // want `op constant value \(testdata/helloWorld.Do\) does not match function name \(testdata/helloWorld.DoNot\)`

	return fmt.Errorf("%s: There is no try", op)
}
