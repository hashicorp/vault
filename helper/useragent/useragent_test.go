package useragent

import (
	"testing"
)

func TestUserAgent(t *testing.T) {
	projectURL = "https://vault-test.com"
	rt = "go5.0"
	versionFunc = func() string { return "1.2.3" }

	act := String()

	exp := "Vault/1.2.3 (+https://vault-test.com; go5.0)"
	if exp != act {
		t.Errorf("expected %q to be %q", act, exp)
	}
}
