package storagepacker

import (
	"encoding/hex"
	"errors"
	"github.com/hashicorp/vault/sdk/helper/cryptoutil"
	"sort"
	"strings"
)

type itemRequest struct {
	// Item ID, provided by client
	ID string

	// Storage key == hash of ID
	Key string

	// Stored object, nil if not found
	Value *Item

	// Bucket responsible for this key
	Bucket *Bucket
}

func GetItemIDHash(itemID string) string {
	return hex.EncodeToString(cryptoutil.Blake2b256Hash(itemID))
}

// Given a list of keys, calculate their keys and sort the
// resulting array of itemRequests by key.
func (s *StoragePackerV2) keysForIDs(ids []string) []*itemRequest {
	requests := make([]*itemRequest, 0, len(ids))
	for _, id := range ids {
		requests = append(requests, &itemRequest{
			ID:     id,
			Key:    GetItemIDHash(id),
			Value:  nil,
			Bucket: nil,
		})
	}
	sort.Slice(requests, func(i, j int) bool {
		return requests[i].Key < requests[j].Key
	})
	return requests
}

func checkForDuplicateIds(ids []string) (bool, string) {
	idsSeen := make(map[string]bool, len(ids))
	for _, id := range ids {
		if _, found := idsSeen[id]; found {
			return true, id
		}
		idsSeen[id] = true
	}
	return false, ""
}

// Return the topmost bucket in the tree for a given key.
// Used as a defult if the cache is empty or bypassed.
func (s *StoragePackerV2) firstKey(cacheKey string) (string, error) {
	rootShardLength := s.BaseBucketBits / 4
	if len(cacheKey) < rootShardLength {
		return cacheKey, errors.New("Key too short.")
	}
	return cacheKey[0 : s.BaseBucketBits/4], nil
}

// Buckets keys have / in them.
// Entries in the radix tree do not.
// Lock hashing uses the latter form.
func (s *StoragePackerV2) GetCacheKey(key string) string {
	return strings.Replace(key, "/", "", -1)
}
