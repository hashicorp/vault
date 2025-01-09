// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keysutil

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	mathrand "math/rand"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/helper/cryptoutil"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/copystructure"
	"golang.org/x/crypto/ed25519"
)

// Ordering of these items needs to match the iota order defined in policy.go. Ordering changes
// should never occur, as it would lead to a key type change within existing stored policies.
var allTestKeyTypes = []KeyType{
	KeyType_AES256_GCM96, KeyType_ECDSA_P256, KeyType_ED25519, KeyType_RSA2048,
	KeyType_RSA4096, KeyType_ChaCha20_Poly1305, KeyType_ECDSA_P384, KeyType_ECDSA_P521, KeyType_AES128_GCM96,
	KeyType_RSA3072, KeyType_MANAGED_KEY, KeyType_HMAC, KeyType_AES128_CMAC, KeyType_AES256_CMAC, KeyType_ML_DSA,
	KeyType_HYBRID,
}

func TestPolicy_KeyTypes(t *testing.T) {
	// Make sure the iota value never change for key types, as existing storage would be affected
	for i, keyType := range allTestKeyTypes {
		if int(keyType) != i {
			t.Fatalf("iota of keytype %s changed, expected %d got %d", keyType.String(), i, keyType)
		}
	}

	// Make sure we have a string presentation for all types
	for _, keyType := range allTestKeyTypes {
		if strings.Contains(keyType.String(), "unknown") {
			t.Fatalf("keytype with iota of %d should not contain 'unknown', missing in String() switch statement", keyType)
		}
	}
}

func TestPolicy_HmacCmacSupported(t *testing.T) {
	// Test HMAC supported feature
	for _, keyType := range allTestKeyTypes {
		switch keyType {
		case KeyType_MANAGED_KEY:
			if keyType.HMACSupported() {
				t.Fatalf("hmac should not have been not be supported for keytype %s", keyType.String())
			}
			if keyType.CMACSupported() {
				t.Fatalf("cmac should not have been be supported for keytype %s", keyType.String())
			}
		case KeyType_AES128_CMAC, KeyType_AES256_CMAC:
			if keyType.HMACSupported() {
				t.Fatalf("hmac should have been not be supported for keytype %s", keyType.String())
			}
			if !keyType.CMACSupported() {
				t.Fatalf("cmac should have been be supported for keytype %s", keyType.String())
			}
		default:
			if !keyType.HMACSupported() {
				t.Fatalf("hmac should have been supported for keytype %s", keyType.String())
			}
			if keyType.CMACSupported() {
				t.Fatalf("cmac should not have been supported for keytype %s", keyType.String())
			}
		}
	}
}

func TestPolicy_CMACKeyUpgrade(t *testing.T) {
	ctx := context.Background()
	lm, _ := NewLockManager(false, 0)
	storage := &logical.InmemStorage{}
	p, upserted, err := lm.GetPolicy(ctx, PolicyRequest{
		Upsert:  true,
		Storage: storage,
		KeyType: KeyType_AES256_CMAC,
		Name:    "test",
	}, rand.Reader)
	if err != nil {
		t.Fatalf("failed loading policy: %v", err)
	}
	if p == nil {
		t.Fatal("nil policy")
	}
	if !upserted {
		t.Fatal("expected an upsert")
	}

	// This verifies we don't have a hmac key
	_, err = p.HMACKey(1)
	if err == nil {
		t.Fatal("cmac key should not return an hmac key but did on initial creation")
	}

	if p.NeedsUpgrade() {
		t.Fatal("cmac key should not require an upgrade after initial key creation")
	}

	err = p.Upgrade(ctx, storage, rand.Reader)
	if err != nil {
		t.Fatalf("an error was returned from upgrade method: %v", err)
	}
	p.Unlock()

	// Now reload our policy from disk and make sure we still don't have a hmac key
	p, upserted, err = lm.GetPolicy(ctx, PolicyRequest{
		Upsert:  true,
		Storage: storage,
		KeyType: KeyType_AES256_CMAC,
		Name:    "test",
	}, rand.Reader)
	if err != nil {
		t.Fatalf("failed loading policy: %v", err)
	}
	if p == nil {
		t.Fatal("nil policy")
	}
	if upserted {
		t.Fatal("expected the key to exist but upserted was true")
	}

	p.Unlock()

	_, err = p.HMACKey(1)
	if err == nil {
		t.Fatal("cmac key should not return an hmac key post upgrade")
	}
}

func TestPolicy_KeyEntryMapUpgrade(t *testing.T) {
	now := time.Now()
	old := map[int]KeyEntry{
		1: {
			Key:                []byte("samplekey"),
			HMACKey:            []byte("samplehmackey"),
			CreationTime:       now,
			FormattedPublicKey: "sampleformattedpublickey",
		},
		2: {
			Key:                []byte("samplekey2"),
			HMACKey:            []byte("samplehmackey2"),
			CreationTime:       now.Add(10 * time.Second),
			FormattedPublicKey: "sampleformattedpublickey2",
		},
	}

	oldEncoded, err := jsonutil.EncodeJSON(old)
	if err != nil {
		t.Fatal(err)
	}

	var new keyEntryMap
	err = jsonutil.DecodeJSON(oldEncoded, &new)
	if err != nil {
		t.Fatal(err)
	}

	newEncoded, err := jsonutil.EncodeJSON(&new)
	if err != nil {
		t.Fatal(err)
	}

	if string(oldEncoded) != string(newEncoded) {
		t.Fatalf("failed to upgrade key entry map;\nold: %q\nnew: %q", string(oldEncoded), string(newEncoded))
	}
}

func Test_KeyUpgrade(t *testing.T) {
	lockManagerWithCache, _ := NewLockManager(true, 0)
	lockManagerWithoutCache, _ := NewLockManager(false, 0)
	testKeyUpgradeCommon(t, lockManagerWithCache)
	testKeyUpgradeCommon(t, lockManagerWithoutCache)
}

func testKeyUpgradeCommon(t *testing.T, lm *LockManager) {
	ctx := context.Background()

	storage := &logical.InmemStorage{}
	p, upserted, err := lm.GetPolicy(ctx, PolicyRequest{
		Upsert:  true,
		Storage: storage,
		KeyType: KeyType_AES256_GCM96,
		Name:    "test",
	}, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	if p == nil {
		t.Fatal("nil policy")
	}
	if !upserted {
		t.Fatal("expected an upsert")
	}
	if !lm.useCache {
		p.Unlock()
	}

	testBytes := make([]byte, len(p.Keys["1"].Key))
	copy(testBytes, p.Keys["1"].Key)

	p.Key = p.Keys["1"].Key
	p.Keys = nil
	p.MigrateKeyToKeysMap()
	if p.Key != nil {
		t.Fatal("policy.Key is not nil")
	}
	if len(p.Keys) != 1 {
		t.Fatal("policy.Keys is the wrong size")
	}
	if !reflect.DeepEqual(testBytes, p.Keys["1"].Key) {
		t.Fatal("key mismatch")
	}
}

func Test_ArchivingUpgrade(t *testing.T) {
	lockManagerWithCache, _ := NewLockManager(true, 0)
	lockManagerWithoutCache, _ := NewLockManager(false, 0)
	testArchivingUpgradeCommon(t, lockManagerWithCache)
	testArchivingUpgradeCommon(t, lockManagerWithoutCache)
}

func testArchivingUpgradeCommon(t *testing.T, lm *LockManager) {
	ctx := context.Background()

	// First, we generate a policy and rotate it a number of times. Each time
	// we'll ensure that we have the expected number of keys in the archive and
	// the main keys object, which without changing the min version should be
	// zero and latest, respectively

	storage := &logical.InmemStorage{}
	p, _, err := lm.GetPolicy(ctx, PolicyRequest{
		Upsert:  true,
		Storage: storage,
		KeyType: KeyType_AES256_GCM96,
		Name:    "test",
	}, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	if p == nil {
		t.Fatal("nil policy")
	}
	if !lm.useCache {
		p.Unlock()
	}

	// Store the initial key in the archive
	keysArchive := []KeyEntry{{}, p.Keys["1"]}
	checkKeys(t, ctx, p, storage, keysArchive, "initial", 1, 1, 1)

	for i := 2; i <= 10; i++ {
		err = p.Rotate(ctx, storage, rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		keysArchive = append(keysArchive, p.Keys[strconv.Itoa(i)])
		checkKeys(t, ctx, p, storage, keysArchive, "rotate", i, i, i)
	}

	// Now, wipe the archive and set the archive version to zero
	err = storage.Delete(ctx, "archive/test")
	if err != nil {
		t.Fatal(err)
	}
	p.ArchiveVersion = 0

	// Store it, but without calling persist, so we don't trigger
	// handleArchiving()
	buf, err := p.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	// Write the policy into storage
	err = storage.Put(ctx, &logical.StorageEntry{
		Key:   "policy/" + p.Name,
		Value: buf,
	})
	if err != nil {
		t.Fatal(err)
	}

	// If we're caching, expire from the cache since we modified it
	// under-the-hood
	if lm.useCache {
		lm.cache.Delete("test")
	}

	// Now get the policy again; the upgrade should happen automatically
	p, _, err = lm.GetPolicy(ctx, PolicyRequest{
		Storage: storage,
		Name:    "test",
	}, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	if p == nil {
		t.Fatal("nil policy")
	}
	if !lm.useCache {
		p.Unlock()
	}

	checkKeys(t, ctx, p, storage, keysArchive, "upgrade", 10, 10, 10)

	// Let's check some deletion logic while we're at it

	// The policy should be in there
	if lm.useCache {
		_, ok := lm.cache.Load("test")
		if !ok {
			t.Fatal("nil policy in cache")
		}
	}

	// First we'll do this wrong, by not setting the deletion flag
	err = lm.DeletePolicy(ctx, storage, "test")
	if err == nil {
		t.Fatal("got nil error, but should not have been able to delete since we didn't set the deletion flag on the policy")
	}

	// The policy should still be in there
	if lm.useCache {
		_, ok := lm.cache.Load("test")
		if !ok {
			t.Fatal("nil policy in cache")
		}
	}

	p, _, err = lm.GetPolicy(ctx, PolicyRequest{
		Storage: storage,
		Name:    "test",
	}, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	if p == nil {
		t.Fatal("policy nil after bad delete")
	}
	if !lm.useCache {
		p.Unlock()
	}

	// Now do it properly
	p.DeletionAllowed = true
	err = p.Persist(ctx, storage)
	if err != nil {
		t.Fatal(err)
	}
	err = lm.DeletePolicy(ctx, storage, "test")
	if err != nil {
		t.Fatal(err)
	}

	// The policy should *not* be in there
	if lm.useCache {
		_, ok := lm.cache.Load("test")
		if ok {
			t.Fatal("non-nil policy in cache")
		}
	}

	p, _, err = lm.GetPolicy(ctx, PolicyRequest{
		Storage: storage,
		Name:    "test",
	}, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	if p != nil {
		t.Fatal("policy not nil after delete")
	}
}

func Test_Archiving(t *testing.T) {
	lockManagerWithCache, _ := NewLockManager(true, 0)
	lockManagerWithoutCache, _ := NewLockManager(false, 0)
	testArchivingUpgradeCommon(t, lockManagerWithCache)
	testArchivingUpgradeCommon(t, lockManagerWithoutCache)
}

func testArchivingCommon(t *testing.T, lm *LockManager) {
	ctx := context.Background()

	// First, we generate a policy and rotate it a number of times. Each time
	// we'll ensure that we have the expected number of keys in the archive and
	// the main keys object, which without changing the min version should be
	// zero and latest, respectively

	storage := &logical.InmemStorage{}
	p, _, err := lm.GetPolicy(ctx, PolicyRequest{
		Upsert:  true,
		Storage: storage,
		KeyType: KeyType_AES256_GCM96,
		Name:    "test",
	}, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	if p == nil {
		t.Fatal("nil policy")
	}
	if !lm.useCache {
		p.Unlock()
	}

	// Store the initial key in the archive
	keysArchive := []KeyEntry{{}, p.Keys["1"]}
	checkKeys(t, ctx, p, storage, keysArchive, "initial", 1, 1, 1)

	for i := 2; i <= 10; i++ {
		err = p.Rotate(ctx, storage, rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		keysArchive = append(keysArchive, p.Keys[strconv.Itoa(i)])
		checkKeys(t, ctx, p, storage, keysArchive, "rotate", i, i, i)
	}

	// Move the min decryption version up
	for i := 1; i <= 10; i++ {
		p.MinDecryptionVersion = i

		err = p.Persist(ctx, storage)
		if err != nil {
			t.Fatal(err)
		}
		// We expect to find:
		// * The keys in archive are the same as the latest version
		// * The latest version is constant
		// * The number of keys in the policy itself is from the min
		// decryption version up to the latest version, so for e.g. 7 and
		// 10, you'd need 7, 8, 9, and 10 -- IOW, latest version - min
		// decryption version plus 1 (the min decryption version key
		// itself)
		checkKeys(t, ctx, p, storage, keysArchive, "minadd", 10, 10, p.LatestVersion-p.MinDecryptionVersion+1)
	}

	// Move the min decryption version down
	for i := 10; i >= 1; i-- {
		p.MinDecryptionVersion = i

		err = p.Persist(ctx, storage)
		if err != nil {
			t.Fatal(err)
		}
		// We expect to find:
		// * The keys in archive are never removed so same as the latest version
		// * The latest version is constant
		// * The number of keys in the policy itself is from the min
		// decryption version up to the latest version, so for e.g. 7 and
		// 10, you'd need 7, 8, 9, and 10 -- IOW, latest version - min
		// decryption version plus 1 (the min decryption version key
		// itself)
		checkKeys(t, ctx, p, storage, keysArchive, "minsub", 10, 10, p.LatestVersion-p.MinDecryptionVersion+1)
	}
}

func checkKeys(t *testing.T,
	ctx context.Context,
	p *Policy,
	storage logical.Storage,
	keysArchive []KeyEntry,
	action string,
	archiveVer, latestVer, keysSize int,
) {
	// Sanity check
	if len(keysArchive) != latestVer+1 {
		t.Fatalf("latest expected key version is %d, expected test keys archive size is %d, "+
			"but keys archive is of size %d", latestVer, latestVer+1, len(keysArchive))
	}

	archive, err := p.LoadArchive(ctx, storage)
	if err != nil {
		t.Fatal(err)
	}

	badArchiveVer := false
	if archiveVer == 0 {
		if len(archive.Keys) != 0 || p.ArchiveVersion != 0 {
			badArchiveVer = true
		}
	} else {
		// We need to subtract one because we have the indexes match key
		// versions, which start at 1. So for an archive version of 1, we
		// actually have two entries -- a blank 0 entry, and the key at spot 1
		if archiveVer != len(archive.Keys)-1 || archiveVer != p.ArchiveVersion {
			badArchiveVer = true
		}
	}
	if badArchiveVer {
		t.Fatalf(
			"expected archive version %d, found length of archive keys %d and policy archive version %d",
			archiveVer, len(archive.Keys), p.ArchiveVersion,
		)
	}

	if latestVer != p.LatestVersion {
		t.Fatalf(
			"expected latest version %d, found %d",
			latestVer, p.LatestVersion,
		)
	}

	if keysSize != len(p.Keys) {
		t.Fatalf(
			"expected keys size %d, found %d, action is %s, policy is \n%#v\n",
			keysSize, len(p.Keys), action, p,
		)
	}

	for i := p.MinDecryptionVersion; i <= p.LatestVersion; i++ {
		if _, ok := p.Keys[strconv.Itoa(i)]; !ok {
			t.Fatalf(
				"expected key %d, did not find it in policy keys", i,
			)
		}
	}

	for i := p.MinDecryptionVersion; i <= p.LatestVersion; i++ {
		ver := strconv.Itoa(i)
		if !p.Keys[ver].CreationTime.Equal(keysArchive[i].CreationTime) {
			t.Fatalf("key %d not equivalent between policy keys and test keys archive; policy keys:\n%#v\ntest keys archive:\n%#v\n", i, p.Keys[ver], keysArchive[i])
		}
		polKey := p.Keys[ver]
		polKey.CreationTime = keysArchive[i].CreationTime
		p.Keys[ver] = polKey
		if !reflect.DeepEqual(p.Keys[ver], keysArchive[i]) {
			t.Fatalf("key %d not equivalent between policy keys and test keys archive; policy keys:\n%#v\ntest keys archive:\n%#v\n", i, p.Keys[ver], keysArchive[i])
		}
	}

	for i := 1; i < len(archive.Keys); i++ {
		if !reflect.DeepEqual(archive.Keys[i].Key, keysArchive[i].Key) {
			t.Fatalf("key %d not equivalent between policy archive and test keys archive; policy archive:\n%#v\ntest keys archive:\n%#v\n", i, archive.Keys[i].Key, keysArchive[i].Key)
		}
	}
}

func Test_StorageErrorSafety(t *testing.T) {
	ctx := context.Background()
	lm, _ := NewLockManager(true, 0)

	storage := &logical.InmemStorage{}
	p, _, err := lm.GetPolicy(ctx, PolicyRequest{
		Upsert:  true,
		Storage: storage,
		KeyType: KeyType_AES256_GCM96,
		Name:    "test",
	}, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	if p == nil {
		t.Fatal("nil policy")
	}

	// Store the initial key in the archive
	keysArchive := []KeyEntry{{}, p.Keys["1"]}
	checkKeys(t, ctx, p, storage, keysArchive, "initial", 1, 1, 1)

	// We use checkKeys here just for sanity; it doesn't really handle cases of
	// errors below so we do more targeted testing later
	for i := 2; i <= 5; i++ {
		err = p.Rotate(ctx, storage, rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		keysArchive = append(keysArchive, p.Keys[strconv.Itoa(i)])
		checkKeys(t, ctx, p, storage, keysArchive, "rotate", i, i, i)
	}

	underlying := storage.Underlying()
	underlying.FailPut(true)

	priorLen := len(p.Keys)

	err = p.Rotate(ctx, storage, rand.Reader)
	if err == nil {
		t.Fatal("expected error")
	}

	if len(p.Keys) != priorLen {
		t.Fatal("length of keys should not have changed")
	}
}

func Test_BadUpgrade(t *testing.T) {
	ctx := context.Background()
	lm, _ := NewLockManager(true, 0)
	storage := &logical.InmemStorage{}
	p, _, err := lm.GetPolicy(ctx, PolicyRequest{
		Upsert:  true,
		Storage: storage,
		KeyType: KeyType_AES256_GCM96,
		Name:    "test",
	}, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	if p == nil {
		t.Fatal("nil policy")
	}

	orig, err := copystructure.Copy(p)
	if err != nil {
		t.Fatal(err)
	}
	orig.(*Policy).l = p.l

	p.Key = p.Keys["1"].Key
	p.Keys = nil
	p.MinDecryptionVersion = 0

	if err := p.Upgrade(ctx, storage, rand.Reader); err != nil {
		t.Fatal(err)
	}

	k := p.Keys["1"]
	o := orig.(*Policy).Keys["1"]
	k.CreationTime = o.CreationTime
	k.HMACKey = o.HMACKey
	p.Keys["1"] = k
	p.versionPrefixCache = sync.Map{}

	if !reflect.DeepEqual(orig, p) {
		t.Fatalf("not equal:\n%#v\n%#v", orig, p)
	}

	// Do it again with a failing storage call
	underlying := storage.Underlying()
	underlying.FailPut(true)

	p.Key = p.Keys["1"].Key
	p.Keys = nil
	p.MinDecryptionVersion = 0

	if err := p.Upgrade(ctx, storage, rand.Reader); err == nil {
		t.Fatal("expected error")
	}

	if p.MinDecryptionVersion == 1 {
		t.Fatal("min decryption version was changed")
	}
	if p.Keys != nil {
		t.Fatal("found upgraded keys")
	}
	if p.Key == nil {
		t.Fatal("non-upgraded key not found")
	}
}

func Test_BadArchive(t *testing.T) {
	ctx := context.Background()
	lm, _ := NewLockManager(true, 0)
	storage := &logical.InmemStorage{}
	p, _, err := lm.GetPolicy(ctx, PolicyRequest{
		Upsert:  true,
		Storage: storage,
		KeyType: KeyType_AES256_GCM96,
		Name:    "test",
	}, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	if p == nil {
		t.Fatal("nil policy")
	}

	for i := 2; i <= 10; i++ {
		err = p.Rotate(ctx, storage, rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
	}

	p.MinDecryptionVersion = 5
	if err := p.Persist(ctx, storage); err != nil {
		t.Fatal(err)
	}
	if p.ArchiveVersion != 10 {
		t.Fatalf("unexpected archive version %d", p.ArchiveVersion)
	}
	if len(p.Keys) != 6 {
		t.Fatalf("unexpected key length %d", len(p.Keys))
	}

	// Set back
	p.MinDecryptionVersion = 1
	if err := p.Persist(ctx, storage); err != nil {
		t.Fatal(err)
	}
	if p.ArchiveVersion != 10 {
		t.Fatalf("unexpected archive version %d", p.ArchiveVersion)
	}
	if len(p.Keys) != 10 {
		t.Fatalf("unexpected key length %d", len(p.Keys))
	}

	// Run it again but we'll turn off storage along the way
	p.MinDecryptionVersion = 5
	if err := p.Persist(ctx, storage); err != nil {
		t.Fatal(err)
	}
	if p.ArchiveVersion != 10 {
		t.Fatalf("unexpected archive version %d", p.ArchiveVersion)
	}
	if len(p.Keys) != 6 {
		t.Fatalf("unexpected key length %d", len(p.Keys))
	}

	underlying := storage.Underlying()
	underlying.FailPut(true)

	// Set back, which should cause p.Keys to be changed if the persist works,
	// but it doesn't
	p.MinDecryptionVersion = 1
	if err := p.Persist(ctx, storage); err == nil {
		t.Fatal("expected error during put")
	}
	if p.ArchiveVersion != 10 {
		t.Fatalf("unexpected archive version %d", p.ArchiveVersion)
	}
	// Here's the expected change
	if len(p.Keys) != 6 {
		t.Fatalf("unexpected key length %d", len(p.Keys))
	}
}

func Test_Import(t *testing.T) {
	ctx := context.Background()
	storage := &logical.InmemStorage{}
	testKeys, err := generateTestKeys()
	if err != nil {
		t.Fatalf("error generating test keys: %s", err)
	}

	tests := map[string]struct {
		policy      Policy
		key         []byte
		shouldError bool
	}{
		"import AES key": {
			policy: Policy{
				Name: "test-aes-key",
				Type: KeyType_AES256_GCM96,
			},
			key:         testKeys[KeyType_AES256_GCM96],
			shouldError: false,
		},
		"import RSA key": {
			policy: Policy{
				Name: "test-rsa-key",
				Type: KeyType_RSA2048,
			},
			key:         testKeys[KeyType_RSA2048],
			shouldError: false,
		},
		"import ECDSA key": {
			policy: Policy{
				Name: "test-ecdsa-key",
				Type: KeyType_ECDSA_P256,
			},
			key:         testKeys[KeyType_ECDSA_P256],
			shouldError: false,
		},
		"import ED25519 key": {
			policy: Policy{
				Name: "test-ed25519-key",
				Type: KeyType_ED25519,
			},
			key:         testKeys[KeyType_ED25519],
			shouldError: false,
		},
		"import incorrect key type": {
			policy: Policy{
				Name: "test-ed25519-key",
				Type: KeyType_ED25519,
			},
			key:         testKeys[KeyType_AES256_GCM96],
			shouldError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if err := test.policy.Import(ctx, storage, test.key, rand.Reader); (err != nil) != test.shouldError {
				t.Fatalf("error importing key: %s", err)
			}
		})
	}
}

func generateTestKeys() (map[KeyType][]byte, error) {
	keyMap := make(map[KeyType][]byte)

	rsaKey, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	rsaKeyBytes, err := x509.MarshalPKCS8PrivateKey(rsaKey)
	if err != nil {
		return nil, err
	}
	keyMap[KeyType_RSA2048] = rsaKeyBytes

	rsaKey, err = cryptoutil.GenerateRSAKey(rand.Reader, 3072)
	if err != nil {
		return nil, err
	}
	rsaKeyBytes, err = x509.MarshalPKCS8PrivateKey(rsaKey)
	if err != nil {
		return nil, err
	}
	keyMap[KeyType_RSA3072] = rsaKeyBytes

	rsaKey, err = cryptoutil.GenerateRSAKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	rsaKeyBytes, err = x509.MarshalPKCS8PrivateKey(rsaKey)
	if err != nil {
		return nil, err
	}
	keyMap[KeyType_RSA4096] = rsaKeyBytes

	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	ecdsaKeyBytes, err := x509.MarshalPKCS8PrivateKey(ecdsaKey)
	if err != nil {
		return nil, err
	}
	keyMap[KeyType_ECDSA_P256] = ecdsaKeyBytes

	_, ed25519Key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	ed25519KeyBytes, err := x509.MarshalPKCS8PrivateKey(ed25519Key)
	if err != nil {
		return nil, err
	}
	keyMap[KeyType_ED25519] = ed25519KeyBytes

	aesKey := make([]byte, 32)
	_, err = rand.Read(aesKey)
	if err != nil {
		return nil, err
	}
	keyMap[KeyType_AES256_GCM96] = aesKey

	return keyMap, nil
}

func BenchmarkSymmetric(b *testing.B) {
	ctx := context.Background()
	lm, _ := NewLockManager(true, 0)
	storage := &logical.InmemStorage{}
	p, _, _ := lm.GetPolicy(ctx, PolicyRequest{
		Upsert:  true,
		Storage: storage,
		KeyType: KeyType_AES256_GCM96,
		Name:    "test",
	}, rand.Reader)
	key, _ := p.GetKey(nil, 1, 32)
	pt := make([]byte, 10)
	ad := make([]byte, 10)
	for i := 0; i < b.N; i++ {
		ct, _ := p.SymmetricEncryptRaw(1, key, pt,
			SymmetricOpts{
				AdditionalData: ad,
			})
		pt2, _ := p.SymmetricDecryptRaw(key, ct, SymmetricOpts{
			AdditionalData: ad,
		})
		if !bytes.Equal(pt, pt2) {
			b.Fail()
		}
	}
}

func saltOptions(options SigningOptions, saltLength int) SigningOptions {
	return SigningOptions{
		HashAlgorithm: options.HashAlgorithm,
		Marshaling:    options.Marshaling,
		SaltLength:    saltLength,
		SigAlgorithm:  options.SigAlgorithm,
	}
}

func manualVerify(depth int, t *testing.T, p *Policy, input []byte, sig *SigningResult, options SigningOptions) {
	tabs := strings.Repeat("\t", depth)
	t.Log(tabs, "Manually verifying signature with options:", options)

	tabs = strings.Repeat("\t", depth+1)
	verified, err := p.VerifySignatureWithOptions(nil, input, sig.Signature, &options)
	if err != nil {
		t.Fatal(tabs, "❌ Failed to manually verify signature:", err)
	}
	if !verified {
		t.Fatal(tabs, "❌ Failed to manually verify signature")
	}
}

func autoVerify(depth int, t *testing.T, p *Policy, input []byte, sig *SigningResult, options SigningOptions) {
	tabs := strings.Repeat("\t", depth)
	t.Log(tabs, "Automatically verifying signature with options:", options)

	tabs = strings.Repeat("\t", depth+1)
	verified, err := p.VerifySignature(nil, input, options.HashAlgorithm, options.SigAlgorithm, options.Marshaling, sig.Signature)
	if err != nil {
		t.Fatal(tabs, "❌ Failed to automatically verify signature:", err)
	}
	if !verified {
		t.Fatal(tabs, "❌ Failed to automatically verify signature")
	}
}

func autoVerifyDecrypt(depth int, t *testing.T, p *Policy, input []byte, ct string, factories ...any) {
	tabs := strings.Repeat("\t", depth)
	t.Log(tabs, "Automatically decrypting with options:", factories)

	tabs = strings.Repeat("\t", depth+1)
	ptb64, err := p.DecryptWithFactory(nil, nil, ct, factories...)
	if err != nil {
		t.Fatal(tabs, "❌ Failed to automatically verify signature:", err)
	}

	pt, err := base64.StdEncoding.DecodeString(ptb64)
	if err != nil {
		t.Fatal(tabs, "❌ Failed decoding plaintext:", err)
	}
	if !bytes.Equal(input, pt) {
		t.Fatal(tabs, "❌ Failed to automatically decrypt")
	}
}

func Test_RSA_PSS(t *testing.T) {
	t.Log("Testing RSA PSS")
	mathrand.Seed(time.Now().UnixNano())

	var userError errutil.UserError
	ctx := context.Background()
	storage := &logical.InmemStorage{}
	// https://crypto.stackexchange.com/a/1222
	input := []byte("the ancients say the longer the salt, the more provable the security")
	sigAlgorithm := "pss"

	tabs := make(map[int]string)
	for i := 1; i <= 6; i++ {
		tabs[i] = strings.Repeat("\t", i)
	}

	test_RSA_PSS := func(t *testing.T, p *Policy, rsaKey *rsa.PrivateKey, hashType HashType,
		marshalingType MarshalingType,
	) {
		unsaltedOptions := SigningOptions{
			HashAlgorithm: hashType,
			Marshaling:    marshalingType,
			SigAlgorithm:  sigAlgorithm,
		}
		cryptoHash := CryptoHashMap[hashType]
		minSaltLength := p.minRSAPSSSaltLength()
		maxSaltLength := p.maxRSAPSSSaltLength(rsaKey.N.BitLen(), cryptoHash)
		hash := cryptoHash.New()
		hash.Write(input)
		input = hash.Sum(nil)

		// 1. Make an "automatic" signature with the given key size and hash algorithm,
		// but an automatically chosen salt length.
		t.Log(tabs[3], "Make an automatic signature")
		sig, err := p.Sign(0, nil, input, hashType, sigAlgorithm, marshalingType)
		if err != nil {
			// A bit of a hack but FIPS go does not support some hash types
			if isUnsupportedGoHashType(hashType, err) {
				t.Skip(tabs[4], "skipping test as FIPS Go does not support hash type")
				return
			}
			t.Fatal(tabs[4], "❌ Failed to automatically sign:", err)
		}

		// 1.1 Verify this automatic signature using the *inferred* salt length.
		autoVerify(4, t, p, input, sig, unsaltedOptions)

		// 1.2. Verify this automatic signature using the *correct, given* salt length.
		manualVerify(4, t, p, input, sig, saltOptions(unsaltedOptions, maxSaltLength))

		// 1.3. Try to verify this automatic signature using *incorrect, given* salt lengths.
		t.Log(tabs[4], "Test incorrect salt lengths")
		incorrectSaltLengths := []int{minSaltLength, maxSaltLength - 1}
		for _, saltLength := range incorrectSaltLengths {
			t.Log(tabs[5], "Salt length:", saltLength)
			saltedOptions := saltOptions(unsaltedOptions, saltLength)

			verified, _ := p.VerifySignatureWithOptions(nil, input, sig.Signature, &saltedOptions)
			if verified {
				t.Fatal(tabs[6], "❌ Failed to invalidate", verified, "signature using incorrect salt length:", err)
			}
		}

		// 2. Rule out boundary, invalid salt lengths.
		t.Log(tabs[3], "Test invalid salt lengths")
		invalidSaltLengths := []int{minSaltLength - 1, maxSaltLength + 1}
		for _, saltLength := range invalidSaltLengths {
			t.Log(tabs[4], "Salt length:", saltLength)
			saltedOptions := saltOptions(unsaltedOptions, saltLength)

			// 2.1. Fail to sign.
			t.Log(tabs[5], "Try to make a manual signature")
			_, err := p.SignWithOptions(0, nil, input, &saltedOptions)
			if !errors.As(err, &userError) {
				t.Fatal(tabs[6], "❌ Failed to reject invalid salt length:", err)
			}

			// 2.2. Fail to verify.
			t.Log(tabs[5], "Try to verify an automatic signature using an invalid salt length")
			_, err = p.VerifySignatureWithOptions(nil, input, sig.Signature, &saltedOptions)
			if !errors.As(err, &userError) {
				t.Fatal(tabs[6], "❌ Failed to reject invalid salt length:", err)
			}
		}

		// 3. For three possible valid salt lengths...
		t.Log(tabs[3], "Test three possible valid salt lengths")
		midSaltLength := mathrand.Intn(maxSaltLength-1) + 1 // [1, maxSaltLength)
		validSaltLengths := []int{minSaltLength, midSaltLength, maxSaltLength}
		for _, saltLength := range validSaltLengths {
			t.Log(tabs[4], "Salt length:", saltLength)
			saltedOptions := saltOptions(unsaltedOptions, saltLength)

			// 3.1. Make a "manual" signature with the given key size, hash algorithm, and salt length.
			t.Log(tabs[5], "Make a manual signature")
			sig, err := p.SignWithOptions(0, nil, input, &saltedOptions)
			if err != nil {
				t.Fatal(tabs[6], "❌ Failed to manually sign:", err)
			}

			// 3.2. Verify this manual signature using the *correct, given* salt length.
			manualVerify(6, t, p, input, sig, saltedOptions)

			// 3.3. Verify this manual signature using the *inferred* salt length.
			autoVerify(6, t, p, input, sig, unsaltedOptions)
		}
	}

	rsaKeyTypes := []KeyType{KeyType_RSA2048, KeyType_RSA3072, KeyType_RSA4096}
	testKeys, err := generateTestKeys()
	if err != nil {
		t.Fatalf("error generating test keys: %s", err)
	}

	// 1. For each standard RSA key size 2048, 3072, and 4096...
	for _, rsaKeyType := range rsaKeyTypes {
		t.Log("Key size: ", rsaKeyType)
		p := &Policy{
			Name: fmt.Sprint(rsaKeyType), // NOTE: crucial to create a new key per key size
			Type: rsaKeyType,
		}

		rsaKeyBytes := testKeys[rsaKeyType]
		err := p.Import(ctx, storage, rsaKeyBytes, rand.Reader)
		if err != nil {
			t.Fatal(tabs[1], "❌ Failed to import key:", err)
		}
		rsaKeyAny, err := x509.ParsePKCS8PrivateKey(rsaKeyBytes)
		if err != nil {
			t.Fatalf("error parsing test keys: %s", err)
		}
		rsaKey := rsaKeyAny.(*rsa.PrivateKey)

		// 2. For each hash algorithm...
		for hashAlgorithm, hashType := range HashTypeMap {
			t.Log(tabs[1], "Hash algorithm:", hashAlgorithm)
			if hashAlgorithm == "none" {
				continue
			}

			// 3. For each marshaling type...
			for marshalingName, marshalingType := range MarshalingTypeMap {
				t.Log(tabs[2], "Marshaling type:", marshalingName)
				testName := fmt.Sprintf("%s-%s-%s", rsaKeyType, hashAlgorithm, marshalingName)
				t.Run(testName, func(t *testing.T) { test_RSA_PSS(t, p, rsaKey, hashType, marshalingType) })
			}
		}
	}
}

func Test_RSA_PKCS1Encryption(t *testing.T) {
	t.Log("Testing RSA PKCS#1v1.5 padded encryption")

	ctx := context.Background()
	storage := &logical.InmemStorage{}
	// https://crypto.stackexchange.com/a/1222
	pt := []byte("Sphinx of black quartz, judge my vow")
	input := base64.StdEncoding.EncodeToString(pt)

	tabs := make(map[int]string)
	for i := 1; i <= 6; i++ {
		tabs[i] = strings.Repeat("\t", i)
	}

	test_RSA_PKCS1 := func(t *testing.T, p *Policy, rsaKey *rsa.PrivateKey, padding PaddingScheme) {
		// 1. Make a signature with the given key size and hash algorithm.
		t.Log(tabs[3], "Make an automatic signature")
		ct, err := p.EncryptWithFactory(0, nil, nil, string(input), padding)
		if err != nil {
			t.Fatal(tabs[4], "❌ Failed to automatically encrypt:", err)
		}

		// 1.1 Verify this signature using the *inferred* salt length.
		autoVerifyDecrypt(4, t, p, pt, ct, padding)
	}

	rsaKeyTypes := []KeyType{KeyType_RSA2048, KeyType_RSA3072, KeyType_RSA4096}
	testKeys, err := generateTestKeys()
	if err != nil {
		t.Fatalf("error generating test keys: %s", err)
	}

	// 1. For each standard RSA key size 2048, 3072, and 4096...
	for _, rsaKeyType := range rsaKeyTypes {
		t.Log("Key size: ", rsaKeyType)
		p := &Policy{
			Name: fmt.Sprint(rsaKeyType), // NOTE: crucial to create a new key per key size
			Type: rsaKeyType,
		}

		rsaKeyBytes := testKeys[rsaKeyType]
		err := p.Import(ctx, storage, rsaKeyBytes, rand.Reader)
		if err != nil {
			t.Fatal(tabs[1], "❌ Failed to import key:", err)
		}
		rsaKeyAny, err := x509.ParsePKCS8PrivateKey(rsaKeyBytes)
		if err != nil {
			t.Fatalf("error parsing test keys: %s", err)
		}
		rsaKey := rsaKeyAny.(*rsa.PrivateKey)
		for _, padding := range []PaddingScheme{PaddingScheme_OAEP, PaddingScheme_PKCS1v15, ""} {
			t.Run(fmt.Sprintf("%s/%s", rsaKeyType.String(), padding), func(t *testing.T) { test_RSA_PKCS1(t, p, rsaKey, padding) })
		}
	}
}

func Test_RSA_PKCS1Signing(t *testing.T) {
	t.Log("Testing RSA PKCS#1v1.5 signatures")

	ctx := context.Background()
	storage := &logical.InmemStorage{}
	// https://crypto.stackexchange.com/a/1222
	input := []byte("Sphinx of black quartz, judge my vow")
	sigAlgorithm := "pkcs1v15"

	tabs := make(map[int]string)
	for i := 1; i <= 6; i++ {
		tabs[i] = strings.Repeat("\t", i)
	}

	test_RSA_PKCS1 := func(t *testing.T, p *Policy, rsaKey *rsa.PrivateKey, hashType HashType,
		marshalingType MarshalingType,
	) {
		unsaltedOptions := SigningOptions{
			HashAlgorithm: hashType,
			Marshaling:    marshalingType,
			SigAlgorithm:  sigAlgorithm,
		}
		cryptoHash := CryptoHashMap[hashType]

		// PKCS#1v1.5 NoOID uses a direct input and assumes it is pre-hashed.
		if hashType != 0 {
			hash := cryptoHash.New()
			hash.Write(input)
			input = hash.Sum(nil)
		}

		// 1. Make a signature with the given key size and hash algorithm.
		t.Log(tabs[3], "Make an automatic signature")
		sig, err := p.Sign(0, nil, input, hashType, sigAlgorithm, marshalingType)
		if err != nil {
			// A bit of a hack but FIPS go does not support some hash types
			if isUnsupportedGoHashType(hashType, err) {
				t.Skip(tabs[4], "skipping test as FIPS Go does not support hash type")
				return
			}
			t.Fatal(tabs[4], "❌ Failed to automatically sign:", err)
		}

		// 1.1 Verify this signature using the *inferred* salt length.
		autoVerify(4, t, p, input, sig, unsaltedOptions)
	}

	rsaKeyTypes := []KeyType{KeyType_RSA2048, KeyType_RSA3072, KeyType_RSA4096}
	testKeys, err := generateTestKeys()
	if err != nil {
		t.Fatalf("error generating test keys: %s", err)
	}

	// 1. For each standard RSA key size 2048, 3072, and 4096...
	for _, rsaKeyType := range rsaKeyTypes {
		t.Log("Key size: ", rsaKeyType)
		p := &Policy{
			Name: fmt.Sprint(rsaKeyType), // NOTE: crucial to create a new key per key size
			Type: rsaKeyType,
		}

		rsaKeyBytes := testKeys[rsaKeyType]
		err := p.Import(ctx, storage, rsaKeyBytes, rand.Reader)
		if err != nil {
			t.Fatal(tabs[1], "❌ Failed to import key:", err)
		}
		rsaKeyAny, err := x509.ParsePKCS8PrivateKey(rsaKeyBytes)
		if err != nil {
			t.Fatalf("error parsing test keys: %s", err)
		}
		rsaKey := rsaKeyAny.(*rsa.PrivateKey)

		// 2. For each hash algorithm...
		for hashAlgorithm, hashType := range HashTypeMap {
			t.Log(tabs[1], "Hash algorithm:", hashAlgorithm)

			// 3. For each marshaling type...
			for marshalingName, marshalingType := range MarshalingTypeMap {
				t.Log(tabs[2], "Marshaling type:", marshalingName)
				testName := fmt.Sprintf("%s-%s-%s", rsaKeyType, hashAlgorithm, marshalingName)
				t.Run(testName, func(t *testing.T) { test_RSA_PKCS1(t, p, rsaKey, hashType, marshalingType) })
			}
		}
	}
}

// Normal Go builds support all the hash functions for RSA_PSS signatures but the
// FIPS Go build does not support at this time the SHA3 hashes as FIPS 140_2 does
// not accept them.
func isUnsupportedGoHashType(hashType HashType, err error) bool {
	switch hashType {
	case HashTypeSHA3224, HashTypeSHA3256, HashTypeSHA3384, HashTypeSHA3512:
		return strings.Contains(err.Error(), "unsupported hash function")
	}

	return false
}
