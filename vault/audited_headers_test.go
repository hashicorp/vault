package vault

import (
	"reflect"
	"testing"
)

func mockAuditedHeaderConfig(t *testing.T) *AuditedHeadersConfig {
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "foo/")
	return &AuditedHeadersConfig{
		Headers: make(map[string]*auditedHeaderSettings),
		view:    view,
	}
}

func TestAuditedHeaderConfig_CRUD(t *testing.T) {
	conf := mockAuditedHeaderConfig(t)

	testAuditedHeaderConfig_Add(t, conf)
	testAuditedHeaderConfig_Remove(t, conf)
}

func testAuditedHeaderConfig_Add(t *testing.T, conf *AuditedHeadersConfig) {
	err := conf.add("X-Test-Header", false)
	if err != nil {
		t.Fatalf("Error when adding header to config: %s", err)
	}

	settings, ok := conf.Headers["X-Test-Header"]
	if !ok {
		t.Fatal("Expected header to be found in config")
	}

	if settings.HMAC {
		t.Fatal("Expected HMAC to be set to false, got true")
	}

	out, err := conf.view.Get(auditedHeadersEntry)
	if err != nil {
		t.Fatalf("Could not retrieve headers entry from config: %s", err)
	}

	headers := make(map[string]*auditedHeaderSettings)
	err = out.DecodeJSON(&headers)
	if err != nil {
		t.Fatalf("Error decoding header view: %s", err)
	}

	expected := map[string]*auditedHeaderSettings{
		"X-Test-Header": &auditedHeaderSettings{
			HMAC: false,
		},
	}

	if !reflect.DeepEqual(headers, expected) {
		t.Fatalf("Expected config didn't match actual. Expected: %#v, Got: %#v", expected, headers)
	}

	err = conf.add("X-Vault-Header", true)
	if err != nil {
		t.Fatalf("Error when adding header to config: %s", err)
	}

	settings, ok = conf.Headers["X-Vault-Header"]
	if !ok {
		t.Fatal("Expected header to be found in config")
	}

	if !settings.HMAC {
		t.Fatal("Expected HMAC to be set to true, got false")
	}

	out, err = conf.view.Get(auditedHeadersEntry)
	if err != nil {
		t.Fatalf("Could not retrieve headers entry from config: %s", err)
	}

	headers = make(map[string]*auditedHeaderSettings)
	err = out.DecodeJSON(&headers)
	if err != nil {
		t.Fatalf("Error decoding header view: %s", err)
	}

	expected["X-Vault-Header"] = &auditedHeaderSettings{
		HMAC: true,
	}

	if !reflect.DeepEqual(headers, expected) {
		t.Fatalf("Expected config didn't match actual. Expected: %#v, Got: %#v", expected, headers)
	}

}

func testAuditedHeaderConfig_Remove(t *testing.T, conf *AuditedHeadersConfig) {
	err := conf.remove("X-Test-Header")
	if err != nil {
		t.Fatalf("Error when adding header to config: %s", err)
	}

	_, ok := conf.Headers["X-Test-Header"]
	if ok {
		t.Fatal("Expected header to not be found in config")
	}

	out, err := conf.view.Get(auditedHeadersEntry)
	if err != nil {
		t.Fatalf("Could not retrieve headers entry from config: %s", err)
	}

	headers := make(map[string]*auditedHeaderSettings)
	err = out.DecodeJSON(&headers)
	if err != nil {
		t.Fatalf("Error decoding header view: %s", err)
	}

	expected := map[string]*auditedHeaderSettings{
		"X-Vault-Header": &auditedHeaderSettings{
			HMAC: true,
		},
	}

	if !reflect.DeepEqual(headers, expected) {
		t.Fatalf("Expected config didn't match actual. Expected: %#v, Got: %#v", expected, headers)
	}

	err = conf.remove("X-Vault-Header")
	if err != nil {
		t.Fatalf("Error when adding header to config: %s", err)
	}

	_, ok = conf.Headers["X-Vault-Header"]
	if ok {
		t.Fatal("Expected header to not be found in config")
	}

	out, err = conf.view.Get(auditedHeadersEntry)
	if err != nil {
		t.Fatalf("Could not retrieve headers entry from config: %s", err)
	}

	headers = make(map[string]*auditedHeaderSettings)
	err = out.DecodeJSON(&headers)
	if err != nil {
		t.Fatalf("Error decoding header view: %s", err)
	}

	expected = make(map[string]*auditedHeaderSettings)

	if !reflect.DeepEqual(headers, expected) {
		t.Fatalf("Expected config didn't match actual. Expected: %#v, Got: %#v", expected, headers)
	}

}
