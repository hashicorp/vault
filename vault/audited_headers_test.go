package vault

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/helper/salt"
)

func mockAuditedHeadersConfig(t *testing.T) *AuditedHeadersConfig {
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "foo/")
	return &AuditedHeadersConfig{
		Headers: make(map[string]*auditedHeaderSettings),
		view:    view,
	}
}

func TestAuditedHeadersConfig_CRUD(t *testing.T) {
	conf := mockAuditedHeadersConfig(t)

	testAuditedHeadersConfig_Add(t, conf)
	testAuditedHeadersConfig_Remove(t, conf)
}

func testAuditedHeadersConfig_Add(t *testing.T, conf *AuditedHeadersConfig) {
	err := conf.add(context.Background(), "X-Test-Header", false)
	if err != nil {
		t.Fatalf("Error when adding header to config: %s", err)
	}

	settings, ok := conf.Headers["x-test-header"]
	if !ok {
		t.Fatal("Expected header to be found in config")
	}

	if settings.HMAC {
		t.Fatal("Expected HMAC to be set to false, got true")
	}

	out, err := conf.view.Get(context.Background(), auditedHeadersEntry)
	if err != nil {
		t.Fatalf("Could not retrieve headers entry from config: %s", err)
	}

	headers := make(map[string]*auditedHeaderSettings)
	err = out.DecodeJSON(&headers)
	if err != nil {
		t.Fatalf("Error decoding header view: %s", err)
	}

	expected := map[string]*auditedHeaderSettings{
		"x-test-header": &auditedHeaderSettings{
			HMAC: false,
		},
	}

	if !reflect.DeepEqual(headers, expected) {
		t.Fatalf("Expected config didn't match actual. Expected: %#v, Got: %#v", expected, headers)
	}

	err = conf.add(context.Background(), "X-Vault-Header", true)
	if err != nil {
		t.Fatalf("Error when adding header to config: %s", err)
	}

	settings, ok = conf.Headers["x-vault-header"]
	if !ok {
		t.Fatal("Expected header to be found in config")
	}

	if !settings.HMAC {
		t.Fatal("Expected HMAC to be set to true, got false")
	}

	out, err = conf.view.Get(context.Background(), auditedHeadersEntry)
	if err != nil {
		t.Fatalf("Could not retrieve headers entry from config: %s", err)
	}

	headers = make(map[string]*auditedHeaderSettings)
	err = out.DecodeJSON(&headers)
	if err != nil {
		t.Fatalf("Error decoding header view: %s", err)
	}

	expected["x-vault-header"] = &auditedHeaderSettings{
		HMAC: true,
	}

	if !reflect.DeepEqual(headers, expected) {
		t.Fatalf("Expected config didn't match actual. Expected: %#v, Got: %#v", expected, headers)
	}

}

func testAuditedHeadersConfig_Remove(t *testing.T, conf *AuditedHeadersConfig) {
	err := conf.remove(context.Background(), "X-Test-Header")
	if err != nil {
		t.Fatalf("Error when adding header to config: %s", err)
	}

	_, ok := conf.Headers["x-Test-HeAder"]
	if ok {
		t.Fatal("Expected header to not be found in config")
	}

	out, err := conf.view.Get(context.Background(), auditedHeadersEntry)
	if err != nil {
		t.Fatalf("Could not retrieve headers entry from config: %s", err)
	}

	headers := make(map[string]*auditedHeaderSettings)
	err = out.DecodeJSON(&headers)
	if err != nil {
		t.Fatalf("Error decoding header view: %s", err)
	}

	expected := map[string]*auditedHeaderSettings{
		"x-vault-header": &auditedHeaderSettings{
			HMAC: true,
		},
	}

	if !reflect.DeepEqual(headers, expected) {
		t.Fatalf("Expected config didn't match actual. Expected: %#v, Got: %#v", expected, headers)
	}

	err = conf.remove(context.Background(), "x-VaulT-Header")
	if err != nil {
		t.Fatalf("Error when adding header to config: %s", err)
	}

	_, ok = conf.Headers["x-vault-header"]
	if ok {
		t.Fatal("Expected header to not be found in config")
	}

	out, err = conf.view.Get(context.Background(), auditedHeadersEntry)
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

func TestAuditedHeadersConfig_ApplyConfig(t *testing.T) {
	conf := mockAuditedHeadersConfig(t)

	conf.add(context.Background(), "X-TesT-Header", false)
	conf.add(context.Background(), "X-Vault-HeAdEr", true)

	reqHeaders := map[string][]string{
		"X-Test-Header":  []string{"foo"},
		"X-Vault-Header": []string{"bar", "bar"},
		"Content-Type":   []string{"json"},
	}

	hashFunc := func(ctx context.Context, s string) (string, error) { return "hashed", nil }

	result, err := conf.ApplyConfig(context.Background(), reqHeaders, hashFunc)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string][]string{
		"x-test-header":  []string{"foo"},
		"x-vault-header": []string{"hashed", "hashed"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected headers did not match actual: Expected %#v\n Got %#v\n", expected, result)
	}

	//Make sure we didn't edit the reqHeaders map
	reqHeadersCopy := map[string][]string{
		"X-Test-Header":  []string{"foo"},
		"X-Vault-Header": []string{"bar", "bar"},
		"Content-Type":   []string{"json"},
	}

	if !reflect.DeepEqual(reqHeaders, reqHeadersCopy) {
		t.Fatalf("Req headers were changed, expected %#v\n got %#v", reqHeadersCopy, reqHeaders)
	}

}

func BenchmarkAuditedHeaderConfig_ApplyConfig(b *testing.B) {
	conf := &AuditedHeadersConfig{
		Headers: make(map[string]*auditedHeaderSettings),
		view:    nil,
	}

	conf.Headers = map[string]*auditedHeaderSettings{
		"X-Test-Header":  &auditedHeaderSettings{false},
		"X-Vault-Header": &auditedHeaderSettings{true},
	}

	reqHeaders := map[string][]string{
		"X-Test-Header":  []string{"foo"},
		"X-Vault-Header": []string{"bar", "bar"},
		"Content-Type":   []string{"json"},
	}

	salter, err := salt.NewSalt(context.Background(), nil, nil)
	if err != nil {
		b.Fatal(err)
	}

	hashFunc := func(ctx context.Context, s string) (string, error) { return salter.GetIdentifiedHMAC(s), nil }

	// Reset the timer since we did a lot above
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conf.ApplyConfig(context.Background(), reqHeaders, hashFunc)
	}
}
