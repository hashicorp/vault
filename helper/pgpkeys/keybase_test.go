// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pgpkeys

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

func TestFetchKeybasePubkeys(t *testing.T) {
	testset := []string{"keybase:jefferai", "keybase:hashicorp"}
	ret, err := FetchKeybasePubkeys(testset)
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
		"0f801f518ec853daff611e836528efcac6caa3db",
		"c874011f0ab405110d02105534365d9472d7468f",
	}

	if !reflect.DeepEqual(fingerprints, exp) {
		t.Fatalf("fingerprints do not match; expected \n%#v\ngot\n%#v\n", exp, fingerprints)
	}
}
