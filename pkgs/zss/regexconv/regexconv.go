package regexconv

import "regexp"

var pythonNamedGroupPattern = regexp.MustCompile(`\(\?P<[^>]+>`)

// GoRE2ToECMA262 converts Go/RE2 regex syntax to ECMA-262 syntax.
// For now, it only removes Go-supported Python-style named capture groups.
func GoRE2ToECMA262(pattern string) string {
	return pythonNamedGroupPattern.ReplaceAllString(pattern, "(?:")
}
