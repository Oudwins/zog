package regexconv

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoRE2ToECMA262ConvertsPythonStyleNamedGroups(t *testing.T) {
	pattern := regexp.MustCompile(`^(?P<id>[a-z]+)/(?P<version>[0-9]+)$`).String()

	assert.Equal(t, `^(?:[a-z]+)/(?:[0-9]+)$`, GoRE2ToECMA262(pattern))
}

func TestGoRE2ToECMA262LeavesOtherSyntaxUnchanged(t *testing.T) {
	pattern := `^(?<id>[a-z]+)/(?:[0-9]+)$`

	assert.Equal(t, pattern, GoRE2ToECMA262(pattern))
}
