package api

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenAccessor(t *testing.T) {
	testCases := []struct {
		name        string
		f           func() (interface{}, error)
		expectedRet interface{}
		expectedErr error
	}{
		{
			name: "token_accessor_success",
			f: func() (interface{}, error) {
				return (&Secret{
					Auth: &SecretAuth{Accessor: "some_accessor"},
					Data: map[string]interface{}{"accessor": "some_other_accessor"},
				}).TokenAccessor()
			},
			expectedRet: "some_accessor",
			expectedErr: nil,
		},
		{
			name: "token_accessor_other_accessor",
			f: func() (interface{}, error) {
				return (&Secret{
					Data: map[string]interface{}{"accessor": "some_other_accessor"},
				}).TokenAccessor()
			},
			expectedRet: "some_other_accessor",
			expectedErr: nil,
		},
		{
			name: "token_accessor_empty_accessor",
			f: func() (interface{}, error) {
				return (&Secret{
					Auth: &SecretAuth{Accessor: ""},
					Data: map[string]interface{}{"accessor": "some_other_accessor"},
				}).TokenAccessor()
			},
			expectedRet: "some_other_accessor",
			expectedErr: nil,
		},
		{
			name: "token_accessor_nil",
			f: func() (interface{}, error) {
				return (&Secret{
					Data: map[string]interface{}{"some_accessor": "accessor"},
				}).TokenAccessor()
			},
			expectedRet: "",
			expectedErr: nil,
		},
		{
			name: "token_accessor_wrong_type",
			f: func() (interface{}, error) {
				return (&Secret{
					Data: map[string]interface{}{"accessor": 2},
				}).TokenAccessor()
			},
			expectedRet: "",
			expectedErr: errors.New("token found but in the wrong format"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ret, err := testCase.f()
			assert.Equal(t, testCase.expectedRet, ret)
			assert.Equal(t, err, testCase.expectedErr)
		})
	}
}

func TestTokenAccessors(t *testing.T) {
	var emptyRet []string

	testCases := []struct {
		name        string
		f           func() (interface{}, error)
		expectedRet interface{}
		expectedErr error
	}{
		{
			name: "token_accessors_success",
			f: func() (interface{}, error) {
				return (&Secret{
					Data: map[string]interface{}{"keys": []string{"key1", "key2"}},
				}).TokenAccessors()
			},
			expectedRet: []string{"key1", "key2"},
			expectedErr: nil,
		},
		{
			name: "token_accessors_empty_accessor",
			f: func() (interface{}, error) {
				return (&Secret{
					Auth: &SecretAuth{Accessor: ""},
					Data: map[string]interface{}{"keys": []string{"key1", "key2"}},
				}).TokenAccessors()
			},
			expectedRet: []string{"key1", "key2"},
			expectedErr: nil,
		},
		{
			name: "token_accessors_nil",
			f: func() (interface{}, error) {
				return (&Secret{}).TokenAccessors()
			},
			expectedRet: emptyRet,
			expectedErr: nil,
		},
		{
			name: "token_accessors_wrong_type",
			f: func() (interface{}, error) {
				return (&Secret{
					Data: map[string]interface{}{"accessor": 2},
				}).TokenAccessors()
			},
			expectedRet: emptyRet,
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ret, err := testCase.f()
			assert.Exactly(t, testCase.expectedRet, ret)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}
