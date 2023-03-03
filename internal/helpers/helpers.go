package helpers

import (
	"regexp"
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
