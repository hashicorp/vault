package http

import (
	"encoding/hex"
	"net/http"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/vault"
)

func TestSysInit_get(t *testing.T) {
	core := vault.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	{
		// Pre-init
		resp, err := http.Get(addr + "/v1/sys/init")
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		var actual map[string]interface{}
		expected := map[string]interface{}{
			"initialized": false,
		}
		testResponseStatus(t, resp, 200)
		testResponseBody(t, resp, &actual)
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("bad: %#v", actual)
		}
	}

	vault.TestCoreInit(t, core)

	{
		// Post-init
		resp, err := http.Get(addr + "/v1/sys/init")
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		var actual map[string]interface{}
		expected := map[string]interface{}{
			"initialized": true,
		}
		testResponseStatus(t, resp, 200)
		testResponseBody(t, resp, &actual)
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("bad: %#v", actual)
		}
	}
}

// Test to check if the API errors out when wrong number of PGP keys are
// supplied
func TestSysInit_pgpKeysEntries(t *testing.T) {
	core := vault.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp := testHttpPut(t, "", addr+"/v1/sys/init", map[string]interface{}{
		"secret_shares":   5,
		"secret_threhold": 3,
		"pgp_keys":        []string{"pgpkey1"},
	})
	testResponseStatus(t, resp, 400)
}

// Test to check if the API errors out when wrong number of PGP keys are
// supplied for recovery config
func TestSysInit_pgpKeysEntriesForRecovery(t *testing.T) {
	core := vault.TestCoreNewSeal(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp := testHttpPut(t, "", addr+"/v1/sys/init", map[string]interface{}{
		"secret_shares":      1,
		"secret_threshold":   1,
		"stored_shares":      1,
		"recovery_shares":    5,
		"recovery_threshold": 3,
		"recovery_pgp_keys":  []string{"pgpkey1"},
	})
	testResponseStatus(t, resp, 400)
}

func TestSysInit_SecretSharesMetadata(t *testing.T) {
	core := vault.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	var actual map[string]interface{}

	// basic case of getting unseal metadata back
	initInput := map[string]interface{}{
		"secret_shares":    5,
		"secret_threshold": 3,
	}
	resp := testHttpPut(t, "", addr+"/v1/sys/init", initInput)
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	secretSharesMetadata := actual["secret_shares_metadata"].([]interface{})
	if len(secretSharesMetadata) != 5 {
		t.Fatalf("bad: length of secret_shares_metadata: expected: 5, actual: %d", len(secretSharesMetadata))
	}

	for _, item := range secretSharesMetadata {
		metadata := item.(map[string]interface{})
		if metadata["id"].(string) == "" {
			t.Fatalf("bad: missing identifier in secretSharesMetadata: %#v", secretSharesMetadata)
		}
	}
}

// Test key identifiers to have valid IDs
func TestSysInit_KeyIdentifiersCase1(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	var actual InitKeySharesIdentifiersResponse
	resp := testHttpGet(t, token, addr+"/v1/sys/init/unseal-shares-identifiers")
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	if len(actual.KeyIdentifiers) != 3 {
		t.Fatalf("bad: number of key identifiers; expected: 5, actual: %d", len(actual.KeyIdentifiers))
	}

	for _, identifier := range actual.KeyIdentifiers {
		if identifier.ID == "" {
			t.Fatalf("bad: key identifier is empty")
		}
	}
}

// Test key identifiers to have valid IDs and names
func TestSysInit_KeyIdentifiersCase2(t *testing.T) {
	bc, rc := vault.TestSealDefConfigs()
	core, _, _, token, keyIdentifiers := vault.TestCoreUnsealedWithConfigs(t, bc, rc)

	if len(keyIdentifiers) != 5 {
		t.Fatalf("bad: number of key identifiers; expected: 5, actual: %d", len(keyIdentifiers))
	}

	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	var actual InitKeySharesIdentifiersResponse
	resp := testHttpGet(t, token, addr+"/v1/sys/init/unseal-shares-identifiers")
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	nameList := []string{"first", "second", "third", "forth", "fifth"}

	for _, identifier := range actual.KeyIdentifiers {
		if identifier.ID == "" {
			t.Fatalf("bad: key identifier is empty")
		}
		if identifier.Name == "" {
			t.Fatalf("bad: key identifier name is empty")
		}
		nameList = strutil.StrListDelete(nameList, identifier.Name)
	}

	if len(nameList) != 0 {
		t.Fatalf("bad: length of key identifier names list: expected 0, actual: %d", len(nameList))
	}
}

func TestSysInit_SecretSharesMetadataKeyIdentifierNames(t *testing.T) {
	core := vault.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	var actual map[string]interface{}

	// basic case of getting unseal metadata back
	initInput := map[string]interface{}{
		"secret_shares":                  5,
		"secret_threshold":               3,
		"secret_shares_identifier_names": "first,second,third,forth,fifth",
	}

	// set the key identifier names and check if the associated metadata has
	// names in it
	resp := testHttpPut(t, "", addr+"/v1/sys/init", initInput)
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	secretSharesMetadata := actual["secret_shares_metadata"].([]interface{})
	nameList := []string{"first", "second", "third", "forth", "fifth"}

	if len(secretSharesMetadata) != 5 {
		t.Fatalf("bad: length of secret_shares_metadata: expected: 5, actual: %d", len(secretSharesMetadata))
	}

	for _, item := range secretSharesMetadata {
		metadata := item.(map[string]interface{})
		if metadata["id"].(string) == "" {
			t.Fatalf("bad: missing identifier in secretSharesMetadata: %#v", secretSharesMetadata)
		}
		if metadata["name"].(string) == "" {
			t.Fatalf("invalid key identifier name")
		}
		nameList = strutil.StrListDelete(nameList, metadata["name"].(string))
	}

	if len(nameList) != 0 {
		t.Fatalf("bad: length of key identifier names list: expected 0, actual: %d", len(nameList))
	}
}

func TestSysInit_put(t *testing.T) {
	core := vault.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp := testHttpPut(t, "", addr+"/v1/sys/init", map[string]interface{}{
		"secret_shares":    5,
		"secret_threshold": 3,
	})

	var actual map[string]interface{}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	keysRaw, ok := actual["keys"]
	if !ok {
		t.Fatalf("no keys: %#v", actual)
	}

	if _, ok := actual["root_token"]; !ok {
		t.Fatal("no root token")
	}

	for _, key := range keysRaw.([]interface{}) {
		keySlice, err := hex.DecodeString(key.(string))
		if err != nil {
			t.Fatalf("bad: %s", err)
		}

		if _, err := core.Unseal(keySlice); err != nil {
			t.Fatalf("bad: %s", err)
		}
	}

	seal, err := core.Sealed()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if seal {
		t.Fatal("should not be sealed")
	}
}
