package activedirectory

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/logical/framework"
)

func TestCertificateValidation(t *testing.T) {

	// certificate should default to "" without error if it doesn't exist
	fd := fieldDataWithSchema()
	config, err := NewConfiguration(hclog.NewNullLogger(), fd)
	if err != nil {
		t.FailNow()
	}
	if config.Certificate != "" {
		t.FailNow()
	}

	// certificate should cause an error if a bad one is provided
	fd.Raw = map[string]interface{}{
		"certificate": "cats",
		"dn":          "example,com",
	}
	config, err = NewConfiguration(hclog.NewNullLogger(), fd)
	if err == nil {
		t.FailNow()
	}

	// valid certificates should pass inspection
	fd.Raw = map[string]interface{}{
		"certificate": validCertificate,
		"dn":          "example,com",
	}
	config, err = NewConfiguration(hclog.NewNullLogger(), fd)
	if err != nil {
		t.FailNow()
	}
}

func TestTLSDefaultsTo12(t *testing.T) {
	fd := fieldDataWithSchema()
	config, err := NewConfiguration(hclog.NewNullLogger(), fd)
	if err != nil {
		t.FailNow()
	}
	expected := uint16(771)
	if config.TLSMinVersion != expected || config.TLSMaxVersion != expected {
		t.FailNow()
	}
}

func TestTLSSessionDefaultsToStarting(t *testing.T) {
	fd := fieldDataWithSchema()
	config, err := NewConfiguration(hclog.NewNullLogger(), fd)
	if err != nil {
		t.FailNow()
	}
	if !config.StartTLS {
		t.FailNow()
	}
}

func TestTLSSessionDefaultsToSecure(t *testing.T) {
	fd := fieldDataWithSchema()
	config, err := NewConfiguration(hclog.NewNullLogger(), fd)
	if err != nil {
		t.FailNow()
	}
	if config.InsecureTLS {
		t.FailNow()
	}
}

func TestRootDomainName(t *testing.T) {
	fd := fieldDataWithSchema()
	fd.Raw = map[string]interface{}{}
	_, err := NewConfiguration(hclog.NewNullLogger(), fd)
	if err == nil {
		t.FailNow()
	}
	fd.Raw = map[string]interface{}{
		"urls": "ldap://138.91.247.105",
		"dn":   "example,com",
	}
	config, err := NewConfiguration(hclog.NewNullLogger(), fd)
	if err != nil {
		t.FailNow()
	}
	if config.RootDomainName != "example,com" {
		t.FailNow()
	}
}

func TestGetTLSConfigs(t *testing.T) {
	fd := fieldDataWithSchema()
	fd.Raw = map[string]interface{}{
		"urls": "ldap://138.91.247.105",
		"dn":   "example,com",
	}
	config, err := NewConfiguration(hclog.NewNullLogger(), fd)
	if err != nil {
		t.FailNow()
	}
	tlsConfigs, err := config.GetTLSConfigs()
	if err != nil {
		t.FailNow()
	}
	if len(tlsConfigs) != 1 {
		t.FailNow()
	}

	for u, tlsConfig := range tlsConfigs {
		if u.String() != "ldap://138.91.247.105" {
			t.FailNow()
		}

		if tlsConfig.InsecureSkipVerify {
			t.FailNow()
		}

		if tlsConfig.ServerName != "138.91.247.105" {
			t.FailNow()
		}

		expected := uint16(771)
		if tlsConfig.MinVersion != expected || tlsConfig.MaxVersion != expected {
			t.FailNow()
		}
	}
}

func fieldDataWithSchema() *framework.FieldData {
	return &framework.FieldData{
		Schema: map[string]*framework.FieldSchema{
			"username": {
				Type:        framework.TypeString,
				Description: "Username with sufficient permissions in Active Directory to administer passwords.",
			},

			"password": {
				Type:        framework.TypeString,
				Description: "Password for username with sufficient permissions in Active Directory to administer passwords.",
			},

			"urls": {
				Type:        framework.TypeCommaStringSlice,
				Default:     "ldap://127.0.0.1",
				Description: "LDAP URL to connect to (default: ldap://127.0.0.1). Multiple URLs can be specified by concatenating them with commas; they will be tried in-order.",
			},

			"certificate": {
				Type:        framework.TypeString,
				Description: "CA certificate to use when verifying LDAP server certificate, must be x509 PEM encoded.",
			},

			"dn": {
				Type:        framework.TypeString,
				Description: "The root distinguished name to bind to when managing service accounts.",
			},

			"insecure_tls": {
				Type:        framework.TypeBool,
				Description: "Skip LDAP server SSL Certificate verification - VERY insecure.",
			},

			"starttls": {
				Type:        framework.TypeBool,
				Default:     true,
				Description: "Issue a StartTLS command after establishing unencrypted connection.",
			},

			"tls_min_version": {
				Type:        framework.TypeString,
				Default:     "tls12",
				Description: "Minimum TLS version to use. Accepted values are 'tls10', 'tls11' or 'tls12'. Defaults to 'tls12'.",
			},

			"tls_max_version": {
				Type:        framework.TypeString,
				Default:     "tls12",
				Description: "Maximum TLS version to use. Accepted values are 'tls10', 'tls11' or 'tls12'. Defaults to 'tls12'.",
			},
		},
		Raw: map[string]interface{}{
			"dn": "example,com",
		},
	}
}

const validCertificate = `
-----BEGIN CERTIFICATE-----
MIIF7zCCA9egAwIBAgIJAOY2qjn64Qq5MA0GCSqGSIb3DQEBCwUAMIGNMQswCQYD
VQQGEwJVUzEQMA4GA1UECAwHTm93aGVyZTERMA8GA1UEBwwIVGltYnVrdHUxEjAQ
BgNVBAoMCVRlc3QgRmFrZTENMAsGA1UECwwETm9uZTEPMA0GA1UEAwwGTm9ib2R5
MSUwIwYJKoZIhvcNAQkBFhZkb25vdHRydXN0QG5vd2hlcmUuY29tMB4XDTE4MDQw
MzIwNDQwOFoXDTE5MDQwMzIwNDQwOFowgY0xCzAJBgNVBAYTAlVTMRAwDgYDVQQI
DAdOb3doZXJlMREwDwYDVQQHDAhUaW1idWt0dTESMBAGA1UECgwJVGVzdCBGYWtl
MQ0wCwYDVQQLDAROb25lMQ8wDQYDVQQDDAZOb2JvZHkxJTAjBgkqhkiG9w0BCQEW
FmRvbm90dHJ1c3RAbm93aGVyZS5jb20wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAw
ggIKAoICAQDzQPGErqjaoFcuUV6QFpSMU6w8wO8F0othik+rrlKERmrGonUGsoum
WqRe6L4ZnxBvCKB6EWjvf894TXOF2cpUnjDAyBePISyPkRBEJS6VS2SEC4AJzmVu
a+P+fZr4Hf7/bEcUr7Ax37yGVZ5i5ByNHgZkBlPxKiGWSmAqIDRZLp9gbu2EkG9q
NOjNLPU+QI2ov6U/laGS1vbE2LahTYeT5yscu9LpllxzFv4lM1f4wYEaM3HuOxzT
l86cGmEr9Q2N4PZ2T0O/s6D4but7c6Bz2XPXy9nWb5bqu0n5bJEpbRFrkryW1ozh
L9uVVz4dyW10pFBJtE42bqA4PRCDQsUof7UfsQF11D1ThrDfKsQa8PxrYdGUHUG9
GFF1MdTTwaoT90RI582p+6XYV+LNlXcdfyNZO9bMThu9fnCvT7Ey0TKU4MfPrlfT
aIhZmyaHt6mL5p881UPDIvy7paTLgL+C1orLjZAiT//c4Zn+0qG0//Cirxr020UF
3YiEFk2H0bBVwOHoOGw4w5HrvLdyy0ZLDSPQbzkSZ0RusHb5TjiyhtTk/h9vvJv7
u1fKJub4MzgrBRi16ejFdiWoVuMXRC6fu/ERy3+9DH6LURerbPrdroYypUmTe9N6
XPeaF1Tc+WO7O/yW96mV7X/D211qjkOtwboZC5kjogVbaZgGzjHCVwIDAQABo1Aw
TjAdBgNVHQ4EFgQU2zWT3HeiMBzusz7AggVqVEL5g0UwHwYDVR0jBBgwFoAU2zWT
3HeiMBzusz7AggVqVEL5g0UwDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOC
AgEAwTGcppY86mNRE43uOimeApTfqHJv+lGDTjEoJCZZmzmtxFe6O9+Vk4bH/8/i
gVQvqzBpaWXRt9OhqlFMK7OkX4ZvqXmnShmxib1dz1XxGhbwSec9ca8bill59Jqa
bIOq2SXVMcFD0GwFxfJRBVzHHuB6AwV9B2QN61zeB1oxNGJrUOo80jVkB7+MWMyD
bQqiFCHWGMa6BG4N91KGOTveZCGdBvvVw5j6lt731KjbvL2hB1UHioucOweKLfa4
QWDImTEjgV68699wKERNL0DCpeD7PcP/L3SY2RJzdyC1CSR7O8yU4lQK7uZGusgB
Mgup+yUaSjxasIqYMebNDDocr5kdwG0+2r2gQdRwc5zLX6YDBn6NLSWjRnY04ZuK
P1cF68rWteWpzJu8bmkJ5r2cqskqrnVK+zz8xMQyEaj548Bnt51ARLHOftR9jkSU
NJWh7zOLZ1r2UUKdDlrMoh3GQO3rvnCJJ16NBM1dB7TUyhMhtF6UOE62BSKdHtQn
d6TqelcRw9WnDsb9IPxRwaXhvGljnYVAgXXlJEI/6nxj2T4wdmL1LWAr6C7DuWGz
8qIvxc4oAau4DsZs2+BwolCFtYc98OjWGcBStBfZz/YYXM+2hKjbONKFxWdEPxGR
Beq3QOqp2+dga36IzQybzPQ8QtotrpSJ3q82zztEvyWiJ7E=
-----END CERTIFICATE-----
`
