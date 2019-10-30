package openapi

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenericMarshal(t *testing.T) {
	t.Run("noext", func(t *testing.T) {
		o := &OpenAPI{
			OpenAPI: "3.0.2",
		}
		j, err := json.Marshal(o)
		assert.NoError(t, err)
		t.Log(string(j))
	})

	t.Run("valid_ext", func(t *testing.T) {
		o := &OpenAPI{
			OpenAPI:    "3.0.2",
			Extensions: make(Extensions),
		}
		o.Extensions["foo"] = "bar"
		j, err := json.Marshal(o)
		assert.NoError(t, err)
		t.Log(string(j))
	})
}
