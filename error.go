package runway

import (
	"errors"
	"fmt"
	"os"
)

type AppError struct {
	ExitCode int
	err      error
}

var a error

func (e AppError) Error() string {
	return e.err.Error()
}

var ErrSecretToShort = AppError{1, errors.New("secret must be at least 16 characters long")}

func Terminate(err AppError) {
	fmt.Fprintln(os.Stderr, "Fatal Error:", err.Error())
	os.Exit(err.ExitCode)
}
