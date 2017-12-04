package vault

import (
	"reflect"
	"testing"
)

func mockPassthroughHeadersConfig(t *testing.T) *PassthroughHeadersConfig {
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "foo/")
	return &PassthroughHeadersConfig{
		Headers: make(map[string]*passthroughHeaderSettings),
		view:    view,
	}
}

func TestPassthroughHeadersConfig_CRUD(t *testing.T) {
	conf := mockPassthroughHeadersConfig(t)

	testPassthroughHeadersConfig_Add(t, conf)
	testPassthroughHeadersConfig_Remove(t, conf)
}

func testPassthroughHeadersConfig_Add(t *testing.T, conf *PassthroughHeadersConfig) {
	err := conf.add("X-Test-Header", []string{"foo"})
	if err != nil {
		t.Fatalf("Error when adding header to config: %s", err)
	}

	settings, ok := conf.Headers["x-test-header"]
	if !ok {
		t.Fatal("Expected header to be found in config")
	}

	if len(settings.Backends) != 1 || settings.Backends[0] != "foo" {
		t.Fatalf("Expected Backends to be set to [foo], got %v", settings.Backends)
	}

	out, err := conf.view.Get(passthroughHeadersEntry)
	if err != nil {
		t.Fatalf("Could not retrieve headers entry from config: %s", err)
	}

	headers := make(map[string]*passthroughHeaderSettings)
	err = out.DecodeJSON(&headers)
	if err != nil {
		t.Fatalf("Error decoding header view: %s", err)
	}

	expected := map[string]*passthroughHeaderSettings{
		"x-test-header": &passthroughHeaderSettings{
			Backends: []string{"foo"},
		},
	}

	if !reflect.DeepEqual(headers, expected) {
		t.Fatalf("Expected config didn't match actual. Expected: %#v, Got: %#v", expected, headers)
	}

	err = conf.add("X-Vault-Header", []string{"foo", "bar"})
	if err != nil {
		t.Fatalf("Error when adding header to config: %s", err)
	}

	settings, ok = conf.Headers["x-vault-header"]
	if !ok {
		t.Fatal("Expected header to be found in config")
	}

	if len(settings.Backends) != 2 {
		t.Fatalf("Expected Backends to have two values, got %v", settings.Backends)
	}

	out, err = conf.view.Get(passthroughHeadersEntry)
	if err != nil {
		t.Fatalf("Could not retrieve headers entry from config: %s", err)
	}

	headers = make(map[string]*passthroughHeaderSettings)
	err = out.DecodeJSON(&headers)
	if err != nil {
		t.Fatalf("Error decoding header view: %s", err)
	}

	expected["x-vault-header"] = &passthroughHeaderSettings{
		Backends: []string{"foo", "bar"},
	}

	if !reflect.DeepEqual(headers, expected) {
		t.Fatalf("Expected config didn't match actual. Expected: %#v, Got: %#v", expected, headers)
	}

}

func testPassthroughHeadersConfig_Remove(t *testing.T, conf *PassthroughHeadersConfig) {
	err := conf.remove("X-Test-Header")
	if err != nil {
		t.Fatalf("Error when adding header to config: %s", err)
	}

	_, ok := conf.Headers["x-Test-HeAder"]
	if ok {
		t.Fatal("Expected header to not be found in config")
	}

	_, ok = conf.Headers["x-test-header"]
	if ok {
		t.Fatal("Expected header to not be found in config")
	}

	out, err := conf.view.Get(passthroughHeadersEntry)
	if err != nil {
		t.Fatalf("Could not retrieve headers entry from config: %s", err)
	}

	headers := make(map[string]*passthroughHeaderSettings)
	err = out.DecodeJSON(&headers)
	if err != nil {
		t.Fatalf("Error decoding header view: %s", err)
	}

	expected := map[string]*passthroughHeaderSettings{
		"x-vault-header": &passthroughHeaderSettings{
			Backends: []string{"foo", "bar"},
		},
	}

	if !reflect.DeepEqual(headers, expected) {
		t.Fatalf("Expected config didn't match actual. Expected: %#v, Got: %#v", expected, headers)
	}

	err = conf.remove("x-VaulT-Header")
	if err != nil {
		t.Fatalf("Error when adding header to config: %s", err)
	}

	_, ok = conf.Headers["x-vault-header"]
	if ok {
		t.Fatal("Expected header to not be found in config")
	}

	out, err = conf.view.Get(passthroughHeadersEntry)
	if err != nil {
		t.Fatalf("Could not retrieve headers entry from config: %s", err)
	}

	headers = make(map[string]*passthroughHeaderSettings)
	err = out.DecodeJSON(&headers)
	if err != nil {
		t.Fatalf("Error decoding header view: %s", err)
	}

	expected = make(map[string]*passthroughHeaderSettings)

	if !reflect.DeepEqual(headers, expected) {
		t.Fatalf("Expected config didn't match actual. Expected: %#v, Got: %#v", expected, headers)
	}
}

func TestPassthroughHeadersConfig_ApplyConfig(t *testing.T) {
	conf := mockPassthroughHeadersConfig(t)

	conf.add("X-TesT-Header", []string{"auth/foo/"})
	conf.add("X-Vault-HeAdEr", []string{"auth/bar/"})

	reqHeaders := map[string][]string{
		"X-Test-Header":  []string{"foo"},
		"X-Vault-Header": []string{"bar", "bar"},
		"Content-Type":   []string{"json"},
	}

	result, err := conf.ApplyConfig(reqHeaders, "auth/foo/baz")
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string][]string{
		"x-test-header":  []string{"foo"},
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

func BenchmarkPassthroughHeaderConfig_ApplyConfig(b *testing.B) {
	conf := &PassthroughHeadersConfig{
		Headers: map[string]*passthroughHeaderSettings{
			"X-Test-Header":  &passthroughHeaderSettings{[]string{"foo", "bar"}},
			"X-Vault-Header": &passthroughHeaderSettings{[]string{"foo", "bar"}},
		},
		view:    nil,
	}

	reqHeaders := map[string][]string{
		"X-Test-Header":  []string{"foo"},
		"X-Vault-Header": []string{"bar", "bar"},
		"Content-Type":   []string{"json"},
	}

	// Reset the timer since we did a lot above
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conf.ApplyConfig(reqHeaders, "foo/bar")
	}
}
