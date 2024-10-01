package scram

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestScram tests the Hash method. The hashed password string should have a SCRAM-SHA-256 prefix.
func TestScram(t *testing.T) {
	tcs := map[string]struct {
		Password string
	}{
		"empty-password":  {Password: ""},
		"simple-password": {Password: "password"},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			got, err := Hash(tc.Password)
			assert.NoError(t, err)
			assert.True(t, strings.HasPrefix(got, "SCRAM-SHA-256$4096:"))
			assert.Len(t, got, 133)
		})
	}
}
