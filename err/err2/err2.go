package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func main() {
	if err := doSomething(); err != nil {
		switch errors.Cause(err).(type) {
		case *SomeError:
			fmt.Fprintln(os.Stderr, "*SomeError")
		default:
		}
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

type SomeError struct {
	Code    int
	Message string
}

func (s *SomeError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", s.Code, s.Message)
}

func doSomething() error {
	return open()
}

func open() error {
	return &SomeError{Code: 400, Message: "invalid open"}
}
