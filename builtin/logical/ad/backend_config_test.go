package ad

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/builtin/logical/ad/config"
	"github.com/hashicorp/vault/logical"
)

var (
	ctx     = context.Background()
	storage = &logical.InmemStorage{}
)

func TestConfigWriteReadDelete(t *testing.T) {

	b, err := Factory(ctx, &logical.BackendConfig{
		Logger:      hclog.NewNullLogger(),
		StorageView: storage,
	})
	if err != nil {
		t.Error(err)
	}

	// Write
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      config.BackendPath,
		Storage:   storage,
		Data: map[string]interface{}{
			"username":    "tester",
			"password":    "pa$$w0rd",
			"urls":        "ldap://138.91.247.105",
			"certificate": validCertificate,
			"dn":          "example,com",
		},
	}
	resp, err := b.HandleRequest(ctx, req)
	if err != nil {
		t.Error(err)
	}
	verifyResponse(t, resp)
	configIsStoredAsExpected(t)

	// Read
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      config.BackendPath,
		Storage:   storage,
	}
	resp, err = b.HandleRequest(ctx, req)
	if err != nil {
		t.Error(err)
	}
	verifyResponse(t, resp)
	configIsStoredAsExpected(t)

	// Delete
	req = &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      config.BackendPath,
		Storage:   storage,
	}
	resp, err = b.HandleRequest(ctx, req)
	if err != nil {
		t.Error(err)
	}
	if resp != nil {
		t.Error("expected a nil resp, to provide a 204 with no body as the outer response")
	}
	entry, err := storage.Get(ctx, config.StorageKey)
	if err != nil {
		t.Error(err)
	}
	if entry != nil {
		t.Error("config should no longer be stored")
	}
}

func configIsStoredAsExpected(t *testing.T) {
	entry, err := storage.Get(ctx, config.StorageKey)
	if err != nil {
		t.Error(err)
	}

	engineConf := &config.EngineConf{}
	if err := entry.DecodeJSON(engineConf); err != nil {
		t.Error(err)
	}

	if engineConf.ADConf.Certificate != validCertificate {
		t.Errorf("expected Certificate to be the given valid certificate but received \"%s\"", engineConf.ADConf.Certificate)
	}

	if engineConf.ADConf.RootDomainName != "example,com" {
		t.Errorf("expected RootDomainName to be \"example,com\" but received \"%s\"", engineConf.ADConf.RootDomainName)
	}

	if engineConf.ADConf.InsecureTLS {
		t.Error("expected InsecureTLS to be false but received true")
	}

	if engineConf.ADConf.Password != "pa$$w0rd" {
		t.Errorf("expected Password to be \"pa$$w0rd\" but received \"%s\"", engineConf.ADConf.Password)
	}

	if !engineConf.ADConf.StartTLS {
		t.Error("expected StartTLS to be true but received false")
	}

	if engineConf.ADConf.URLs[0] != "ldap://138.91.247.105" {
		t.Errorf("expected the url to be \"ldap://138.91.247.105\" but received \"%s\"", engineConf.ADConf.URLs[0])
	}

	if engineConf.ADConf.TLSMinVersion != 771 {
		t.Errorf("expected TLSMinVersion to be \"771\" but received \"%d\"", engineConf.ADConf.TLSMinVersion)
	}

	if engineConf.ADConf.TLSMaxVersion != 771 {
		t.Errorf("expected TLSMaxVersion to be \"771\" but received \"%d\"", engineConf.ADConf.TLSMaxVersion)
	}

	if engineConf.ADConf.Username != "tester" {
		t.Errorf("expected Username to be \"tester\" but received \"%s\"", engineConf.ADConf.Username)
	}

	if engineConf.PasswordConf.DefaultPasswordTTL != config.DefaultPasswordTTLs {
		t.Errorf("received unexpected DefaultPasswordTTL of \"%d\"", engineConf.PasswordConf.DefaultPasswordTTL)
	}

	if engineConf.PasswordConf.MaxPasswordTTL != config.DefaultPasswordTTLs {
		t.Errorf("received unexpected MaxPasswordTTL of \"%d\"", engineConf.PasswordConf.MaxPasswordTTL)
	}

	if engineConf.PasswordConf.PasswordLength != config.DefaultPasswordLength {
		t.Errorf("received unexpected PasswordLength of \"%d\"", engineConf.PasswordConf.PasswordLength)
	}
}

func verifyResponse(t *testing.T, resp *logical.Response) {

	// Did we get the response data we expect?
	if resp.Data["certificate"] != "\n-----BEGIN CERTIFICATE-----\nMIIF7zCCA9egAwIBAgIJAOY2qjn64Qq5MA0GCSqGSIb3DQEBCwUAMIGNMQswCQYD\nVQQGEwJVUzEQMA4GA1UECAwHTm93aGVyZTERMA8GA1UEBwwIVGltYnVrdHUxEjAQ\nBgNVBAoMCVRlc3QgRmFrZTENMAsGA1UECwwETm9uZTEPMA0GA1UEAwwGTm9ib2R5\nMSUwIwYJKoZIhvcNAQkBFhZkb25vdHRydXN0QG5vd2hlcmUuY29tMB4XDTE4MDQw\nMzIwNDQwOFoXDTE5MDQwMzIwNDQwOFowgY0xCzAJBgNVBAYTAlVTMRAwDgYDVQQI\nDAdOb3doZXJlMREwDwYDVQQHDAhUaW1idWt0dTESMBAGA1UECgwJVGVzdCBGYWtl\nMQ0wCwYDVQQLDAROb25lMQ8wDQYDVQQDDAZOb2JvZHkxJTAjBgkqhkiG9w0BCQEW\nFmRvbm90dHJ1c3RAbm93aGVyZS5jb20wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAw\nggIKAoICAQDzQPGErqjaoFcuUV6QFpSMU6w8wO8F0othik+rrlKERmrGonUGsoum\nWqRe6L4ZnxBvCKB6EWjvf894TXOF2cpUnjDAyBePISyPkRBEJS6VS2SEC4AJzmVu\na+P+fZr4Hf7/bEcUr7Ax37yGVZ5i5ByNHgZkBlPxKiGWSmAqIDRZLp9gbu2EkG9q\nNOjNLPU+QI2ov6U/laGS1vbE2LahTYeT5yscu9LpllxzFv4lM1f4wYEaM3HuOxzT\nl86cGmEr9Q2N4PZ2T0O/s6D4but7c6Bz2XPXy9nWb5bqu0n5bJEpbRFrkryW1ozh\nL9uVVz4dyW10pFBJtE42bqA4PRCDQsUof7UfsQF11D1ThrDfKsQa8PxrYdGUHUG9\nGFF1MdTTwaoT90RI582p+6XYV+LNlXcdfyNZO9bMThu9fnCvT7Ey0TKU4MfPrlfT\naIhZmyaHt6mL5p881UPDIvy7paTLgL+C1orLjZAiT//c4Zn+0qG0//Cirxr020UF\n3YiEFk2H0bBVwOHoOGw4w5HrvLdyy0ZLDSPQbzkSZ0RusHb5TjiyhtTk/h9vvJv7\nu1fKJub4MzgrBRi16ejFdiWoVuMXRC6fu/ERy3+9DH6LURerbPrdroYypUmTe9N6\nXPeaF1Tc+WO7O/yW96mV7X/D211qjkOtwboZC5kjogVbaZgGzjHCVwIDAQABo1Aw\nTjAdBgNVHQ4EFgQU2zWT3HeiMBzusz7AggVqVEL5g0UwHwYDVR0jBBgwFoAU2zWT\n3HeiMBzusz7AggVqVEL5g0UwDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOC\nAgEAwTGcppY86mNRE43uOimeApTfqHJv+lGDTjEoJCZZmzmtxFe6O9+Vk4bH/8/i\ngVQvqzBpaWXRt9OhqlFMK7OkX4ZvqXmnShmxib1dz1XxGhbwSec9ca8bill59Jqa\nbIOq2SXVMcFD0GwFxfJRBVzHHuB6AwV9B2QN61zeB1oxNGJrUOo80jVkB7+MWMyD\nbQqiFCHWGMa6BG4N91KGOTveZCGdBvvVw5j6lt731KjbvL2hB1UHioucOweKLfa4\nQWDImTEjgV68699wKERNL0DCpeD7PcP/L3SY2RJzdyC1CSR7O8yU4lQK7uZGusgB\nMgup+yUaSjxasIqYMebNDDocr5kdwG0+2r2gQdRwc5zLX6YDBn6NLSWjRnY04ZuK\nP1cF68rWteWpzJu8bmkJ5r2cqskqrnVK+zz8xMQyEaj548Bnt51ARLHOftR9jkSU\nNJWh7zOLZ1r2UUKdDlrMoh3GQO3rvnCJJ16NBM1dB7TUyhMhtF6UOE62BSKdHtQn\nd6TqelcRw9WnDsb9IPxRwaXhvGljnYVAgXXlJEI/6nxj2T4wdmL1LWAr6C7DuWGz\n8qIvxc4oAau4DsZs2+BwolCFtYc98OjWGcBStBfZz/YYXM+2hKjbONKFxWdEPxGR\nBeq3QOqp2+dga36IzQybzPQ8QtotrpSJ3q82zztEvyWiJ7E=\n-----END CERTIFICATE-----\n" {
		t.Errorf("expected certificate to be the given one but received \"%s\"", resp.Data["certificate"])
	}

	if resp.Data["dn"] != "example,com" {
		t.Errorf("expected dn to be \"example,com\" but received \"%s\"", resp.Data["dn"])
	}

	if resp.Data["insecure_tls"].(bool) {
		t.Errorf("expected insecure_tls to be false but received true")
	}

	if resp.Data["password"] != "pa$$w0rd" {
		t.Errorf("expected password to be \"pa$$w0rd\" but received \"%s\"", resp.Data["password"])
	}

	if !resp.Data["starttls"].(bool) {
		t.Errorf("expected starttls to be true but received false")
	}

	if fmt.Sprintf("%s", resp.Data["urls"]) != `[ldap://138.91.247.105]` {
		t.Errorf("expected url to be \"[ldap://138.91.247.105]\" but received \"%s\"", fmt.Sprintf("%s", resp.Data["urls"]))
	}

	if resp.Data["tlsminversion"].(uint16) != 771 {
		t.Errorf("expected tlsminversion to be \"771\" but received \"%s\"", resp.Data["tlsminversion"])
	}

	if resp.Data["tlsmaxversion"].(uint16) != 771 {
		t.Errorf("expected tlsmaxversion to be \"771\" but received \"%s\"", resp.Data["tlsmaxversion"])
	}

	if resp.Data["username"] != "tester" {
		t.Errorf("expected username to be \"tester\" but received \"%s\"", resp.Data["username"])
	}

	if resp.Data["default_password_ttl"] != config.DefaultPasswordTTLs {
		t.Errorf("received unexpected default_password_ttl of \"%d\"", resp.Data["default_password_ttl"])
	}

	if resp.Data["max_password_ttl"] != config.DefaultPasswordTTLs {
		t.Errorf("received unexpected default_password_ttl of \"%d\"", resp.Data["max_password_ttl"])
	}

	if resp.Data["password_length"] != config.DefaultPasswordLength {
		t.Errorf("received unexpected default_password_ttl of \"%d\"", resp.Data["password_length"])
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
