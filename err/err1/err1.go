package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func main() {
	err := doSomething()
	if err != nil {
		Debugf("%+v\n", err)
	}
}

func Debugf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, "[DEBUG] "+format+"\n", args...)
}

func doSomething() error {
	err := read()
	if err != nil {
		return errors.Wrap(err, "faild")
	}
	return nil
}

func read() error {
	return fmt.Errorf("エラー")
}
