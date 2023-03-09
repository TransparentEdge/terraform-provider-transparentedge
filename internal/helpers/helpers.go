package helpers

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func SanitizeStringForDiff(config string) string {
	config = regexp.MustCompile(`[\t\r\n]+`).ReplaceAllString(strings.TrimSpace(config), "\n")
	output := ""

	for _, v := range strings.Split(config, "\n") {
		output += strings.TrimSpace(v) + "\n"
	}

	return output
}

func GetIntEnv(key string, fallback int) (int, error) {
	s, ok := os.LookupEnv(key)
	if !ok {
		return fallback, fmt.Errorf("Variable %s not set", key)
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return fallback, err
	}
	return v, nil
}

func GetEnvBool(key string, fallback bool) (bool, error) {
	s, ok := os.LookupEnv(key)
	if !ok {
		return fallback, fmt.Errorf("Variable %s not set", key)
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return fallback, err
	}
	return v, nil
}
