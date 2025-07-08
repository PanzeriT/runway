package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type exitCode int

const (
	Success exitCode = iota
	ErrInvalidFlag
	ErrUnableToDetermineWorkingDirectory
	ErrMissingConfiguration
	ErrCreatingFile
	ErrWritingFile
)

func terminate(reason string, code exitCode) {
	fmt.Println(reason)
	os.Exit(int(code))
}

func CheckError(err error, code exitCode, message ...string) {
	if err == nil {
		return
	}

	if len(message) == 0 {
		terminate(err.Error(), code)
	}

	var sb strings.Builder

	for i, m := range message {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(m)
	}

	terminate(sb.String(), code)
}

func debug(a ...any) {
	if os.Getenv("DEBUG") == "" {
		fmt.Println(a...)
	}
}

func makePathAbsolute(path, base string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(base, path)
}
