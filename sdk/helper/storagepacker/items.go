package storagepacker

import (
	"sort"
)

// Item is used to store or return byte slices in the storage packer.
type Item struct {
	// ID is the Item ID, provided by client
	ID string

	// Key is the identifier used internally, a hash of ID
	key string

	// Value is the stored object, nil if not found
	Value []byte
}

// ItemsForIDs calculates keys and returns an Item for each input ID.
func (s *StoragePackerV2) itemsForIDs(ids []string) []*Item {
	requests := make([]*Item, 0, len(ids))
	for _, id := range ids {
		requests = append(requests, &Item{
			ID:    id,
			key:   GetItemIDHash(id),
			Value: nil,
		})
	}
	return requests
}

// computerKeysForItems calculates keys for each input Item (and returns
// the original slice.)
func (s *StoragePackerV2) computeKeysForItems(items []*Item) []*Item {
	for _, i := range items {
		i.key = GetItemIDHash(i.ID)
	}
	return items
}

// Sort the requests in key order, nondestructively (so we can refer
// back to the original order.)
func sortRequests(requests []*Item) []*Item {
	duplicate := append([]*Item{}, requests...)
	sort.Slice(duplicate, func(i, j int) bool {
		return duplicate[i].key < duplicate[j].key
	})
	return duplicate
}
