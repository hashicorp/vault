package scram

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestScram(t *testing.T) {
	tcs := map[string]struct {
		Password string
	}{
		"empty-password":  {Password: ""},
		"simple-password": {Password: "password"},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			got, err := Encrypt(tc.Password)
			assert.NoError(t, err)
			assert.True(t, strings.HasPrefix(got, "SCRAM-SHA-256$4096:"))
			assert.Len(t, got, 133)

		})
	}
}
