package generator

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type exitCode int

const (
	Success exitCode = iota
	ErrInvalidFlag
	ErrUnableToDetermineWorkingDirectory
	ErrMissingConfiguration
	ErrCreatingFile
	ErrWritingFile
	ErrCopyingFile
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

func copyFSFile(srcFS fs.FS, src, dst string) error {
	srcFile, err := srcFS.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	content, err := io.ReadAll(srcFile)
	if err != nil {
		return err
	}

	modified := bytes.ReplaceAll(content,
		[]byte("MODULE_NAME"),
		[]byte(viper.GetString("module_name")+"/internal/server/admin"),
	)
	return os.WriteFile(dst, modified, 0644)
}
