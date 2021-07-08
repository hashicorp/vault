package configutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePrefixFilters(t *testing.T) {
	prefixFilters := []string{"", "+vault.abc", "-vault.abc", "vault.abc"}

	allowedPrefixes, blockedPrefixes := parsePrefixFilter(prefixFilters)

	assert.Equal(t, len(allowedPrefixes), 1)
	assert.Equal(t, allowedPrefixes[0], prefixFilters[1][1:])

	assert.Equal(t, len(blockedPrefixes), 1)
	assert.Equal(t, blockedPrefixes[0], prefixFilters[2][1:])
}
