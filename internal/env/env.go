package env

import (
	"os"
	"path/filepath"
)

func DefaultEnvFile() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	file := filepath.Join(cwd, ".env")

	return file, nil
}
