package ldaputil

import (
	"testing"
)

func TestLDAPEscape(t *testing.T) {
	testcases := map[string]string{
		"#test":       "\\#test",
		"test,hello":  "test\\,hello",
		"test,hel+lo": "test\\,hel\\+lo",
		"test\\hello": "test\\\\hello",
		"  test  ":    "\\  test \\ ",
	}

	for test, answer := range testcases {
		res := EscapeLDAPValue(test)
		if res != answer {
			t.Errorf("Failed to escape %s: %s != %s\n", test, res, answer)
		}
	}
}

func TestGetTLSConfigs(t *testing.T) {
	config := testConfig()
	if err := config.Validate(); err != nil {
		t.Fatal(err)
	}
	tlsConfig, err := getTLSConfig(config, "138.91.247.105")
	if err != nil {
		t.Fatal(err)
	}
	if tlsConfig == nil {
		t.Fatal("expected 1 TLS config because there's 1 url")
	}
	if tlsConfig.InsecureSkipVerify {
		t.Fatal("InsecureSkipVerify should be false because we should default to the most secure connection")
	}
	if tlsConfig.ServerName != "138.91.247.105" {
		t.Fatalf("expected ServerName of \"138.91.247.105\" but received %q", tlsConfig.ServerName)
	}
	expected := uint16(771)
	if tlsConfig.MinVersion != expected || tlsConfig.MaxVersion != expected {
		t.Fatal("expected TLS min and max version of 771 which corresponds with TLS 1.2 since TLS 1.1 and 1.0 have known vulnerabilities")
	}
}
