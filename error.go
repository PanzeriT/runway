package runway

import (
	"errors"
	"os"
)

type AppError struct {
	ExitCode  int
	Err       error
	DetailErr error
}

var a error

func (e AppError) Error() string {
	return e.Err.Error()
}

var (
	ErrSecretToShort       = AppError{1, errors.New("secret must be at least 16 characters long"), nil}
	ErrCannotRegisterModel = AppError{2, errors.New("cannot register model"), nil}
)

func NewAppError(base AppError, err error) AppError {
	e := AppError{
		ExitCode:  base.ExitCode,
		Err:       base.Err,
		DetailErr: err,
	}
	return e
}

func Terminate(err AppError) {
	if err.DetailErr != nil {
		logger.Info("fatal error", "message", err.Error(), "detail_message", err.DetailErr.Error(), "exit_code", err.ExitCode)
	} else {
		logger.Info("fatal error", "message", err.Error(), "exit_code", err.ExitCode)
	}

	os.Exit(err.ExitCode)
}
