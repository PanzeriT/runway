package exit

import (
	"log/slog"
	"os"
)

type ExitCode int

const (
	Ok ExitCode = iota
	ConfigurationError
	DatabaseConnectionError
	DatabaseMigrationError
	ServerError
)

func WithCode(code ExitCode, err error) {
	if code == Ok {
		slog.Info("application exits normally", "code", code)
	} else {
		slog.Error("application exits with error", "code", code, "err", err)
	}

	os.Exit(int(code))
}
