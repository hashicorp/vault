package command

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"reflect"
	"regexp"
	"sort"
	"testing"

	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/vault"

	"github.com/keybase/go-crypto/openpgp"
	"github.com/keybase/go-crypto/openpgp/packet"
)

func getPubKeyFiles(t *testing.T) (string, []string, error) {
	tempDir, err := ioutil.TempDir("", "vault-test")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %s", err)
	}

	pubFiles := []string{
		tempDir + "/pubkey1",
		tempDir + "/pubkey2",
		tempDir + "/pubkey3",
		tempDir + "/aapubkey1",
	}
	decoder := base64.StdEncoding
	pub1Bytes, err := decoder.DecodeString(pgpkeys.TestPubKey1)
	if err != nil {
		t.Fatalf("Error decoding bytes for public key 1: %s", err)
	}
	err = ioutil.WriteFile(pubFiles[0], pub1Bytes, 0755)
	if err != nil {
		t.Fatalf("Error writing pub key 1 to temp file: %s", err)
	}
	pub2Bytes, err := decoder.DecodeString(pgpkeys.TestPubKey2)
	if err != nil {
		t.Fatalf("Error decoding bytes for public key 2: %s", err)
	}
	err = ioutil.WriteFile(pubFiles[1], pub2Bytes, 0755)
	if err != nil {
		t.Fatalf("Error writing pub key 2 to temp file: %s", err)
	}
	pub3Bytes, err := decoder.DecodeString(pgpkeys.TestPubKey3)
	if err != nil {
		t.Fatalf("Error decoding bytes for public key 3: %s", err)
	}
	err = ioutil.WriteFile(pubFiles[2], pub3Bytes, 0755)
	if err != nil {
		t.Fatalf("Error writing pub key 3 to temp file: %s", err)
	}
	err = ioutil.WriteFile(pubFiles[3], []byte(pgpkeys.TestAAPubKey1), 0755)
	if err != nil {
		t.Fatalf("Error writing aa pub key 1 to temp file: %s", err)
	}

	return tempDir, pubFiles, nil
}

func testPGPDecrypt(tb testing.TB, privKey, enc string) string {
	tb.Helper()

	privKeyBytes, err := base64.StdEncoding.DecodeString(privKey)
	if err != nil {
		tb.Fatal(err)
	}

	ptBuf := bytes.NewBuffer(nil)
	entity, err := openpgp.ReadEntity(packet.NewReader(bytes.NewBuffer(privKeyBytes)))
	if err != nil {
		tb.Fatal(err)
	}

	var rootBytes []byte
	rootBytes, err = base64.StdEncoding.DecodeString(enc)
	if err != nil {
		tb.Fatal(err)
	}

	entityList := &openpgp.EntityList{entity}
	md, err := openpgp.ReadMessage(bytes.NewBuffer(rootBytes), entityList, nil, nil)
	if err != nil {
		tb.Fatal(err)
	}
	ptBuf.ReadFrom(md.UnverifiedBody)
	return ptBuf.String()
}

func parseDecryptAndTestUnsealKeys(t *testing.T,
	input, rootToken string,
	fingerprints bool,
	backupKeys map[string][]string,
	backupKeysB64 map[string][]string,
	core *vault.Core) {

	decoder := base64.StdEncoding
	priv1Bytes, err := decoder.DecodeString(pgpkeys.TestPrivKey1)
	if err != nil {
		t.Fatalf("Error decoding bytes for private key 1: %s", err)
	}
	priv2Bytes, err := decoder.DecodeString(pgpkeys.TestPrivKey2)
	if err != nil {
		t.Fatalf("Error decoding bytes for private key 2: %s", err)
	}
	priv3Bytes, err := decoder.DecodeString(pgpkeys.TestPrivKey3)
	if err != nil {
		t.Fatalf("Error decoding bytes for private key 3: %s", err)
	}

	privBytes := [][]byte{
		priv1Bytes,
		priv2Bytes,
		priv3Bytes,
	}

	testFunc := func(bkeys map[string][]string) {
		var re *regexp.Regexp
		if fingerprints {
			re, err = regexp.Compile("\\s*Key\\s+\\d+\\s+fingerprint:\\s+([0-9a-fA-F]+);\\s+value:\\s+(.*)")
		} else {
			re, err = regexp.Compile("\\s*Key\\s+\\d+:\\s+(.*)")
		}
		if err != nil {
			t.Fatalf("Error compiling regex: %s", err)
		}
		matches := re.FindAllStringSubmatch(input, -1)
		if len(matches) != 4 {
			t.Fatalf("Unexpected number of keys returned, got %d, matches was \n\n%#v\n\n, input was \n\n%s\n\n", len(matches), matches, input)
		}

		encodedKeys := []string{}
		matchedFingerprints := []string{}
		for _, tuple := range matches {
			if fingerprints {
				if len(tuple) != 3 {
					t.Fatalf("Key not found: %#v", tuple)
				}
				matchedFingerprints = append(matchedFingerprints, tuple[1])
				encodedKeys = append(encodedKeys, tuple[2])
			} else {
				if len(tuple) != 2 {
					t.Fatalf("Key not found: %#v", tuple)
				}
				encodedKeys = append(encodedKeys, tuple[1])
			}
		}

		if bkeys != nil && len(matchedFingerprints) != 0 {
			testMap := map[string][]string{}
			for i, v := range matchedFingerprints {
				testMap[v] = append(testMap[v], encodedKeys[i])
				sort.Strings(testMap[v])
			}
			if !reflect.DeepEqual(testMap, bkeys) {
				t.Fatalf("test map and backup map do not match, test map is\n%#v\nbackup map is\n%#v", testMap, bkeys)
			}
		}

		unsealKeys := []string{}
		ptBuf := bytes.NewBuffer(nil)
		for i, privKeyBytes := range privBytes {
			if i > 2 {
				break
			}
			ptBuf.Reset()
			entity, err := openpgp.ReadEntity(packet.NewReader(bytes.NewBuffer(privKeyBytes)))
			if err != nil {
				t.Fatalf("Error parsing private key %d: %s", i, err)
			}
			var keyBytes []byte
			keyBytes, err = base64.StdEncoding.DecodeString(encodedKeys[i])
			if err != nil {
				t.Fatalf("Error decoding key %d: %s", i, err)
			}
			entityList := &openpgp.EntityList{entity}
			md, err := openpgp.ReadMessage(bytes.NewBuffer(keyBytes), entityList, nil, nil)
			if err != nil {
				t.Fatalf("Error decrypting with key %d (%s): %s", i, encodedKeys[i], err)
			}
			ptBuf.ReadFrom(md.UnverifiedBody)
			unsealKeys = append(unsealKeys, ptBuf.String())
		}

		err = core.Seal(rootToken)
		if err != nil {
			t.Fatalf("Error sealing vault with provided root token: %s", err)
		}

		for i, unsealKey := range unsealKeys {
			unsealBytes, err := hex.DecodeString(unsealKey)
			if err != nil {
				t.Fatalf("Error hex decoding unseal key %s: %s", unsealKey, err)
			}
			unsealed, err := core.Unseal(unsealBytes)
			if err != nil {
				t.Fatalf("Error using unseal key %s: %s", unsealKey, err)
			}
			if i >= 2 && !unsealed {
				t.Fatalf("Error: Provided two unseal keys but core is not unsealed")
			}
		}
	}

	testFunc(backupKeysB64)
}
