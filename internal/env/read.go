package env

import (
	"os"
	"strings"
)

func Load(filename string, overload bool) error {
	envMap, err := parseFile(filename)
	if err != nil {
		return err
	}

	currentEnv := map[string]bool{}
	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}

	for key, value := range envMap {
		if !currentEnv[key] || overload {
			_ = os.Setenv(key, value)
		}
	}

	return nil
}

func parseFile(filename string) (envMap map[string]string, err error) {
	envMap = make(map[string]string)

	var src []byte
	src, err = os.ReadFile(filename)
	if err != nil {
		return
	}

	err = parseBytes(src, envMap)

	return
}
