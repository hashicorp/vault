package logical

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextDisableReplicationStatusEndpointsValue(t *testing.T) {

	testcases := []struct {
		name          string
		ctx           context.Context
		expectedValue bool
		expectedOk    bool
	}{
		{
			name:          "without-value",
			ctx:           context.Background(),
			expectedValue: false,
			expectedOk:    false,
		},
		{
			name:          "with-nil",
			ctx:           context.WithValue(context.Background(), ctxKeyDisableReplicationStatusEndpoints{}, nil),
			expectedValue: false,
			expectedOk:    false,
		},
		{
			name:          "with-incompatible-value",
			ctx:           context.WithValue(context.Background(), ctxKeyDisableReplicationStatusEndpoints{}, "true"),
			expectedValue: false,
			expectedOk:    false,
		},
		{
			name:          "with-bool-true",
			ctx:           context.WithValue(context.Background(), ctxKeyDisableReplicationStatusEndpoints{}, true),
			expectedValue: true,
			expectedOk:    true,
		},
		{
			name:          "with-bool-false",
			ctx:           context.WithValue(context.Background(), ctxKeyDisableReplicationStatusEndpoints{}, false),
			expectedValue: false,
			expectedOk:    true,
		},
	}

	for _, testcase := range testcases {
		value, ok := ContextDisableReplicationStatusEndpointsValue(testcase.ctx)
		assert.Equal(t, testcase.expectedValue, value, testcase.name)
		assert.Equal(t, testcase.expectedOk, ok, testcase.name)
	}
}

func TestCreateContextDisableReplicationStatusEndpoints(t *testing.T) {
	ctx := CreateContextDisableReplicationStatusEndpoints(context.Background(), true)

	value := ctx.Value(ctxKeyDisableReplicationStatusEndpoints{})

	assert.NotNil(t, ctx)
	assert.NotNil(t, value)
	assert.IsType(t, bool(false), value)
	assert.Equal(t, true, value.(bool))

	ctx = CreateContextDisableReplicationStatusEndpoints(context.Background(), false)

	value = ctx.Value(ctxKeyDisableReplicationStatusEndpoints{})

	assert.NotNil(t, ctx)
	assert.NotNil(t, value)
	assert.IsType(t, bool(false), value)
	assert.Equal(t, false, value.(bool))
}

func TestContextMaxRequestSizeValue(t *testing.T) {

	testcases := []struct {
		name          string
		ctx           context.Context
		expectedValue int64
		expectedOk    bool
	}{
		{
			name:          "without-value",
			ctx:           context.Background(),
			expectedValue: 0,
			expectedOk:    false,
		},
		{
			name:          "with-nil",
			ctx:           context.WithValue(context.Background(), ctxKeyMaxRequestSize{}, nil),
			expectedValue: 0,
			expectedOk:    false,
		},
		{
			name:          "with-incompatible-value",
			ctx:           context.WithValue(context.Background(), ctxKeyMaxRequestSize{}, "6666"),
			expectedValue: 0,
			expectedOk:    false,
		},
		{
			name:          "with-int64-8888",
			ctx:           context.WithValue(context.Background(), ctxKeyMaxRequestSize{}, int64(8888)),
			expectedValue: 8888,
			expectedOk:    true,
		},
		{
			name:          "with-int64-zero",
			ctx:           context.WithValue(context.Background(), ctxKeyMaxRequestSize{}, int64(0)),
			expectedValue: 0,
			expectedOk:    true,
		},
	}

	for _, testcase := range testcases {
		value, ok := ContextMaxRequestSizeValue(testcase.ctx)
		assert.Equal(t, testcase.expectedValue, value, testcase.name)
		assert.Equal(t, testcase.expectedOk, ok, testcase.name)
	}
}

func TestCreateContextMaxRequestSize(t *testing.T) {
	ctx := CreateContextMaxRequestSize(context.Background(), int64(8888))

	value := ctx.Value(ctxKeyMaxRequestSize{})

	assert.NotNil(t, ctx)
	assert.NotNil(t, value)
	assert.IsType(t, int64(0), value)
	assert.Equal(t, int64(8888), value.(int64))

	ctx = CreateContextMaxRequestSize(context.Background(), int64(0))

	value = ctx.Value(ctxKeyMaxRequestSize{})

	assert.NotNil(t, ctx)
	assert.NotNil(t, value)
	assert.IsType(t, int64(0), value)
	assert.Equal(t, int64(0), value.(int64))
}

func TestContextOriginalRequestPathValue(t *testing.T) {

	testcases := []struct {
		name          string
		ctx           context.Context
		expectedValue string
		expectedOk    bool
	}{
		{
			name:          "without-value",
			ctx:           context.Background(),
			expectedValue: "",
			expectedOk:    false,
		},
		{
			name:          "with-nil",
			ctx:           context.WithValue(context.Background(), ctxKeyOriginalRequestPath{}, nil),
			expectedValue: "",
			expectedOk:    false,
		},
		{
			name:          "with-incompatible-value",
			ctx:           context.WithValue(context.Background(), ctxKeyOriginalRequestPath{}, 6666),
			expectedValue: "",
			expectedOk:    false,
		},
		{
			name:          "with-string-value",
			ctx:           context.WithValue(context.Background(), ctxKeyOriginalRequestPath{}, "test"),
			expectedValue: "test",
			expectedOk:    true,
		},
		{
			name:          "with-empty-string",
			ctx:           context.WithValue(context.Background(), ctxKeyOriginalRequestPath{}, ""),
			expectedValue: "",
			expectedOk:    true,
		},
	}

	for _, testcase := range testcases {
		value, ok := ContextOriginalRequestPathValue(testcase.ctx)
		assert.Equal(t, testcase.expectedValue, value, testcase.name)
		assert.Equal(t, testcase.expectedOk, ok, testcase.name)
	}
}

func TestCreateContextOriginalRequestPath(t *testing.T) {
	ctx := CreateContextOriginalRequestPath(context.Background(), "test")

	value := ctx.Value(ctxKeyOriginalRequestPath{})

	assert.NotNil(t, ctx)
	assert.NotNil(t, value)
	assert.IsType(t, string(""), value)
	assert.Equal(t, "test", value.(string))

	ctx = CreateContextOriginalRequestPath(context.Background(), "")

	value = ctx.Value(ctxKeyOriginalRequestPath{})

	assert.NotNil(t, ctx)
	assert.NotNil(t, value)
	assert.IsType(t, string(""), value)
	assert.Equal(t, "", value.(string))
}
