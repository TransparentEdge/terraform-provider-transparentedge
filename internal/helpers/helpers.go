package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func GetIntEnv(key string, fallback int) (int, error) {
	s, ok := os.LookupEnv(key)
	if !ok {
		return fallback, fmt.Errorf("variable %s not set", key)
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
		return fallback, fmt.Errorf("variable %s not set", key)
	}

	v, err := strconv.ParseBool(s)
	if err != nil {
		return fallback, err
	}

	return v, nil
}

func SplitAndSort(input string) []string {
	// Split the input string
	words := strings.FieldsFunc(input, func(r rune) bool {
		return r == ' ' || r == '\n' || r == '\r' || r == '\t'
	})

	// Remove duplicate words
	uniqueWords := make(map[string]bool)
	for _, word := range words {
		uniqueWords[word] = true
	}

	// Convert the unique words to a sorted slice
	sortedWords := make([]string, 0, len(uniqueWords))
	for word := range uniqueWords {
		sortedWords = append(sortedWords, word)
	}

	slices.Sort(sortedWords)

	return sortedWords
}

func ParseCertReqLogString(input string) string {
	var logData map[string]any

	err := json.Unmarshal([]byte(input), &logData)
	if err != nil {
		return input
	}

	if en, ok := logData["en"].(string); ok {
		if status, ok := logData["status"].(string); ok {
			return fmt.Sprintf("%s: %s", status, en)
		}

		return en
	}

	return input
}
