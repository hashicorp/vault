package transit

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/helper/random"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestTransit_Random(t *testing.T) {
	var b *backend
	sysView := logical.TestSystemView()
	storage := &logical.InmemStorage{}
	sysView.CachingDisabledVal = true

	b, _ = Backend(context.Background(), &logical.BackendConfig{
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
			resp, err := b.HandleRequest(context.Background(), req)
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
			t.Fatal("length of output random bytes not what is expected")
		}
		if reflect.DeepEqual(rand1, rand2) {
			t.Fatal("found identical ouputs")
		}
	}

	for _, source := range []string{"", "platform", "seal", "all"} {
		req.Data["source"] = source
		req.Data["bytes"] = 32
		req.Data["format"] = "base64"
		req.Path = "random"
		// Test defaults
		doRequest(req, false, "base64", 32)

		// Test size selection in the path
		req.Path = "random/24"
		req.Data["format"] = "hex"
		doRequest(req, false, "hex", 24)

		if source != "" {
			// Test source selection in the path
			req.Path = fmt.Sprintf("random/%s", source)
			req.Data["format"] = "hex"
			doRequest(req, false, "hex", 32)

			req.Path = fmt.Sprintf("random/%s/24", source)
			req.Data["format"] = "hex"
			doRequest(req, false, "hex", 24)
		}

		// Test bad input/format
		req.Path = "random"
		req.Data["format"] = "base92"
		doRequest(req, true, "", 0)

		req.Data["format"] = "hex"
		req.Data["bytes"] = -1
		doRequest(req, true, "", 0)

		req.Data["format"] = "hex"
		req.Data["bytes"] = random.APIMaxBytes + 1

		doRequest(req, true, "", 0)
	}
}
