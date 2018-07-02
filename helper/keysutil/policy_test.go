package keysutil

import (
	"context"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/logical"
	"github.com/mitchellh/copystructure"
)

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
	testKeyUpgradeCommon(t, NewLockManager(false))
	testKeyUpgradeCommon(t, NewLockManager(true))
}

func testKeyUpgradeCommon(t *testing.T, lm *LockManager) {
	ctx := context.Background()

	storage := &logical.InmemStorage{}
	p, upserted, err := lm.GetPolicy(ctx, PolicyRequest{
		Upsert:  true,
		Storage: storage,
		KeyType: KeyType_AES256_GCM96,
		Name:    "test",
	})
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
	testArchivingUpgradeCommon(t, NewLockManager(false))
	testArchivingUpgradeCommon(t, NewLockManager(true))
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
	})
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
	keysArchive := []KeyEntry{KeyEntry{}, p.Keys["1"]}
	checkKeys(t, ctx, p, storage, keysArchive, "initial", 1, 1, 1)

	for i := 2; i <= 10; i++ {
		err = p.Rotate(ctx, storage)
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
	})
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
	})
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
	})
	if err != nil {
		t.Fatal(err)
	}
	if p != nil {
		t.Fatal("policy not nil after delete")
	}
}

func Test_Archiving(t *testing.T) {
	testArchivingCommon(t, NewLockManager(false))
	testArchivingCommon(t, NewLockManager(true))
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
	})
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
	keysArchive := []KeyEntry{KeyEntry{}, p.Keys["1"]}
	checkKeys(t, ctx, p, storage, keysArchive, "initial", 1, 1, 1)

	for i := 2; i <= 10; i++ {
		err = p.Rotate(ctx, storage)
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
	archiveVer, latestVer, keysSize int) {

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
		// Travis has weird time zone issues and gets super unhappy
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
	lm := NewLockManager(false)

	storage := &logical.InmemStorage{}
	p, _, err := lm.GetPolicy(ctx, PolicyRequest{
		Upsert:  true,
		Storage: storage,
		KeyType: KeyType_AES256_GCM96,
		Name:    "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	if p == nil {
		t.Fatal("nil policy")
	}

	// Store the initial key in the archive
	keysArchive := []KeyEntry{KeyEntry{}, p.Keys["1"]}
	checkKeys(t, ctx, p, storage, keysArchive, "initial", 1, 1, 1)

	// We use checkKeys here just for sanity; it doesn't really handle cases of
	// errors below so we do more targeted testing later
	for i := 2; i <= 5; i++ {
		err = p.Rotate(ctx, storage)
		if err != nil {
			t.Fatal(err)
		}
		keysArchive = append(keysArchive, p.Keys[strconv.Itoa(i)])
		checkKeys(t, ctx, p, storage, keysArchive, "rotate", i, i, i)
	}

	underlying := storage.Underlying()
	underlying.FailPut(true)

	priorLen := len(p.Keys)

	err = p.Rotate(ctx, storage)
	if err == nil {
		t.Fatal("expected error")
	}

	if len(p.Keys) != priorLen {
		t.Fatal("length of keys should not have changed")
	}
}

func Test_BadUpgrade(t *testing.T) {
	ctx := context.Background()
	lm := NewLockManager(false)
	storage := &logical.InmemStorage{}
	p, _, err := lm.GetPolicy(ctx, PolicyRequest{
		Upsert:  true,
		Storage: storage,
		KeyType: KeyType_AES256_GCM96,
		Name:    "test",
	})
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

	if err := p.Upgrade(ctx, storage); err != nil {
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

	if err := p.Upgrade(ctx, storage); err == nil {
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
	lm := NewLockManager(false)
	storage := &logical.InmemStorage{}
	p, _, err := lm.GetPolicy(ctx, PolicyRequest{
		Upsert:  true,
		Storage: storage,
		KeyType: KeyType_AES256_GCM96,
		Name:    "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	if p == nil {
		t.Fatal("nil policy")
	}

	for i := 2; i <= 10; i++ {
		err = p.Rotate(ctx, storage)
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
