package command

import (
	"flag"
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestFlagSet(t *testing.T) {
	cases := []struct {
		Flags    FlagSetFlags
		Expected []string
	}{
		{
			FlagSetNone,
			[]string{},
		},
		{
			FlagSetServer,
			[]string{"address", "ca-cert", "ca-path", "client-cert", "client-key", "insecure", "tls-skip-verify"},
		},
	}

	for i, tc := range cases {
		var m Meta
		fs := m.FlagSet("foo", tc.Flags)

		actual := make([]string, 0, 0)
		fs.VisitAll(func(f *flag.Flag) {
			actual = append(actual, f.Name)
		})
		sort.Strings(actual)
		sort.Strings(tc.Expected)

		if !reflect.DeepEqual(actual, tc.Expected) {
			t.Fatalf("%d: flags: %#v\n\nExpected: %#v\nGot: %#v",
				i, tc.Flags, tc.Expected, actual)
		}
	}
}

func TestEnvSettings(t *testing.T) {
	os.Setenv("VAULT_CACERT", "/path/to/fake/cert.crt")
	os.Setenv("VAULT_CAPATH", "/path/to/fake/certs")
	os.Setenv("VAULT_CLIENT_CERT", "/path/to/fake/client.crt")
	os.Setenv("VAULT_CLIENT_KEY", "/path/to/fake/client.key")
	os.Setenv("VAULT_SKIP_VERIFY", "true")
	defer os.Setenv("VAULT_CACERT", "")
	defer os.Setenv("VAULT_CAPATH", "")
	defer os.Setenv("VAULT_CLIENT_CERT", "")
	defer os.Setenv("VAULT_CLIENT_KEY", "")
	defer os.Setenv("VAULT_SKIP_VERIFY", "")
	var m Meta

	// Err is ignored as it is expected that the test settings
	// will cause errors; just check the flag settings
	m.Client()

	if m.flagCACert != "/path/to/fake/cert.crt" {
		t.Fatalf("bad: %s", m.flagAddress)
	}
	if m.flagCAPath != "/path/to/fake/certs" {
		t.Fatalf("bad: %s", m.flagAddress)
	}
	if m.flagClientCert != "/path/to/fake/client.crt" {
		t.Fatalf("bad: %s", m.flagAddress)
	}
	if m.flagClientKey != "/path/to/fake/client.key" {
		t.Fatalf("bad: %s", m.flagAddress)
	}
	if m.flagInsecure != true {
		t.Fatalf("bad: %s", m.flagAddress)
	}
}
