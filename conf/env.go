package conf

import (
	"os"
	"strconv"
)

// EnvInt returns the integer environment variable as config.
func EnvInt(name string, defaultValue int) (value int) {
	value, _ = strconv.Atoi(os.Getenv(name))
	if value == 0 {
		value = defaultValue
	}

	return
}

// EnvStr returns the string environment variable as config.
func EnvStr(name, defaultValue string) (value string) {
	value = os.Getenv(name)
	if value == "" {
		value = defaultValue
	}

	return
}
