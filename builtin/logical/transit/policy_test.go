package transit

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/logical"
)

var (
	keysArchive []KeyEntry
)

func resetKeysArchive() {
	keysArchive = []KeyEntry{KeyEntry{}}
}

func Test_KeyUpgrade(t *testing.T) {
	storage := &logical.InmemStorage{}
	policies := &policyCache{
		cache: map[string]*lockingPolicy{},
	}
	lp, err := policies.generatePolicy(storage, "test", false)
	if err != nil {
		t.Fatal(err)
	}
	if lp == nil {
		t.Fatal("nil policy")
	}

	policy := lp.policy

	testBytes := make([]byte, len(policy.Keys[1].Key))
	copy(testBytes, policy.Keys[1].Key)

	policy.Key = policy.Keys[1].Key
	policy.Keys = nil
	policy.migrateKeyToKeysMap()
	if policy.Key != nil {
		t.Fatal("policy.Key is not nil")
	}
	if len(policy.Keys) != 1 {
		t.Fatal("policy.Keys is the wrong size")
	}
	if !reflect.DeepEqual(testBytes, policy.Keys[1].Key) {
		t.Fatal("key mismatch")
	}
}

func Test_ArchivingUpgrade(t *testing.T) {
	resetKeysArchive()

	// First, we generate a policy and rotate it a number of times. Each time
	// we'll ensure that we have the expected number of keys in the archive and
	// the main keys object, which without changing the min version should be
	// zero and latest, respectively

	storage := &logical.InmemStorage{}
	policies := &policyCache{
		cache: map[string]*lockingPolicy{},
	}

	lp, err := policies.generatePolicy(storage, "test", false)
	if err != nil {
		t.Fatal(err)
	}
	if lp == nil {
		t.Fatal("policy is nil")
	}

	policy := lp.policy

	// Store the initial key in the archive
	keysArchive = append(keysArchive, policy.Keys[1])
	checkKeys(t, policy, storage, "initial", 1, 1, 1)

	for i := 2; i <= 10; i++ {
		err = policy.rotate(storage)
		if err != nil {
			t.Fatal(err)
		}
		keysArchive = append(keysArchive, policy.Keys[i])
		checkKeys(t, policy, storage, "rotate", i, i, i)
	}

	// Now, wipe the archive and set the archive version to zero
	err = storage.Delete("archive/test")
	if err != nil {
		t.Fatal(err)
	}
	policy.ArchiveVersion = 0

	// Store it, but without calling persist, so we don't trigger
	// handleArchiving()
	buf, err := policy.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	// Write the policy into storage
	err = storage.Put(&logical.StorageEntry{
		Key:   "policy/" + policy.Name,
		Value: buf,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Expire from the cache since we modified it under-the-hood
	delete(policies.cache, "test")

	// Now get the policy again; the upgrade should happen automatically
	lp, err = policies.getPolicy(&logical.Request{
		Storage: storage,
	}, "test")
	if err != nil {
		t.Fatal(err)
	}
	if lp == nil {
		t.Fatal("policy is nil")
	}

	policy = lp.policy

	checkKeys(t, policy, storage, "upgrade", 10, 10, 10)
}

func Test_Archiving(t *testing.T) {
	resetKeysArchive()

	// First, we generate a policy and rotate it a number of times. Each time
	// we'll ensure that we have the expected number of keys in the archive and
	// the main keys object, which without changing the min version should be
	// zero and latest, respectively

	storage := &logical.InmemStorage{}

	policies := &policyCache{
		cache: map[string]*lockingPolicy{},
	}

	lp, err := policies.generatePolicy(storage, "test", false)
	if err != nil {
		t.Fatal(err)
	}
	if lp == nil {
		t.Fatal("policy is nil")
	}

	policy := lp.policy

	// Store the initial key in the archive
	keysArchive = append(keysArchive, policy.Keys[1])
	checkKeys(t, policy, storage, "initial", 1, 1, 1)

	for i := 2; i <= 10; i++ {
		err = policy.rotate(storage)
		if err != nil {
			t.Fatal(err)
		}
		keysArchive = append(keysArchive, policy.Keys[i])
		checkKeys(t, policy, storage, "rotate", i, i, i)
	}

	// Move the min decryption version up
	for i := 1; i <= 10; i++ {
		policy.MinDecryptionVersion = i

		err = policy.Persist(storage)
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
		checkKeys(t, policy, storage, "minadd", 10, 10, policy.LatestVersion-policy.MinDecryptionVersion+1)
	}

	// Move the min decryption version down
	for i := 10; i >= 1; i-- {
		policy.MinDecryptionVersion = i

		err = policy.Persist(storage)
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
		checkKeys(t, policy, storage, "minsub", 10, 10, policy.LatestVersion-policy.MinDecryptionVersion+1)
	}
}

func checkKeys(t *testing.T,
	policy *Policy,
	storage logical.Storage,
	action string,
	archiveVer, latestVer, keysSize int) {

	// Sanity check
	if len(keysArchive) != latestVer+1 {
		t.Fatalf("latest expected key version is %d, expected test keys archive size is %d, "+
			"but keys archive is of size %d", latestVer, latestVer+1, len(keysArchive))
	}

	archive, err := policy.loadArchive(storage)
	if err != nil {
		t.Fatal(err)
	}

	badArchiveVer := false
	if archiveVer == 0 {
		if len(archive.Keys) != 0 || policy.ArchiveVersion != 0 {
			badArchiveVer = true
		}
	} else {
		// We need to subtract one because we have the indexes match key
		// versions, which start at 1. So for an archive version of 1, we
		// actually have two entries -- a blank 0 entry, and the key at spot 1
		if archiveVer != len(archive.Keys)-1 || archiveVer != policy.ArchiveVersion {
			badArchiveVer = true
		}
	}
	if badArchiveVer {
		t.Fatalf(
			"expected archive version %d, found length of archive keys %d and policy archive version %d",
			archiveVer, len(archive.Keys), policy.ArchiveVersion,
		)
	}

	if latestVer != policy.LatestVersion {
		t.Fatalf(
			"expected latest version %d, found %d",
			latestVer, policy.LatestVersion,
		)
	}

	if keysSize != len(policy.Keys) {
		t.Fatalf(
			"expected keys size %d, found %d, action is %s, policy is \n%#v\n",
			keysSize, len(policy.Keys), action, policy,
		)
	}

	for i := policy.MinDecryptionVersion; i <= policy.LatestVersion; i++ {
		if _, ok := policy.Keys[i]; !ok {
			t.Fatalf(
				"expected key %d, did not find it in policy keys", i,
			)
		}
	}

	for i := policy.MinDecryptionVersion; i <= policy.LatestVersion; i++ {
		if !reflect.DeepEqual(policy.Keys[i], keysArchive[i]) {
			t.Fatalf("key %d not equivalent between policy keys and test keys archive", i)
		}
	}

	for i := 1; i < len(archive.Keys); i++ {
		if !reflect.DeepEqual(archive.Keys[i].Key, keysArchive[i].Key) {
			t.Fatalf("key %d not equivalent between policy archive and test keys archive", i)
		}
	}
}
