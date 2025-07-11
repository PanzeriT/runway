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
	ErrDeletingFile
)

func terminate(reason string, code exitCode) {
	fmt.Println(reason)
	os.Exit(int(code))
}

func CheckError(err error, code exitCode, message ...string) {
	if err == nil {
		return
	}

	fmt.Println("Error:", err, message)
	if len(message) == 0 {
		fmt.Println(code)
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

func makePathAbsolute(path, base string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(base, path)
}
