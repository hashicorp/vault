package pgpkeys

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"reflect"
	"testing"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
)

func TestFetchEmptyKeyserver(t *testing.T) {
	testset := []string{"foo", "bar"}
	if _, err := FetchKeyserverPubkeys("", testset); err == nil {
		t.Fatal("expected to fail with empty keyserver")
	}
}

func TestFetchKeyserverPubkeys(t *testing.T) {
	testset := []string{"jess@docker.com", "paultag@debian.org"}
	ret, err := FetchKeyserverPubkeys(MITKeyserver, testset)
	if err != nil {
		t.Fatalf("bad: %v", err)
	}

	fingerprints := []string{}
	for _, user := range testset {
		data, err := base64.StdEncoding.DecodeString(ret[user])
		if err != nil {
			t.Fatalf("error decoding key for user %s: %v", user, err)
		}
		entity, err := openpgp.ReadEntity(packet.NewReader(bytes.NewBuffer(data)))
		if err != nil {
			t.Fatalf("error parsing key for user %s: %v", user, err)
		}
		fingerprints = append(fingerprints, hex.EncodeToString(entity.PrimaryKey.Fingerprint[:]))
	}

	exp := []string{
		"d4c4dd600d66f65a8efc511e18f3685c0022bff3",
		"8f049ad82c92066c7352d28a7b585b30807c2a87",
	}

	if !reflect.DeepEqual(fingerprints, exp) {
		t.Fatalf("fingerprints do not match; expected \n%#v\ngot\n%#v\n", exp, fingerprints)
	}
}
