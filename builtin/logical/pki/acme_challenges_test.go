package pki

import (
	"testing"
)

type keyAuthorizationTestCase struct {
	keyAuthz   string
	token      string
	thumbprint string
	shouldFail bool
}

var keyAuthorizationTestCases = []keyAuthorizationTestCase{
	{
		// Entirely empty
		"",
		"non-empty-token",
		"non-empty-thumbprint",
		true,
	},
	{
		// Both empty
		".",
		"non-empty-token",
		"non-empty-thumbprint",
		true,
	},
	{
		// Not equal
		"non-.non-",
		"non-empty-token",
		"non-empty-thumbprint",
		true,
	},
	{
		// Empty thumbprint
		"non-.",
		"non-empty-token",
		"non-empty-thumbprint",
		true,
	},
	{
		// Empty token
		".non-",
		"non-empty-token",
		"non-empty-thumbprint",
		true,
	},
	{
		// Wrong order
		"non-empty-thumbprint.non-empty-token",
		"non-empty-token",
		"non-empty-thumbprint",
		true,
	},
	{
		// Too many pieces
		"one.two.three",
		"non-empty-token",
		"non-empty-thumbprint",
		true,
	},
	{
		// Valid
		"non-empty-token.non-empty-thumbprint",
		"non-empty-token",
		"non-empty-thumbprint",
		false,
	},
}

func TestValidateKeyAuthorization(t *testing.T) {
	for index, tc := range keyAuthorizationTestCases {
		isValid, err := ValidateKeyAuthorization(tc.keyAuthz, tc.token, tc.thumbprint)
		if !isValid && err == nil {
			t.Fatalf("[%d] expected failure to give reason via err (%v / %v)", index, isValid, err)
		}

		expectedValid := !tc.shouldFail
		if expectedValid != isValid {
			t.Fatalf("[%d] got ret=%v, expected ret=%v (shouldFail=%v)", index, isValid, expectedValid, tc.shouldFail)
		}
	}
}
