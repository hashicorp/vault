package storagepacker

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/hashicorp/vault/sdk/helper/cryptoutil"
)

func GetItemIDHash(itemID string) string {
	return hex.EncodeToString(cryptoutil.Blake2b256Hash(itemID))
}

// Length of key
const KeyLength = 64

// firstKey returns the topmost bucket in the tree for a given key.
// Used as a default if the cache is empty or bypassed.
func (s *StoragePackerV2) firstKey(cacheKey string) (string, error) {
	rootShardLength := s.BaseBucketBits / 4
	if len(cacheKey) < rootShardLength {
		return cacheKey, errors.New("key too short")
	}
	return cacheKey[0 : s.BaseBucketBits/4], nil
}

// getAllBaseBucketKeys returns all topmost buckets in the tree.
func (s *StoragePackerV2) getAllBaseBucketKeys() []string {
	numBuckets := int(math.Pow(2.0, float64(s.BaseBucketBits)))
	rootBucketLength := s.BaseBucketBits / 4

	// %02x for default configuration, could be %01x, %03x, etc.
	formatString := fmt.Sprintf("%%0%dx", rootBucketLength)

	ret := make([]string, 0, numBuckets)
	for i := 0; i < numBuckets; i++ {
		bucketKey := fmt.Sprintf(formatString, i)
		ret = append(ret, bucketKey)
	}
	return ret
}

// GetCacheKey returns the radix tree key corresponding to a bucket.
// Buckets keys have / in them.
// Entries in the radix tree do not.
// Lock hashing uses the latter form.
func (s *StoragePackerV2) GetCacheKey(key string) string {
	return strings.Replace(key, "/", "", -1)
}
