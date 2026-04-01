package helpers

import (
	"regexp"
	"strings"
)

var cleanStrRe = regexp.MustCompile(`[\t\r\n]+`)

// VCLSemanticEquals returns true if both VCL configurations are equal semantically.
// it strips trailing newlines and whitespaces to match API response.
func VCLSemanticEquals(c1, c2 string) bool {
	return normalizeVCL(c1) == normalizeVCL(c2)
}

// normalizeVCL normalizes the VCL string to be compatible with the API response.
// API response uses django's: https://www.django-rest-framework.org/api-guide/fields/#charfield (basically s.strip()).
func normalizeVCL(s string) string {
	return strings.TrimSpace(s)
}

// NormalizeStringForComparison normalizes a string for semantic comparison by collapsing
// tabs and newlines, trimming each line, and removing leading/trailing whitespace.
func NormalizeStringForComparison(s string) string {
	s = cleanStrRe.ReplaceAllString(strings.TrimSpace(s), "\n")

	var output strings.Builder
	for v := range strings.SplitSeq(s, "\n") {
		output.WriteString(strings.TrimSpace(v) + "\n")
	}

	return output.String()
}
