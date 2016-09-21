package transit

import (
	"encoding/base64"
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestTransit_Random(t *testing.T) {
	var b *backend
	sysView := logical.TestSystemView()
	storage := &logical.InmemStorage{}

	b = Backend(&logical.BackendConfig{
		StorageView: storage,
		System:      sysView,
	})

	req := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "random",
		Data:      map[string]interface{}{},
	}

	doRequest := func(req *logical.Request, errExpected bool, format string, numBytes int) {
		getResponse := func() []byte {
			resp, err := b.HandleRequest(req)
			if err != nil && !errExpected {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("expected non-nil response")
			}
			if errExpected {
				if !resp.IsError() {
					t.Fatalf("bad: got error response: %#v", *resp)
				}
				return nil
			}
			if resp.IsError() {
				t.Fatalf("bad: got error response: %#v", *resp)
			}
			if _, ok := resp.Data["random_bytes"]; !ok {
				t.Fatal("no random_bytes found in response")
			}

			outputStr := resp.Data["random_bytes"].(string)
			var outputBytes []byte
			switch format {
			case "base64":
				outputBytes, err = base64.StdEncoding.DecodeString(outputStr)
			case "hex":
				outputBytes, err = hex.DecodeString(outputStr)
			default:
				t.Fatal("unknown format")
			}
			if err != nil {
				t.Fatal(err)
			}

			return outputBytes
		}

		rand1 := getResponse()
		// Expected error
		if rand1 == nil {
			return
		}
		rand2 := getResponse()
		if len(rand1) != numBytes || len(rand2) != numBytes {
			t.Fatal("length of output random bytes not what is exepcted")
		}
		if reflect.DeepEqual(rand1, rand2) {
			t.Fatal("found identical ouputs")
		}
	}

	// Test defaults
	doRequest(req, false, "base64", 32)

	// Test size selection in the path
	req.Path = "random/24"
	req.Data["format"] = "hex"
	doRequest(req, false, "hex", 24)

	// Test bad input/format
	req.Path = "random"
	req.Data["format"] = "base92"
	doRequest(req, true, "", 0)

	req.Data["format"] = "hex"
	req.Data["bytes"] = -1
	doRequest(req, true, "", 0)
}
