package configo

import (
	"os"
	"strconv"
	"syscall"
)

var ConfigVariables = map[string]string{}

type Source interface {
	GetVariables() (map[string]string, error)
}

func LoadConfigs(sources ...Source) {

	for _, s := range sources {
		vars, err := s.GetVariables()
		if err != nil {
			panic(err)
		}
		for k, v := range vars {
			ConfigVariables[k] = v
			err := os.Setenv(k, v)
			if err != nil {
				panic(err)
			}
		}
	}
}

func EnvString(key, defaultValue string) string {
	if value, ok := getEnv(key); ok {
		return value
	}

	return defaultValue
}

func EnvInt(key string, defaultValue int) int {
	if value, ok := getEnv(key); ok {
		i, err := strconv.Atoi(value)
		if err == nil {
			return i
		}
	}

	return defaultValue
}

func EnvBool(key string, defaultValue bool) bool {
	if value, ok := getEnv(key); ok {
		i, err := strconv.ParseBool(value)
		if err == nil {
			return i
		}
	}

	return defaultValue
}

func getEnv(key string) (string, bool) {
	value, ok := ConfigVariables[key]
	if ok {
		return value, true
	}
	if value, ok = syscall.Getenv(key); ok {
		ConfigVariables[key] = value
		return value, true
	}

	return "", false
}
