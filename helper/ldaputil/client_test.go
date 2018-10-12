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
		"":            "",
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

func TestSIDBytesToString(t *testing.T) {
	testcases := map[string][]byte{
		"S-1-5-21-2127521184-1604012920-1887927527-72713": []byte{0x01, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x15, 0x00, 0x00, 0x00, 0xA0, 0x65, 0xCF, 0x7E, 0x78, 0x4B, 0x9B, 0x5F, 0xE7, 0x7C, 0x87, 0x70, 0x09, 0x1C, 0x01, 0x00},
		"S-1-1-0": []byte{0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00},
		"S-1-5":   []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05},
	}

	for answer, test := range testcases {
		res, err := sidBytesToString(test)
		if err != nil {
			t.Errorf("Failed to conver %#v: %s", test, err)
		} else if answer != res {
			t.Errorf("Failed to convert %#v: %s != %s", test, res, answer)
		}
	}
}
