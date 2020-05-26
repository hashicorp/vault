package gocb

import (
	"errors"
	"fmt"
)

// CouchbaseList represents a list document.
type CouchbaseList struct {
	collection *Collection
	id         string
}

// List returns a new CouchbaseList for the document specified by id.
func (c *Collection) List(id string) *CouchbaseList {
	return &CouchbaseList{
		collection: c,
		id:         id,
	}
}

// Iterator returns an iterable for all items in the list.
func (cl *CouchbaseList) Iterator() ([]interface{}, error) {
	content, err := cl.collection.Get(cl.id, nil)
	if err != nil {
		return nil, err
	}

	var listContents []interface{}
	err = content.Content(&listContents)
	if err != nil {
		return nil, err
	}

	return listContents, nil
}

// At retrieves the value specified at the given index from the list.
func (cl *CouchbaseList) At(index int, valuePtr interface{}) error {
	ops := make([]LookupInSpec, 1)
	ops[0] = GetSpec(fmt.Sprintf("[%d]", index), nil)
	result, err := cl.collection.LookupIn(cl.id, ops, nil)
	if err != nil {
		return err
	}

	return result.ContentAt(0, valuePtr)
}

// RemoveAt removes the value specified at the given index from the list.
func (cl *CouchbaseList) RemoveAt(index int) error {
	ops := make([]MutateInSpec, 1)
	ops[0] = RemoveSpec(fmt.Sprintf("[%d]", index), nil)
	_, err := cl.collection.MutateIn(cl.id, ops, nil)
	if err != nil {
		return err
	}

	return nil
}

// Append appends an item to the list.
func (cl *CouchbaseList) Append(val interface{}) error {
	ops := make([]MutateInSpec, 1)
	ops[0] = ArrayAppendSpec("", val, nil)
	_, err := cl.collection.MutateIn(cl.id, ops, &MutateInOptions{StoreSemantic: StoreSemanticsUpsert})
	if err != nil {
		return err
	}

	return nil
}

// Prepend prepends an item to the list.
func (cl *CouchbaseList) Prepend(val interface{}) error {
	ops := make([]MutateInSpec, 1)
	ops[0] = ArrayPrependSpec("", val, nil)
	_, err := cl.collection.MutateIn(cl.id, ops, &MutateInOptions{StoreSemantic: StoreSemanticsUpsert})
	if err != nil {
		return err
	}

	return nil
}

// IndexOf gets the index of the item in the list.
func (cl *CouchbaseList) IndexOf(val interface{}) (int, error) {
	content, err := cl.collection.Get(cl.id, nil)
	if err != nil {
		return 0, err
	}

	var listContents []interface{}
	err = content.Content(&listContents)
	if err != nil {
		return 0, err
	}

	for i, item := range listContents {
		if item == val {
			return i, nil
		}
	}

	return -1, nil
}

// Size returns the size of the list.
func (cl *CouchbaseList) Size() (int, error) {
	ops := make([]LookupInSpec, 1)
	ops[0] = CountSpec("", nil)
	result, err := cl.collection.LookupIn(cl.id, ops, nil)
	if err != nil {
		return 0, err
	}

	var count int
	err = result.ContentAt(0, &count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Clear clears a list, also removing it.
func (cl *CouchbaseList) Clear() error {
	_, err := cl.collection.Remove(cl.id, nil)
	if err != nil {
		return err
	}

	return nil
}

// CouchbaseMap represents a map document.
type CouchbaseMap struct {
	collection *Collection
	id         string
}

// Map returns a new CouchbaseMap.
func (c *Collection) Map(id string) *CouchbaseMap {
	return &CouchbaseMap{
		collection: c,
		id:         id,
	}
}

// Iterator returns an iterable for all items in the map.
func (cl *CouchbaseMap) Iterator() (map[string]interface{}, error) {
	content, err := cl.collection.Get(cl.id, nil)
	if err != nil {
		return nil, err
	}

	var mapContents map[string]interface{}
	err = content.Content(&mapContents)
	if err != nil {
		return nil, err
	}

	return mapContents, nil
}

// At retrieves the item for the given id from the map.
func (cl *CouchbaseMap) At(id string, valuePtr interface{}) error {
	ops := make([]LookupInSpec, 1)
	ops[0] = GetSpec(fmt.Sprintf("[%s]", id), nil)
	result, err := cl.collection.LookupIn(cl.id, ops, nil)
	if err != nil {
		return err
	}

	return result.ContentAt(0, valuePtr)
}

// Add adds an item to the map.
func (cl *CouchbaseMap) Add(id string, val interface{}) error {
	ops := make([]MutateInSpec, 1)
	ops[0] = UpsertSpec(id, val, nil)
	_, err := cl.collection.MutateIn(cl.id, ops, &MutateInOptions{StoreSemantic: StoreSemanticsUpsert})
	if err != nil {
		return err
	}

	return nil
}

// Remove removes an item from the map.
func (cl *CouchbaseMap) Remove(id string) error {
	ops := make([]MutateInSpec, 1)
	ops[0] = RemoveSpec(id, nil)
	_, err := cl.collection.MutateIn(cl.id, ops, nil)
	if err != nil {
		return err
	}

	return nil
}

// Exists verifies whether or a id exists in the map.
func (cl *CouchbaseMap) Exists(id string) (bool, error) {
	ops := make([]LookupInSpec, 1)
	ops[0] = ExistsSpec(fmt.Sprintf("[%s]", id), nil)
	result, err := cl.collection.LookupIn(cl.id, ops, nil)
	if err != nil {
		return false, err
	}

	return result.Exists(0), nil
}

// Size returns the size of the map.
func (cl *CouchbaseMap) Size() (int, error) {
	ops := make([]LookupInSpec, 1)
	ops[0] = CountSpec("", nil)
	result, err := cl.collection.LookupIn(cl.id, ops, nil)
	if err != nil {
		return 0, err
	}

	var count int
	err = result.ContentAt(0, &count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Keys returns all of the keys within the map.
func (cl *CouchbaseMap) Keys() ([]string, error) {
	content, err := cl.collection.Get(cl.id, nil)
	if err != nil {
		return nil, err
	}

	var mapContents map[string]interface{}
	err = content.Content(&mapContents)
	if err != nil {
		return nil, err
	}

	var keys []string
	for id := range mapContents {
		keys = append(keys, id)
	}

	return keys, nil
}

// Values returns all of the values within the map.
func (cl *CouchbaseMap) Values() ([]interface{}, error) {
	content, err := cl.collection.Get(cl.id, nil)
	if err != nil {
		return nil, err
	}

	var mapContents map[string]interface{}
	err = content.Content(&mapContents)
	if err != nil {
		return nil, err
	}

	var values []interface{}
	for _, val := range mapContents {
		values = append(values, val)
	}

	return values, nil
}

// Clear clears a map, also removing it.
func (cl *CouchbaseMap) Clear() error {
	_, err := cl.collection.Remove(cl.id, nil)
	if err != nil {
		return err
	}

	return nil
}

// CouchbaseSet represents a set document.
type CouchbaseSet struct {
	id         string
	underlying *CouchbaseList
}

// Set returns a new CouchbaseSet.
func (c *Collection) Set(id string) *CouchbaseSet {
	return &CouchbaseSet{
		id:         id,
		underlying: c.List(id),
	}
}

// Iterator returns an iterable for all items in the set.
func (cs *CouchbaseSet) Iterator() ([]interface{}, error) {
	return cs.underlying.Iterator()
}

// Add adds a value to the set.
func (cs *CouchbaseSet) Add(val interface{}) error {
	ops := make([]MutateInSpec, 1)
	ops[0] = ArrayAddUniqueSpec("", val, nil)
	_, err := cs.underlying.collection.MutateIn(cs.id, ops, &MutateInOptions{StoreSemantic: StoreSemanticsUpsert})
	if err != nil {
		return err
	}

	return nil
}

// Remove removes an value from the set.
func (cs *CouchbaseSet) Remove(val string) error {
	for i := 0; i < 16; i++ {
		content, err := cs.underlying.collection.Get(cs.id, nil)
		if err != nil {
			return err
		}

		cas := content.Cas()

		var setContents []interface{}
		err = content.Content(&setContents)
		if err != nil {
			return err
		}

		indexToRemove := -1
		for i, item := range setContents {
			if item == val {
				indexToRemove = i
			}
		}

		if indexToRemove > -1 {
			ops := make([]MutateInSpec, 1)
			ops[0] = RemoveSpec(fmt.Sprintf("[%d]", indexToRemove), nil)
			_, err = cs.underlying.collection.MutateIn(cs.id, ops, &MutateInOptions{Cas: cas})
			if errors.Is(err, ErrCasMismatch) {
				continue
			}
			if err != nil {
				return err
			}
		}
		return nil
	}

	return errors.New("failed to perform operation after 16 retries")
}

// Values returns all of the values within the set.
func (cs *CouchbaseSet) Values() ([]interface{}, error) {
	content, err := cs.underlying.collection.Get(cs.id, nil)
	if err != nil {
		return nil, err
	}

	var setContents []interface{}
	err = content.Content(&setContents)
	if err != nil {
		return nil, err
	}

	return setContents, nil
}

// Contains verifies whether or not a value exists within the set.
func (cs *CouchbaseSet) Contains(val string) (bool, error) {
	content, err := cs.underlying.collection.Get(cs.id, nil)
	if err != nil {
		return false, err
	}

	var setContents []interface{}
	err = content.Content(&setContents)
	if err != nil {
		return false, err
	}

	for _, item := range setContents {
		if item == val {
			return true, nil
		}
	}

	return false, nil
}

// Size returns the size of the set
func (cs *CouchbaseSet) Size() (int, error) {
	return cs.underlying.Size()
}

// Clear clears a set, also removing it.
func (cs *CouchbaseSet) Clear() error {
	err := cs.underlying.Clear()
	if err != nil {
		return err
	}

	return nil
}

// CouchbaseQueue represents a queue document.
type CouchbaseQueue struct {
	id         string
	underlying *CouchbaseList
}

// Queue returns a new CouchbaseQueue.
func (c *Collection) Queue(id string) *CouchbaseQueue {
	return &CouchbaseQueue{
		id:         id,
		underlying: c.List(id),
	}
}

// Iterator returns an iterable for all items in the queue.
func (cs *CouchbaseQueue) Iterator() ([]interface{}, error) {
	return cs.underlying.Iterator()
}

// Push pushes a value onto the queue.
func (cs *CouchbaseQueue) Push(val interface{}) error {
	return cs.underlying.Prepend(val)
}

// Pop pops an items off of the queue.
func (cs *CouchbaseQueue) Pop(valuePtr interface{}) error {
	for i := 0; i < 16; i++ {
		ops := make([]LookupInSpec, 1)
		ops[0] = GetSpec("[-1]", nil)
		content, err := cs.underlying.collection.LookupIn(cs.id, ops, nil)
		if err != nil {
			return err
		}

		cas := content.Cas()
		err = content.ContentAt(0, valuePtr)
		if err != nil {
			return err
		}

		mutateOps := make([]MutateInSpec, 1)
		mutateOps[0] = RemoveSpec("[-1]", nil)
		_, err = cs.underlying.collection.MutateIn(cs.id, mutateOps, &MutateInOptions{Cas: cas})
		if errors.Is(err, ErrCasMismatch) {
			continue
		}
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("failed to perform operation after 16 retries")
}

// Size returns the size of the queue.
func (cs *CouchbaseQueue) Size() (int, error) {
	return cs.underlying.Size()
}

// Clear clears a queue, also removing it.
func (cs *CouchbaseQueue) Clear() error {
	err := cs.underlying.Clear()
	if err != nil {
		return err
	}

	return nil
}
