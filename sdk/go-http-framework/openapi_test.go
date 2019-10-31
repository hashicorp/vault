package openapi

import (
	"encoding/json"
	"errors"
	"strings"
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

func TestResponseHeaderMarshal(t *testing.T) {
	t.Run("content-type-header", func(t *testing.T) {
		o := &Response{
			Headers: map[string]interface{}{
				"conTent-TYpe": Reference{Ref: "foo"},
				"Ref-Val":      Reference{Ref: "bar"},
			},
		}
		j, err := json.Marshal(o)
		assert.NoError(t, err)
		assert.False(t, strings.Contains(strings.ToLower(string(j)), "content-type"))
		assert.True(t, strings.Contains(strings.ToLower(string(j)), "ref-val"))
		t.Log(string(j))
	})

	t.Run("dup-header", func(t *testing.T) {
		o := &Response{
			Headers: map[string]interface{}{
				"Ref-val": Reference{Ref: "foo"},
				"rEf-VaL": Reference{Ref: "bar"},
			},
		}
		_, err := json.Marshal(o)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrDuplicateHeader))
	})

	t.Run("invalid-header-target", func(t *testing.T) {
		o := &Response{
			Headers: map[string]interface{}{
				"Ref-val": true,
			},
		}
		_, err := json.Marshal(o)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInvalidHeaderTarget))
	})

	t.Run("invalid-content-target", func(t *testing.T) {
		o := &Response{
			Content: map[string]interface{}{
				"application/json": true,
			},
		}
		_, err := json.Marshal(o)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInvalidContentTarget))
	})
}
